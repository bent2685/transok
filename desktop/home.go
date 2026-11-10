package desktop

import (
	"fmt"


	"transok/backend/app"
	"transok/backend/services"

	"github.com/egoist/mygo/ui"
)

func (a *App) home(c *ui.Context) {
	ui.Row(c).AlignItems(ui.End).Children(func() {
		ui.Column(c).Grow(1).Gap(2).Children(func() {
			ui.Text(c, a.appInfo["name"]).FontSize(26).Bold()
			ui.Text(c, "局域网文件分享").FontSize(12).TextColor(c.Theme().TextMuted)
		})
		ui.Row(c).Gap(6).Children(func() {
			a.iconButton(c, "device-desktop-search", func() { a.page = discoverRoute })
			a.iconButton(c, "lock", func() {
				a.captchaDraft = a.captcha
				a.captchaOpen = true
			})
			a.iconButton(c, "settings", func() { a.page = settingsRoute })
		})
	})

	a.dropArea(c)
	a.fileList(c)

	label := "开始分享"
	if a.share {
		label = "停止分享"
	}
	ui.PrimaryButton(c, label).Height(44).Radius(10).OnClick(a.toggleShare)

	if a.captchaOpen {
		a.captchaDialog(c)
	}
}

func (a *App) iconButton(c *ui.Context, name string, onClick func()) {
	ui.Button(c, "").Size(32, 32).Radius(6).OnClick(onClick).Children(func() {
		ui.Icon(c, icon(name)).FontSize(16)
	})
}

// dropArea 接住从系统拖进来的文件与文件夹，也提供选择入口。
func (a *App) dropArea(c *ui.Context) {
	area := ui.Column(c).Padding(16, 12).Gap(10).Radius(10).Center().
		BorderWidth(1).BorderStyle(ui.BorderDashed).BorderColor(c.Theme().Border)

	if paths := area.DroppedFiles(); len(paths) > 0 {
		a.addPaths(paths)
	}
	if area.FileDragOver() {
		area.BorderColor(olive).Background(ui.RGBA(124, 126, 44, 0.08))
	}

	area.Children(func() {
		ui.Icon(c, icon("upload")).FontSize(24).TextColor(olive)
		ui.Text(c, "拖入文件或文件夹").FontSize(13).TextColor(c.Theme().TextMuted)

		button := ui.PrimaryButton(c, "上传").Height(32)
		ui.Popover(c, button, &a.uploadMenu, func() {
			ui.Button(c, "选择文件").OnClick(func() {
				a.uploadMenu = false
				a.addPaths(services.File().SelectFiles())
			})
			ui.Button(c, "选择文件夹").OnClick(func() {
				a.uploadMenu = false
				a.addPaths(services.File().SelectFolder())
			})
		})
	})
}

func (a *App) fileList(c *ui.Context) {
	a.list.Key = func(i int) any { return a.files[i].Id }

	ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
		ui.Text(c, "分享清单").FontSize(15).Bold()
		if len(a.files) > 0 {
			ui.Textf(c, "(%d)", len(a.files)).FontSize(13).TextColor(olive)
		}
	})

	ui.List(c, &a.list, len(a.files), func(i int) {
		a.fileRow(c, i)
	}).Grow(1).Children(func() {
		if len(a.files) == 0 {
			ui.Column(c).Padding(0, 0, 24, 0).Gap(6).Center().Children(func() {
				ui.Icon(c, icon("files-off")).FontSize(26).TextColor(c.Theme().TextMuted)
				ui.Text(c, "还没有内容").FontSize(13).TextColor(c.Theme().TextMuted)
			})
		}
	})
}

func (a *App) fileRow(c *ui.Context, i int) {
	f := a.files[i]
	ui.Row(c).Padding(8).Gap(10).Radius(8).AlignItems(ui.Center).
		Background(c.Theme().Surface).Children(func() {
		ui.Column(c).Size(32, 32).Radius(6).Center().Background(c.Theme().SurfaceHover).Children(func() {
			ui.Icon(c, icon(fileIcon(f.Type))).FontSize(17)
		})
		ui.Column(c).Grow(1).Gap(2).Children(func() {
			ui.Text(c, f.Name).FontSize(13).Bold().SingleLine()
			meta := fileMeta(f)
			if f.Note != "" {
				meta += " · " + f.Note
			}
			ui.Text(c, meta).FontSize(11).TextColor(c.Theme().TextMuted).SingleLine()
		})
		if f.Type == "folder" {
			ui.Button(c, "").Size(28, 28).Radius(6).OnClick(func() {
				services.File().OpenInFileManager(f.Path)
			}).Children(func() {
				ui.Icon(c, icon("external-link")).FontSize(15)
			})
		}
		ui.Button(c, "").Size(28, 28).Radius(6).OnClick(func() { a.removeFile(f.Id) }).Children(func() {
			ui.Icon(c, icon("trash")).FontSize(15)
		})
	})
}

func (a *App) captchaDialog(c *ui.Context) {
	ui.Modal(c, &a.captchaOpen, func() {
		ui.Column(c).Width(280).Padding(18).Gap(12).Radius(12).
			Background(c.Theme().Surface).Children(func() {
			ui.Text(c, "设置口令").FontSize(15).Bold()
			ui.Text(c, "被分享方需要用这个口令才能下载").FontSize(12).TextColor(c.Theme().TextMuted)
			ui.TextInput(c.Key("captcha"), &a.captchaDraft).MaxLength(5).AutoFocus().
				OnSubmit(func() { a.captchaOpen = false; a.applyCaptcha() })
			ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
				ui.Button(c, "取消").OnClick(func() { a.captchaOpen = false })
				ui.PrimaryButton(c, "保存").OnClick(func() {
					a.captchaOpen = false
					a.applyCaptcha()
				})
			})
		})
	})
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

func (a *App) removeFile(id string) {
	kept := make([]services.FileInfo, 0, len(a.files))
	for _, f := range a.files {
		if f.Id != id {
			kept = append(kept, f)
		}
	}
	a.files = kept
	services.Storage().Set("share-list", a.files)
}

func (a *App) applyCaptcha() {
	a.captcha = a.captchaDraft
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

func fileIcon(fileType string) string {
	switch fileType {
	case "folder":
		return "folder"
	case "pure-text":
		return "keyboard"
	}
	return "file"
}

func fileMeta(f services.FileInfo) string {
	switch f.Type {
	case "folder":
		return "文件夹"
	case "pure-text":
		if f.Text == "" {
			return "纯文本 · 无内容"
		}
		return "纯文本"
	}
	return humanSize(f.Size)
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
			return fmt.Sprintf("%.1f %s", value, u)
		}
	}
	return fmt.Sprintf("%.1f PB", value/unit)
}
