package desktop

import (
	"fmt"
	"time"

	"transok/backend/app"
	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/google/uuid"
)

// homeHeaderHeight 是标题区预留的高度：内容从这里开始，标题绝对定位在上面。
const homeHeaderHeight float32 = 58

func (a *App) home(c *ui.Context) {
	p := a.palette()

	// 内容整页滚动，顶部留出比标题更高的边距，滚上来的内容会沉到标题下方。
	ui.Scroll(c).Grow(1).Scrollbars(false).TrackScroll(&a.scroll).Children(func() {
		ui.Column(c).FillWidth().Padding(homeHeaderHeight, 16, 0, 16).Children(func() {
			a.dropArea(c, p)
			spacer(c, 18)
			a.filesHeader(c, p)
			spacer(c, 14)
			a.filesList(c, p)
		})
		spacer(c, 130)
	})

	// 标题绝对定位在最上层，容器本身是一层从实心到透明的渐变，
	// 内容滚过时就像沉进背景里。
	ui.Overlay(c, func() {
		band := ui.Column(c).Absolute().Left(0).Right(0).Top(0).
			Height(homeHeaderHeight + 26)
		// 滚起来之后头部整块变实，文案不会和按钮叠在一起；回到顶部再淡回去。
		scrolled := clamp01(a.scroll.Y / 24)
		band.LinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
			{Color: withAlpha(p.Canvas, 1), At: 0},
			{Color: withAlpha(p.Canvas, 1), At: 0.72 + 0.1*scrolled},
			{Color: withAlpha(p.Canvas, 0.6*scrolled), At: 0.86 + 0.06*scrolled},
			{Color: withAlpha(p.Canvas, 0), At: 1},
		}})
		band.Children(func() {
			ui.Column(c).FillWidth().Padding(0, 16).Children(func() {
				if c.TitleBar().Height > 0 {
					ui.Box(c).FillWidth().Height(a.titleBarStrip(c)).DragWindow()
				}
				a.homeHeader(c, p)
			})
		})
	})

	ui.Overlay(c, func() {
		// 底部渐隐：列表滚到按钮下方时被吃掉，落点深色纯黑、浅色纯白。
		fade := ui.Box(c).Absolute().Left(0).Right(0).Bottom(0).Height(150).PassThrough()

		// 分享中这层逐渐透出橄榄色，像是从底部往上打的光。
		on := clamp01(fade.AnimateTo("share-glow", boolFloat(a.share), ui.Animation{
			Duration:  700 * time.Millisecond,
			Ease:      ui.EaseOut,
			NonMotion: true,
		}))
		fade.LinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
			{Color: withAlpha(p.Canvas, 0), At: 0},
			{Color: withAlpha(p.Canvas.Mix(p.Accent, 0.10*on), 0.92), At: 0.5},
			{Color: withAlpha(p.Fade.Mix(p.Accent, 0.20*on), 1), At: 1},
		}})

		if on > 0.01 {
			glow := ui.Box(c).Absolute().Left(0).Right(0).Bottom(0).Height(260).PassThrough()
			glow.RadialGradient(ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 1},
				Radius: 230,
				Stops: []ui.GradientStop{
					{Color: withAlpha(p.Accent, 0.17*on), At: 0},
					{Color: withAlpha(p.Accent, 0.06*on), At: 0.5},
					{Color: withAlpha(p.Accent, 0), At: 1},
				},
			})
		}

		ui.Row(c).Absolute().Left(0).Right(0).Bottom(26).Justify(ui.Center).Children(func() {
			if a.selectMode {
				a.selectBar(c, p)
				return
			}
			a.shareButton(c, p)
		})
	})

	// 抽屉常驻构建，关闭时靠 Exit 动画退场。
	a.captchaDrawer(c, p)
	// 弹窗常驻构建，关闭时靠 Exit 动画退场。
	a.noteDrawer(c, p)
	a.pureTextDrawer(c, p)
	a.shareDrawer(c, p)
	a.deleteDrawer(c, p)
	a.qrDrawer(c, p)
}

