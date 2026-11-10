package desktop

import (
	"fmt"
	"time"

	"transok/backend/app"
	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo/ui"
	"github.com/google/uuid"
)

func (a *App) home(c *ui.Context) {
	p := a.palette()

	ui.Column(c).Fill().Padding(0, 16, 0, 16).Children(func() {
		a.homeHeader(c, p)
		spacer(c, 14)
		a.dropArea(c, p)
		spacer(c, 18)
		a.filesHeader(c, p)
		spacer(c, 14)
		a.filesList(c, p)
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
			a.shareButton(c, p)
		})
	})

	// 抽屉常驻构建，关闭时靠 Exit 动画退场。
	a.captchaDrawer(c, p)
	// 弹窗常驻构建，关闭时靠 Exit 动画退场。
	a.noteDrawer(c, p)
	a.pureTextDrawer(c, p)
}

func (a *App) homeHeader(c *ui.Context, p Palette) {
	header := ui.Row(c).AlignItems(ui.End).DragWindow()
	header.Children(func() {
		ui.Column(c).Grow(1).Gap(3).Children(func() {
			ui.Row(c).AlignItems(ui.End).Gap(4).Children(func() {
				ui.Text(c, a.appInfo["name"]).FontSize(27).FontWeight(800).LetterSpacing(-1.2)
				ui.Box(c).Size(5, 5).Radius(1).Background(p.Accent).Margin(0, 0, 5, 0)
			})
			ui.Text(c, locales.HomeTitle()).FontSize(12.5).TextColor(p.Subtle)
		})
		glassGroup(c, p).Children(func() {
			glassSegment(c, p, "device-desktop-search", func() { a.page = discoverRoute })
			glassSegment(c, p, "lock", func() {
				a.captchaDigits = make([]string, captchaLength)
				for i := 0; i < len(a.captcha) && i < captchaLength; i++ {
					a.captchaDigits[i] = string(a.captcha[i])
				}
				a.focusReset = true
				a.captchaOpen = true
			})
			glassSegment(c, p, "settings", func() { a.page = settingsRoute })
		})
	})
}

// dropArea 接住从系统拖进来的文件与文件夹，也提供选择入口。
func (a *App) dropArea(c *ui.Context, p Palette) {
	// 整块都能点，点了直接开系统的多选文件对话框；文件夹拖进来即可。
	area := ui.ButtonBase(c).Height(96).Padding(0, 16).Radius(12).Gap(12).
		Background(p.Surface).BorderWidth(1).BorderStyle(ui.BorderDashed).
		BorderColor(p.Border).Cursor(ui.CursorPointer).
		OnClick(func() { a.addPaths(services.File().SelectFiles()) })
	if area.Hovered() {
		area.BorderColor(ui.RGBA(124, 126, 44, 0.45))
	}

	if paths := area.DroppedFiles(); len(paths) > 0 {
		a.addPaths(paths)
	}
	if area.FileDragOver() {
		area.BorderColor(p.Accent).Background(ui.RGBA(124, 126, 44, 0.10))
	}

	area.Children(func() {
		tile := ui.Column(c).Size(40, 40).Radius(11).Center().
			Background(ui.RGBA(124, 126, 44, 0.16))
		tile.Children(func() {
			ui.Icon(c, icon("upload")).FontSize(20).TextColor(p.Accent)
		})

		ui.Column(c).Grow(1).Gap(3).Children(func() {
			ui.Text(c, locales.HomeDropHint()).FontSize(13).FontWeight(700).
				TextColor(p.Text).SingleLine()
			ui.Text(c, locales.HomeNoLimit()).FontSize(11).TextColor(p.Muted).SingleLine()
		})

		ui.Icon(c, icon("chevron-right")).FontSize(15).TextColor(p.Muted)
	})
}

func (a *App) filesHeader(c *ui.Context, p Palette) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Grow(1).Gap(6).AlignItems(ui.End).Children(func() {
			ui.Text(c, locales.HomeFileList()).FontSize(17).FontWeight(800)
			if len(a.files) > 0 {
				ui.Textf(c, "(%d)", len(a.files)).FontSize(15).FontWeight(800).TextColor(p.Accent)
			}
		})
		glassGroup(c, p).Children(func() {
			glassChip(c, p, "keyboard", locales.HomePureText(), a.addPureText)
			if len(a.files) > 0 {
				glassChip(c, p, "trash", "", a.clearFiles)
			}
		})
	})
}

func (a *App) filesList(c *ui.Context, p Palette) {
	if len(a.files) == 0 {
		ui.Column(c).Grow(1).Children(func() { a.emptyList(c, p) })
		return
	}

	// 整张卡片一起滚，不把滚动限制在卡片内部；卡片下面留出的空当比悬浮按钮高。
	ui.Scroll(c).Grow(1).Scrollbars(false).Children(func() {
		ui.Column(c).FillWidth().Radius(12).Background(p.Surface).Border(1, p.Border).
			ClipInset(0, 0, 0, 0, 12).Children(func() {
			for i := range a.files {
				a.fileRow(c, p, i)
			}
		})
		spacer(c, 130)
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
		body.OnClick(func() {
			if f.Type == "pure-text" {
				a.openTextDrawer(f)
			}
		})
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

			if f.Type == "folder" {
				ghostButton(c, p, "external-link", func() {
					services.File().OpenInFileManager(f.Path)
				})
			}
			ghostButton(c, p, "trash", func() { a.removeFile(f.Id) })
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

func (a *App) clearFiles() {
	a.files = nil
	a.share = false
	services.Storage().Set("share-list", a.files)
	services.Storage().Set("is-share", false)
	app.Gin().Stop()
}

func (a *App) removeFile(id string) {
	kept := make([]services.FileInfo, 0, len(a.files))
	for _, f := range a.files {
		if f.Id != id {
			kept = append(kept, f)
		}
	}
	a.files = kept
	if len(kept) == 0 {
		a.share = false
		services.Storage().Set("is-share", false)
		app.Gin().Stop()
	}
	services.Storage().Set("share-list", a.files)
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
