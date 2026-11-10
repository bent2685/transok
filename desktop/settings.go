package desktop

import (
	"transok/backend/services"

	"github.com/egoist/mygo/ui"
)

func (a *App) settings(c *ui.Context) {
	a.header(c, "设置")

	ui.Column(c).Grow(1).Gap(14).Children(func() {
		ui.Text(c, "通用").FontSize(12).Bold().TextColor(c.Theme().TextMuted)
		ui.Column(c).Padding(12).Gap(10).Radius(10).Background(c.Theme().Surface).Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "端口").Grow(1).FontSize(13)
				ui.TextInput(c.Key("port"), &a.port).Width(96).FontSize(13).
					OnSubmit(func() { services.Storage().Set("port", a.port) })
			})
			ui.Text(c, "修改端口后需要重新开始分享").FontSize(11).TextColor(c.Theme().TextMuted)
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "设备名称").Grow(1).FontSize(13)
				ui.Text(c, a.deviceName()).FontSize(13).TextColor(c.Theme().TextMuted)
			})
		})

		ui.Text(c, "关于").FontSize(12).Bold().TextColor(c.Theme().TextMuted)
		ui.Column(c).Padding(12).Gap(8).Radius(10).Background(c.Theme().Surface).Children(func() {
			a.aboutRow(c, "应用", a.appInfo["name"])
			a.aboutRow(c, "版本", a.appInfo["version"])
			a.aboutRow(c, "作者", a.appInfo["author"])
			a.aboutRow(c, "邮箱", a.appInfo["email"])
		})
	})
}

func (a *App) aboutRow(c *ui.Context, label, value string) {
	ui.Row(c).Gap(8).Children(func() {
		ui.Text(c, label).Grow(1).FontSize(13)
		ui.Text(c, value).FontSize(13).TextColor(c.Theme().TextMuted)
	})
}

func (a *App) deviceName() string {
	if uname, ok := services.Storage().Get("uname"); ok {
		return uname.(string)
	}
	return "transok"
}