func (a *App) homeHeader(c *ui.Context, p Palette) {
	header := ui.Row(c).AlignItems(ui.Center).DragWindow()
	header.Children(func() {
		ui.Column(c).Grow(1).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(9).Children(func() {
				if icon := appIcon(); icon != nil {
					ui.Image(c, icon).Size(36, 36).Fit(ui.Cover).Radius(9)
				}
				ui.Row(c).AlignItems(ui.End).Gap(4).Children(func() {
					ui.Text(c, a.appInfo["name"]).FontSize(27).FontWeight(800).LetterSpacing(-1.2)
					ui.Box(c).Size(5, 5).Radius(1).Background(p.Accent).Margin(0, 0, 5, 0)
				})
			})
		})
		glassGroup(c, p).Children(func() {
			glassSegment(c, p, "device-desktop-search", func() { a.page = discoverRoute })
			glassSegment(c, p, "share", func() { a.shareOpen = true })
			glassSegment(c, p, "settings", func() { a.page = settingsRoute })
		})
	})
}

// dropArea 接住从系统拖进来的文件与文件夹，也提供选择入口。
func (a *App) dropArea(c *ui.Context, p Palette) {
	// 整块也能点，点了直接开系统的多选文件对话框。
	area := ui.ButtonBase(c).Height(132).Radius(12).Center().
		Background(p.Surface).BorderWidth(1).BorderStyle(ui.BorderDashed).
		BorderColor(p.Border).Cursor(ui.CursorPointer).
		OnClick(func() { a.addPaths(services.File().SelectFiles()) })
	hovered := area.Hovered()
	if hovered {
		area.BorderColor(ui.RGBA(124, 126, 44, 0.45))
	}
	if area.FileDragOver() {
		area.BorderColor(p.Accent).Background(ui.RGBA(124, 126, 44, 0.10))
	}

	// ButtonBase 是横向排布，里面再放一列，三个元素才能换行居中。
	area.Children(func() {
		ui.Column(c).Gap(8).Center().PassThrough().Children(func() {
			// 深色用液态玻璃；浅色下玻璃与白底糊在一起，改成灰色实底。按下由外层整块接住。
			pill := ui.Row(c).Height(34).Padding(0, 16).Radius(11).PassThrough()
			switch {
			case p.Dark:
				pill.Material(glass.Glass{})
				if hovered {
					pill.Background(glassLift(p, 0.08))
				}
			case hovered:
				pill.Background(p.Pressed)
			default:
				pill.Background(p.SurfaceHover)
			}
			pill.Children(func() {
				ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon("plus")).FontSize(14).TextColor(p.Text)
					ui.Text(c, locales.HomeUploadButton()).FontSize(12.5).FontWeight(600).
						TextColor(p.Text).SingleLine()
				})
			})

			ui.Text(c, locales.HomeDropHint()).FontSize(13.5).FontWeight(700).
				TextColor(p.Text).TextAlign(ui.Center)
			ui.Text(c, locales.HomeNoLimit()).FontSize(11.5).TextColor(p.Muted).
				TextAlign(ui.Center)
		})
	})
}

func (a *App) filesHeader(c *ui.Context, p Palette) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Grow(1).Gap(6).AlignItems(ui.Baseline).Children(func() {
			ui.Text(c, locales.HomeFileList()).FontSize(17).FontWeight(800)
			if len(a.files) > 0 {
				ui.Textf(c, "(%d)", len(a.files)).FontSize(15).FontWeight(800).TextColor(p.Accent)
			}
		})
		if a.selectMode {
			glassGroup(c, p).Children(func() {
				glassSegment(c, p, "x", a.leaveSelect)
			})
			return
		}

		glassGroup(c, p).Children(func() {
			glassSegment(c, p, "lock", a.openCaptcha)
			glassSegment(c, p, "text-plus", a.addPureText)
			if len(a.files) > 0 {
				glassSegment(c, p, "trash", a.enterSelect)
			}
		})
	})
}

