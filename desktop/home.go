package desktop

import (
	"fmt"

	"transok/backend/app"
	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo/ui"
)

func (a *App) home(c *ui.Context) {
	p := a.palette()

	ui.Column(c).Fill().Padding(a.topInset(c), 16, 0, 16).Children(func() {
		a.homeHeader(c, p)
		spacer(c, 12)
		a.dropArea(c, p)
		spacer(c, 12)
		a.limitLine(c, p)
		spacer(c, 16)
		a.filesHeader(c, p)
		spacer(c, 14)
		a.filesList(c, p)
	})

	ui.Overlay(c, func() {
		ui.Row(c).Absolute().Left(0).Right(0).Bottom(26).Justify(ui.Center).Children(func() {
			a.shareButton(c, p)
		})
	})

	if a.captchaOpen {
		a.captchaSheet(c, p)
	}
	if a.noteOpen {
		a.noteDialog(c, p)
	}
	if a.textOpen {
		a.pureTextDialog(c, p)
	}
}

// topInset 让内容避开窗口左上角的红黄绿按钮。
func (a *App) topInset(c *ui.Context) float32 {
	top := c.TitleBar().Height + 6
	if top < 38 {
		top = 38
	}
	return top
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
		ui.Row(c).Gap(6).Children(func() {
			iconButton(c, p, "device-desktop-search", func() { a.page = discoverRoute })
			iconButton(c, p, "lock", func() {
				a.captchaDigits = make([]string, captchaLength)
				for i := 0; i < len(a.captcha) && i < captchaLength; i++ {
					a.captchaDigits[i] = string(a.captcha[i])
				}
				a.focusReset = true
				a.captchaOpen = true
			})
			iconButton(c, p, "settings", func() { a.page = settingsRoute })
		})
	})
}

// dropArea 接住从系统拖进来的文件与文件夹，也提供选择入口。
func (a *App) dropArea(c *ui.Context, p Palette) {
	area := ui.Column(c).Height(180).Radius(10).Center().Gap(12).
		Background(p.Surface).BorderWidth(1).BorderStyle(ui.BorderDashed).BorderColor(p.Border)

	if paths := area.DroppedFiles(); len(paths) > 0 {
		a.addPaths(paths)
	}
	if area.FileDragOver() {
		area.BorderColor(p.Accent).Background(ui.RGBA(124, 126, 44, 0.10))
	}

	area.Children(func() {
		tile := ui.Column(c).Size(48, 48).Radius(12).Center().Background(p.Accent)
		tile.Shadow(0, 6, 18, 0, ui.RGBA(124, 126, 44, 0.45))
		tile.Children(func() {
			ui.Icon(c, icon("upload")).FontSize(26).TextColor(p.AccentText)
		})

		ui.Text(c, locales.HomeDropHint()).FontSize(13.5).TextColor(p.Subtle)

		button := ui.ButtonBase(c).Height(38).Padding(0, 20).Radius(8).
			Background(p.Accent).Cursor(ui.CursorPointer).OnClick(func() {
			a.uploadMenu = !a.uploadMenu
		})
		if button.Hovered() {
			button.Background(ui.Hex("#8b8d38"))
		}
		button.Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon("plus")).FontSize(16).TextColor(p.AccentText)
				ui.Text(c, locales.HomeUploadButton()).FontSize(13).FontWeight(600).
					TextColor(p.AccentText).SingleLine()
			})
		})

		ui.Popover(c, button, &a.uploadMenu, func() {
			ui.Button(c, locales.HomeSelectFile()).OnClick(func() {
				a.uploadMenu = false
				a.addPaths(services.File().SelectFiles())
			})
			ui.Button(c, locales.HomeSelectFolder()).OnClick(func() {
				a.uploadMenu = false
				a.addPaths(services.File().SelectFolder())
			})
		})
	})
}

func (a *App) limitLine(c *ui.Context, p Palette) {
	ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Size(4, 4).Radius(999).Background(p.Accent)
		ui.Text(c, locales.HomeNoLimit()).FontSize(11.5).TextColor(p.Muted)
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
		ui.Row(c).Gap(6).Children(func() {
			chip(c, p, "keyboard", locales.HomePureText(), a.addPureText)
			chip(c, p, "trash", "", a.clearFiles)
		})
	})
}

func (a *App) filesList(c *ui.Context, p Palette) {
	a.list.Key = func(i int) any { return a.files[i].Id }

	ui.List(c, &a.list, len(a.files), func(i int) {
		a.fileRow(c, p, i)
	}).Grow(1).Gap(8).Padding(0, 0, 90, 0).Children(func() {
		ui.Column(c).Padding(0, 0, 40, 0).Center().Children(func() {
			ui.Text(c, locales.HomeNoFiles()).FontSize(12.5).TextColor(p.Muted)
		})
	})
}

