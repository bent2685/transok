package desktop

import (
	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func (a *App) settings(c *ui.Context) {
	p := a.palette()

	ui.Column(c).Fill().Padding(0, 16, 16, 16).Children(func() {
		dragHeader(c, p, locales.SettingsTitle(), func() { a.page = homeRoute })

		spacer(c, 22)
		sectionLabel(c, p, locales.SettingsCommon())
		spacer(c, 8)

		card(c, p).Children(func() {
			a.languageRow(c, p)
			divider(c, p)
			a.portRow(c, p)
			divider(c, p)
			a.deviceRow(c, p)
		})

		spacer(c, 22)
		sectionLabel(c, p, locales.SettingsAbout())
		spacer(c, 8)

		card(c, p).Children(func() {
			a.valueRow(c, p, "circle-dot", locales.SettingsAppName(), a.appInfo["name"])
			divider(c, p)
			a.valueRow(c, p, "info-circle", locales.SettingsVersion(), a.appInfo["version"])
			divider(c, p)
			a.valueRow(c, p, "user-circle", locales.SettingsAuthor(), a.appInfo["author"])
			divider(c, p)
			a.valueRow(c, p, "mail", locales.SettingsEmail(), a.appInfo["email"])
			divider(c, p)

			g := ui.Row(c).Gap(12).AlignItems(ui.Center).Padding(11, 14).Height(48).
				Cursor(ui.CursorPointer).OnClick(func() {
				mygo.Shell.OpenExternal("https://github.com/bent2685/transok")
			})
			if g.Hovered() {
				g.Background(glassLift(p, 0.06))
			}
			g.Children(func() {
				ui.Icon(c, icon("brand-github")).FontSize(17).TextColor(p.Text)
				ui.Text(c, "Github").Grow(1).FontSize(13.5).FontWeight(700)
				ui.Text(c, "transok").FontSize(13.5).FontWeight(700)
			})
		})

		spacer(c, 18)
		ui.Row(c).Justify(ui.Center).Children(func() {
			ui.Text(c, locales.SettingsFooter()).FontSize(11.5).TextColor(p.Muted)
		})
	})
}

// dragHeader 是详情页的返回键 + 大标题，兼作窗口拖动区。
func dragHeader(c *ui.Context, p Palette, title string, back func()) {
	header := ui.Row(c).Gap(12).AlignItems(ui.Center).DragWindow()
	header.Children(func() {
		circleButton(c, p, "arrow-left", back)
		ui.Text(c, title).FontSize(27).FontWeight(800).LetterSpacing(-1)
	})
}

// row 是卡片里的一行：图标 + 标题 + 右侧内容。
func (a *App) row(c *ui.Context, p Palette, iconName, label string, trailing func()) ui.Element {
	row := ui.Row(c).Gap(12).AlignItems(ui.Center).Padding(11, 14).Height(48)
	row.Children(func() {
		ui.Icon(c, icon(iconName)).FontSize(17).TextColor(p.Text)
		ui.Text(c, label).Grow(1).FontSize(13.5).FontWeight(700).SingleLine()
		if trailing != nil {
			trailing()
		}
	})
	return row
}

// rowHover 只给能点的行加悬停底色，纯展示的行不加。
func rowHover(row ui.Element, p Palette) ui.Element {
	if row.Hovered() {
		row.Background(glassLift(p, 0.06))
	}
	return row
}

func (a *App) valueRow(c *ui.Context, p Palette, iconName, label, value string) ui.Element {
	return a.row(c, p, iconName, label, func() {
		ui.Text(c, value).FontSize(13).TextColor(p.Text).SingleLine()
	})
}

func (a *App) portRow(c *ui.Context, p Palette) {
	row := a.row(c, p, "plug", locales.SettingsPort(), func() {
		input := a.fieldInput(c, p, "port", &a.port, controlWidth)
		input.OnSubmit(func() { services.Storage().Set("port", a.port) })
		if input.Changed() {
			services.Storage().Set("port", a.port)
		}
	})
	rowHover(row, p)
}

// deviceRow 让设备名就地在输入框里改，改完写回存储，mDNS 广播读的就是它。
func (a *App) deviceRow(c *ui.Context, p Palette) {
	row := a.row(c, p, "device-laptop", locales.SettingsDeviceName(), func() {
		input := a.fieldInput(c, p, "device", &a.device, controlWidth)
		input.OnSubmit(func() { services.Storage().Set("uname", a.device) })
		if input.Changed() {
			services.Storage().Set("uname", a.device)
		}
	})
	rowHover(row, p)
}

// controlWidth 是每行右侧操作区的统一宽度，三个控件对齐成一条竖线。
const controlWidth float32 = 168

// fieldInput 是行右侧那种小输入框，一眼能看出可以改。
func (a *App) fieldInput(c *ui.Context, p Palette, key string, value *string, width float32) ui.Element {
	input := ui.TextInput(c.Key(key), value).Width(width).Height(28).Radius(7).
		FontSize(13).TextAlign(ui.End).Padding(0, 8).
		Background(p.Canvas).Border(1, p.Border)
	if input.Focused() {
		input.Border(1, p.Accent)
	}
	return input
}

func (a *App) languageRow(c *ui.Context, p Palette) {
	row := a.row(c, p, "language", locales.SettingsLanguage(), func() {
		sel := ui.SelectBase(c.Key("language"), &a.language)
		sel.Trigger.Width(controlWidth).Height(30).Padding(0, 10).Radius(8).Gap(6).
			AlignItems(ui.Center).Background(p.Canvas).Border(1, p.Border).
			Cursor(ui.CursorPointer)
		sel.Trigger.Children(func() {
			ui.Text(c, i18n.Autonym(a.language)).FontSize(13).SingleLine()
			ui.Icon(c, icon("chevron-down")).FontSize(14).TextColor(p.Muted)
		})
		sel.Popup(func(panel ui.Element) {
			panel.Margin(4, 0, 0, 0).Padding(4).Radius(10).
				Background(p.Canvas).Border(1, p.Border)
			for _, tag := range languages {
				item := sel.Item(tag).Padding(7, 10).Radius(6).Cursor(ui.CursorPointer)
				if item.Highlighted() {
					item.Background(p.Accent)
				}
				item.Children(func() {
					ui.Text(c, i18n.Autonym(tag)).FontSize(13)
				})
			}
		})
	})
	rowHover(row, p)
}