func (a *App) filesList(c *ui.Context, p Palette) {
	if len(a.files) == 0 {
		ui.Column(c).Height(220).Children(func() { a.emptyList(c, p) })
		return
	}

	// 卡片本身不滚，滚动交给外层的整页滚动容器。
	ui.Column(c).FillWidth().Radius(12).Background(p.Surface).Border(1, p.Border).
		ClipInset(0, 0, 0, 0, 12).Children(func() {
		for i := range a.files {
			a.fileRow(c, p, i)
		}
	})
}

func (a *App) fileRow(c *ui.Context, p Palette, i int) {
	f := a.files[i]

	row := ui.Column(c).FillWidth()
	row.Children(func() {
		if i > 0 {
			// 分割线从文字起始处开始，像 iOS 那样不压到图标。
			ui.Row(c).Padding(0, 0, 0, 56).Children(func() {
				ui.Box(c).Height(1).FillWidth().Background(p.Border)
			})
		}

		body := ui.Row(c).Padding(9, 12).Gap(12).AlignItems(ui.Center)
		if body.Hovered() {
			body.Background(p.SurfaceHover)
		}
		// 右键菜单：在文件管理器中打开、删除这一项。
		body.ContextMenu(func(m *ui.Menu) {
			show := m.Item(locales.HomeOpenInManager()).Disabled(f.Type == "pure-text")
			if show.Chosen() {
				services.File().OpenInFileManager(f.Path)
			}
			m.Separator()
			if m.Item(locales.SelectDelete()).Chosen() {
				a.selected = map[string]bool{f.Id: true}
				a.deleteAsk = true
			}
		})

		body.OnClick(func() {
			if a.selectMode {
				a.toggleSelected(f.Id)
				return
			}
			if f.Type == "pure-text" {
				a.openTextDrawer(f)
			}
		})
		if a.selectMode && a.selected[f.Id] {
			body.Background(ui.RGBA(124, 126, 44, 0.12))
		}

		body.Children(func() {
			tile := ui.Column(c).Size(34, 34).Radius(8).Center().Background(p.SurfaceHover)
			tile.Children(func() {
				ui.Icon(c, icon(fileIcon(f))).FontSize(18).TextColor(fileTint(f))
			})

			ui.Column(c).Grow(1).Gap(3).Children(func() {
				if f.Type == "pure-text" {
					a.pureTextRow(c, p, i)
				} else {
					ui.Text(c, f.Name).FontSize(13.5).FontWeight(700).SingleLine()
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						meta := humanSize(f.Size)
						if f.Type == "folder" {
							meta = locales.HomeFolder()
						}
						ui.Text(c, meta).FontSize(11.5).TextColor(p.Muted)
						a.noteLink(c, p, i)
					})
				}
			})

			if a.selectMode {
				selectBox(c, p, a.selected[f.Id])
			}
		})
	})
}

// emptyList 是清单为空时的占位：居中的图标与说明，本身也接得住拖进来的文件。
func (a *App) emptyList(c *ui.Context, p Palette) {
	box := ui.Column(c).Grow(1).Center().Padding(0, 0, 40, 0)

	if paths := box.DroppedFiles(); len(paths) > 0 {
		a.addPaths(paths)
	}
	over := box.FileDragOver()
	tint := p.Muted
	if over {
		tint = p.Accent
	}

	box.Children(func() {
		tile := ui.Column(c).Size(46, 46).Radius(12).Center().Background(p.SurfaceHover)
		if over {
			tile.Background(ui.RGBA(124, 126, 44, 0.18))
		}
		tile.Children(func() {
			ui.Icon(c, icon("files-off")).FontSize(21).TextColor(tint)
		})

		spacer(c, 12)
		ui.Text(c, locales.HomeEmptyTitle()).FontSize(13).FontWeight(700).
			TextColor(p.Subtle)
		spacer(c, 5)
		ui.Text(c, locales.HomeEmptyHint()).FontSize(11.5).TextColor(p.Muted).
			TextAlign(ui.Center).MaxWidth(220)
	})
}

