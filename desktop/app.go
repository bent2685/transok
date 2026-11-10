package desktop

import (
	"transok/backend/services"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type route int

const (
	homeRoute route = iota
	settingsRoute
	discoverRoute
	allURLRoute
)

// App 是窗口的状态：视图每帧从它重建，所以这里只放数据、widget state
// 和 ui.Handle，绝不放 ui.Element。
type App struct {
	win  *mygo.Window
	dark bool
	themeInit bool
	page route

	appInfo map[string]string
	port    string
	captcha string

	files []services.FileInfo
	share bool
	list  ui.ListState

	uploadMenu   bool
	captchaOpen  bool
	captchaDraft string
}

// New 读取已持久化的状态，视图第一帧就能显示正确的内容。
func New() *App {
	app := &App{appInfo: services.System().GetAppInfo()}

	app.port = "9482"
	if port, ok := services.Storage().Get("port"); ok {
		app.port = port.(string)
	}
	if share, ok := services.Storage().Get("is-share"); ok {
		app.share = share.(bool)
	}
	app.files = services.File().GetShareList()
	app.captcha = services.Share().GetCaptcha()
	return app
}

// Attach 拿到窗口，主题切换、后台更新都通过它重画。
func (a *App) Attach(win *mygo.Window) {
	a.win = win
}

func (a *App) View(c *ui.Context) {
	if !a.themeInit {
		a.dark, a.themeInit = c.Theme().Dark, true
	}
	c.SetTheme(Theme(a.dark))
	ui.Column(c).Fill().Gap(12).Padding(16).Children(func() {
		switch a.page {
		case settingsRoute:
			a.settings(c)
		case discoverRoute:
			a.placeholder(c, "发现设备")
		case allURLRoute:
			a.placeholder(c, "全部分享链接")
		default:
			a.home(c)
		}
	})
}

// SetDark 跟随系统外观，主题换了要重画。
func (a *App) SetDark(dark bool) {
	if a.dark == dark {
		return
	}
	if a.win == nil {
		a.dark = dark
		return
	}
	a.win.Update(func() { a.dark = dark })
}

// header 是三个屏共用的标题栏：返回键 + 标题。
func (a *App) header(c *ui.Context, title string) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Button(c, "").Padding(6).Radius(6).Children(func() {
			ui.Icon(c, icon("arrow-left")).FontSize(16)
		}).OnClick(func() { a.page = homeRoute })
		ui.Text(c, title).FontSize(18).Bold()
	})
}

func (a *App) placeholder(c *ui.Context, title string) {
	a.header(c, title)
	ui.Column(c).Fill().Center().Gap(6).Children(func() {
		ui.Icon(c, icon("file")).FontSize(28).TextColor(c.Theme().TextMuted)
		ui.Text(c, "迁移中").TextColor(c.Theme().TextMuted)
	})
}
