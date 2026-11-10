package desktop

import (
	"transok/locales"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// openCaptcha 把已保存的口令读进格子，再打开抽屉。
func (a *App) openCaptcha() {
	for i := range a.captchaDigits {
		a.captchaDigits[i] = ""
	}
	for i := 0; i < len(a.captcha) && i < captchaLength; i++ {
		a.captchaDigits[i] = string(a.captcha[i])
	}
	a.focusReset = true
	a.captchaOpen = true
}

// shareDrawer 列出本机可访问的分享地址，点一行就复制。
func (a *App) shareDrawer(c *ui.Context, p Palette) {
	a.drawer(c, p, &a.shareOpen, drawerSpec{
		icon:        "share",
		title:       locales.ShareTitle(),
		description: locales.ShareDesc(),
		body: func() {
			urls := a.shareURLs()
			if len(urls) == 0 {
				ui.Text(c, locales.ShareEmpty()).FontSize(12.5).TextColor(p.Muted)
				return
			}
			ui.Column(c).FillWidth().Gap(8).Children(func() {
				for i, url := range urls {
					a.shareURLRow(c, p, i, url)
				}
			})
		},
		footer: func() {
			ui.Spacer(c).Grow(1)
			dialogButton(c, p, locales.ShareSeeAll(), func() {
				a.shareOpen = false
				a.page = allURLRoute
			})
		},
	})
}

func (a *App) shareURLRow(c *ui.Context, p Palette, i int, url string) {
	row := ui.ButtonBase(c).FillWidth().Padding(10, 12).Gap(10).Radius(10).
		Background(p.Canvas).Border(1, p.Border).Cursor(ui.CursorPointer).
		OnClick(func() { mygo.Clipboard.WriteText(url) })
	if row.Hovered() {
		row.Border(1, ui.RGBA(124, 126, 44, 0.45))
	}
	row.Children(func() {
		ui.Textf(c, "%02d", i+1).FontSize(11).TextColor(p.Muted).Width(18)
		ui.Text(c, url).Grow(1).FontSize(12).Font("monospace").SingleLine()
		ui.Icon(c, icon("copy")).FontSize(15).TextColor(p.Accent)
	})
}