func (a *App) pureTextRow(c *ui.Context, p Palette, i int) {
	f := a.files[i]

	name := f.Name
	if name == "" {
		name = locales.HomePureText()
	}
	ui.Text(c, name).FontSize(13.5).FontWeight(700).SingleLine()

	preview := f.Text
	if preview == "" {
		preview = locales.HomeTextNoContent()
	}
	link := ui.ButtonBase(c).Padding(0).Background(ui.Transparent).
		Justify(ui.Start).Gap(4).Cursor(ui.CursorPointer).
		OnClick(func() { a.openTextDrawer(f) })
	link.Children(func() {
		ui.Icon(c, icon("edit")).FontSize(13).TextColor(p.Accent)
		ui.Text(c, preview).FontSize(11.5).TextColor(p.Accent).SingleLine()
	})
}

func (a *App) noteLink(c *ui.Context, p Palette, i int) {
	f := a.files[i]
	label := f.Note
	if label == "" {
		label = locales.HomeNoteAdd()
	}
	link := ui.ButtonBase(c).Padding(0).Background(ui.Transparent).
		Justify(ui.Start).Gap(4).Cursor(ui.CursorPointer).OnClick(func() {
		a.noteID, a.noteDraft, a.noteOpen = f.Id, f.Note, true
	})
	link.Children(func() {
		ui.Icon(c, icon("message-plus")).FontSize(12).TextColor(p.Accent)
		ui.Text(c, label).FontSize(11.5).TextColor(p.Accent).SingleLine()
	})
}

func (a *App) shareButton(c *ui.Context, p Palette) {
	label := locales.HomeShareStart()
	if a.share {
		label = locales.HomeShareStop()
	}

	box := ui.Row(c).Size(56, 56).Center()
	box.Children(func() {
		b := ui.ButtonBase(c).Size(56, 56).Padding(0).Radius(999).
			Background(p.Accent).Cursor(ui.CursorPointer).
			OnClick(a.toggleShare).Label(label)
		b.Shadow(0, 10, 28, -4, ui.RGBA(124, 126, 44, 0.55))
		if b.Hovered() {
			b.Background(ui.Hex("#8b8d38"))
		}

		// 带过冲的弹簧值：0 未分享，1 分享中。
		t := b.AnimateTo("share", boolFloat(a.share), ui.Animation{
			Spring:    ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.28},
			NonMotion: true,
		})
		on := clamp01(t)

		// 分享中：两圈错开的呼吸光晕，先实心后描边，扩张到位再淡出。
		if on > 0.02 {
			halo := box.Loop("share-halo", 1600*time.Millisecond, ui.Linear)
			for i := 0; i < 2; i++ {
				offset := float32(0)
				peak := float32(0.35)
				if i == 1 {
					offset, peak = 0.3125, 0.5
				}
				u := frac(halo + offset)
				opacity, scale := pulse(u, peak)
				ring := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
					Radius(999).PassThrough()
				if i == 0 {
					ring.Background(p.Accent)
				} else {
					ring.Border(1.5, p.Accent)
				}
				ring.Opacity(opacity*on).Scale(scale, scale)
			}
		}

		// 图标交叉切换：一个转出去，另一个从反方向转进来。
		b.Children(func() {
			ui.Box(c).Size(26, 26).Children(func() {
				swapIcon(c, p, "broadcast", 1-on, -90*on, 1-0.4*on)
				swapIcon(c, p, "player-stop", on, 90*(1-on), 0.6+0.4*on)
			})
		})
	})
}

// swapIcon 是叠在同一个位置上的图标，用透明度、角度和缩放做切换。
func swapIcon(c *ui.Context, p Palette, name string, opacity, rotate, scale float32) {
	if opacity <= 0.01 {
		return
	}
	layer := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center()
	layer.Opacity(opacity).Scale(scale, scale)
	layer.Children(func() {
		ui.Icon(c, icon(name)).FontSize(24).TextColor(p.AccentText).Rotate(rotate)
	})
}