func (a *App) fileRow(c *ui.Context, p Palette, i int) {
	f := a.files[i]

	row := ui.Row(c).Padding(10, 11).Gap(10).Radius(10).AlignItems(ui.Center).
		Background(p.Surface).Border(1, p.Border)
	if row.Hovered() {
		row.Border(1, ui.RGBA(124, 126, 44, 0.35))
	}

	row.Children(func() {
		ui.Column(c).Size(36, 36).Radius(8).Center().Background(p.SurfaceHover).Children(func() {
			ui.Icon(c, icon(fileIcon(f))).FontSize(18).TextColor(p.Text)
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
}

func (a *App) pureTextRow(c *ui.Context, p Palette, i int) {
	f := &a.files[i]

	name := ui.TextInput(c.Key("text-name-"+f.Id), &f.Name).FontSize(13.5).
		FontWeight(700).Padding(0).Background(ui.Transparent).BorderWidth(0)
	if name.Changed() {
		services.Storage().Set("share-list", a.files)
	}

	ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
		preview := f.Text
		if preview == "" {
			preview = locales.HomeTextNoContent()
		}
		link := ui.Row(c).Gap(4).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(func() {
			a.textID, a.textDraft, a.textOpen = f.Id, f.Text, true
		})
		link.Children(func() {
			ui.Icon(c, icon("edit")).FontSize(13).TextColor(p.Accent)
			ui.Text(c, preview).FontSize(11.5).TextColor(p.Accent).SingleLine()
		})
	})
}

func (a *App) noteLink(c *ui.Context, p Palette, i int) {
	f := a.files[i]
	label := f.Note
	if label == "" {
		label = locales.HomeNoteAdd()
	}
	link := ui.Row(c).Gap(4).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(func() {
		a.noteID, a.noteDraft, a.noteOpen = f.Id, f.Note, true
	})
	link.Children(func() {
		ui.Icon(c, icon("message-plus")).FontSize(12).TextColor(p.Accent)
		ui.Text(c, label).FontSize(11.5).TextColor(p.Accent).SingleLine()
	})
}

func (a *App) shareButton(c *ui.Context, p Palette) {
	name, label := "broadcast", locales.HomeShareStart()
	if a.share {
		name, label = "player-stop", locales.HomeShareStop()
	}

	b := ui.ButtonBase(c).Size(56, 56).Padding(0).Radius(999).
		Background(p.Accent).Cursor(ui.CursorPointer).OnClick(a.toggleShare).Label(label)
	b.Shadow(0, 10, 28, -4, ui.RGBA(124, 126, 44, 0.55))
	if b.Hovered() {
		b.Background(ui.Hex("#8b8d38"))
	}
	b.Children(func() {
		ui.Icon(c, icon(name)).FontSize(24).TextColor(p.AccentText)
	})
}

func (a *App) noteDialog(c *ui.Context, p Palette) {
	ui.Modal(c, &a.noteOpen, func() {
		ui.Column(c).Width(300).Padding(18).Gap(10).Radius(14).
			Background(p.Surface).Border(1, p.Border).Children(func() {
			ui.Text(c, locales.HomeNoteTitle()).FontSize(15).FontWeight(700)
			ui.TextInput(c.Key("note"), &a.noteDraft).
				Placeholder(locales.HomeNotePlaceholder()).AutoFocus().
				OnSubmit(func() { a.saveNote() })
			ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
				dialogButton(c, p, locales.DialogCancel(), func() { a.noteOpen = false })
				primaryDialogButton(c, p, locales.DialogConfirm(), a.saveNote)
			})
		})
	})
}

func (a *App) pureTextDialog(c *ui.Context, p Palette) {
	ui.Modal(c, &a.textOpen, func() {
		ui.Column(c).Width(300).Padding(18).Gap(10).Radius(14).
			Background(p.Surface).Border(1, p.Border).Children(func() {
			ui.Text(c, locales.HomeEditText()).FontSize(15).FontWeight(700)
			ui.TextInput(c.Key("pure-text"), &a.textDraft).Lines(3, 6).
				Height(110).Placeholder(locales.HomeEnterContent()).AutoFocus()
			ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
				dialogButton(c, p, locales.DialogCancel(), func() { a.textOpen = false })
				primaryDialogButton(c, p, locales.DialogConfirm(), a.saveText)
			})
		})
	})
}

func (a *App) addPureText() {
	id := fmt.Sprintf("text-%d", len(a.files)+1)
	a.files = append(a.files, services.FileInfo{
		Id: id, Name: locales.HomeTextInputTitle(), Path: id, Type: "pure-text",
	})
	services.Storage().Set("share-list", a.files)
	a.textID, a.textDraft, a.textOpen = id, "", true
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
	for i := range a.files {
		if a.files[i].Id == a.textID {
			a.files[i].Text = a.textDraft
			if a.files[i].Name == locales.HomeTextInputTitle() {
				a.files[i].Name = locales.HomeTextNoContent()
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