// pulse 把 0..1 的循环时间换成光晕的透明度与缩放：前 45% 扩张并浮现，之后停在原处淡出。
func pulse(u, peak float32) (opacity, scale float32) {
	if u <= 0.45 {
		k := u / 0.45
		return peak * k, 1 + 0.22*k
	}
	return peak * (1 - (u-0.45)/0.55), 1.22
}

func frac(v float32) float32 {
	v -= float32(int(v))
	if v < 0 {
		v++
	}
	return v
}

func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func boolFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func (a *App) noteDrawer(c *ui.Context, p Palette) {
	a.drawer(c, p, &a.noteOpen, drawerSpec{
		icon:  "message-plus",
		title: locales.HomeNoteTitle(),
		body: func() {
			labeledField(c, p, locales.HomeNoteLabel(), a.noteInput.FocusWithin(c), func() {
				ui.TextInput(c.Key("note"), &a.noteDraft).Bind(&a.noteInput).
					FillWidth().FontSize(13).Padding(0).
					Background(ui.Transparent).BorderWidth(0).
					Placeholder(locales.HomeNotePlaceholder()).AutoFocus().
					OnSubmit(func() { a.saveNote() })
			})
		},
		footer: func() {
			ui.Spacer(c).Grow(1)
			dialogButton(c, p, locales.DialogCancel(), func() { a.noteOpen = false })
			primaryDialogButton(c, p, locales.DialogConfirm(), a.saveNote)
		},
	})
}

func (a *App) pureTextDrawer(c *ui.Context, p Palette) {
	a.drawer(c, p, &a.textOpen, drawerSpec{
		header: func() {
			ui.Text(c, locales.HomeEditText()).FontSize(11).FontWeight(700).
				LetterSpacing(0.4).TextColor(p.Muted)
			spacer(c, 6)

			// 标题就地可改：输入框长得像标题，右边一个笔形图标提示可编辑。
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.TextInput(c.Key("text-name"), &a.textName).Bind(&a.textTitle).
					Grow(1).FontSize(21).FontWeight(800).Padding(0).
					Background(ui.Transparent).BorderWidth(0).
					Placeholder(locales.HomeTextInputTitle()).AutoFocus()
				ui.Icon(c, icon("edit")).FontSize(14).TextColor(p.Muted)
			})
		},
		body: func() {
			labeledField(c, p, locales.HomeTextContentLabel(), a.textContent.FocusWithin(c), func() {
				ui.TextArea(c.Key("text-body"), &a.textDraft).Bind(&a.textContent).
					Lines(6, 14).FillWidth().FontSize(13).Padding(0).
					Background(ui.Transparent).BorderWidth(0).
					Placeholder(locales.HomeEnterContent())
			})
		},
		footer: func() {
			ui.Spacer(c).Grow(1)
			dialogButton(c, p, locales.DialogCancel(), func() { a.textOpen = false })
			primaryDialogButton(c, p, locales.DialogConfirm(), a.saveText)
		},
	})
}

// openTextDrawer 把条目的标题与内容读进草稿，改完点确认才写回。
func (a *App) openTextDrawer(f services.FileInfo) {
	a.textID, a.textNew, a.textName, a.textDraft = f.Id, false, f.Name, f.Text
	a.textOpen = true
}

// addPureText 只是打开一个空白抽屉：确认时才真正落成条目，取消就什么都不留。
func (a *App) addPureText() {
	a.textID, a.textNew, a.textName, a.textDraft = "", true, "", ""
	a.textOpen = true
}

// addPaths 把选中的路径变成分享项，去重后写回存储。
func (a *App) addPaths(paths []string) {
	added := false
	for _, path := range paths {
		if a.hasPath(path) {
			continue
		}
		file := services.File().GetFile(path)
		if file == nil {
			continue
		}
		a.files = append(a.files, *file)
		added = true
	}
	if added {
		services.Storage().Set("share-list", a.files)
	}
}

func (a *App) hasPath(path string) bool {
	for _, f := range a.files {
		if f.Path == path {
			return true
		}
	}
	return false
}

func (a *App) saveNote() {
	for i := range a.files {
		if a.files[i].Id == a.noteID {
			a.files[i].Note = a.noteDraft
		}
	}
	a.noteOpen = false
	services.Storage().Set("share-list", a.files)
}

func (a *App) saveText() {
	if a.textNew && a.textName == "" && a.textDraft == "" {
		// 什么都没填就点确认，等于取消。
		a.textNew, a.textOpen = false, false
		return
	}
	if a.textNew {
		id := uuid.NewString()
		a.files = append(a.files, services.FileInfo{
			Id: id, Path: id, Type: "pure-text", Name: a.textName, Text: a.textDraft,
		})
		a.textNew = false
	} else {
		for i := range a.files {
			if a.files[i].Id == a.textID {
				a.files[i].Name = a.textName
				a.files[i].Text = a.textDraft
			}
		}
	}
	a.textOpen = false
	services.Storage().Set("share-list", a.files)
}

func (a *App) applyCaptcha() {
	a.captcha = joincode(a.captchaDigits)
	services.Share().SetCaptcha(a.captcha)
}

func (a *App) toggleShare() {
	a.share = !a.share
	services.Storage().Set("is-share", a.share)
	if a.share {
		services.Storage().Set("share-list", a.files)
		app.Gin().Start(":" + a.port)
		return
	}
	app.Gin().Stop()
}

func fileIcon(f services.FileInfo) string {
	if f.Type == "folder" {
		return "folder"
	}
	if f.Type == "pure-text" {
		return "keyboard"
	}
	if name, ok := fileIcons[f.Type]; ok {
		return name
	}
	return "file"
}

// fileTint 给每种文件类型一块颜色，图标用白色压在上面，像 iOS 设置里的图标。
func fileTint(f services.FileInfo) ui.Color {
	switch f.Type {
	case "folder":
		return ui.Hex("#3b82f6")
	case "pure-text":
		return ui.Hex("#7c7e2c")
	case "jpg", "jpeg", "png", "gif", "svg", "bmp":
		return ui.Hex("#0ea5e9")
	case "pdf":
		return ui.Hex("#ef4444")
	case "doc", "docx":
		return ui.Hex("#2563eb")
	case "xls", "xlsx":
		return ui.Hex("#16a34a")
	case "ppt", "pptx":
		return ui.Hex("#ea580c")
	case "zip", "rar", "7z", "tar", "gz", "iso":
		return ui.Hex("#d97706")
	case "apk":
		return ui.Hex("#22c55e")
	case "ipa":
		return ui.Hex("#64748b")
	case "ts", "js", "tsx", "jsx", "html", "css", "scss", "json", "yaml", "toml":
		return ui.Hex("#8b5cf6")
	case "md", "txt":
		return ui.Hex("#64748b")
	}
	return ui.Hex("#6b7280")
}

// fileIcons 对齐 Web 端 Uploader/file-type.ts 的映射。
var fileIcons = map[string]string{
	"jpg": "file-type-jpg", "jpeg": "file-type-jpg", "png": "file-type-png",
	"gif": "photo-bolt", "svg": "file-type-svg", "bmp": "file-type-bmp",
	"pdf": "file-type-pdf", "doc": "file-type-doc", "docx": "file-type-doc",
	"xls": "file-type-xls", "xlsx": "file-type-xls",
	"ppt": "file-type-ppt", "pptx": "file-type-ppt",
	"txt": "file-type-txt", "zip": "file-type-zip",
	"rar": "file-zip", "7z": "file-zip", "tar": "file-zip", "gz": "file-zip",
	"ts": "file-type-ts", "js": "file-type-js", "html": "file-type-html",
	"css": "file-type-css", "scss": "file-type-css", "json": "file-dots",
	"yaml": "file-dots", "toml": "file-dots", "md": "markdown",
	"apk": "brand-android", "ipa": "brand-apple",
}

func humanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size)
	for _, u := range units {
		value /= unit
		if value < unit {
			return fmt.Sprintf("%.2f%s", value, u)
		}
	}
	return fmt.Sprintf("%.2fPB", value/unit)
}
