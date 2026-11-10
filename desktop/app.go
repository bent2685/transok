package desktop

import (
	"sync"

	"transok/backend/services"
	"transok/backend/utils/mdns/handlers"

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
	win       *mygo.Window
	dark      bool
	themeInit bool
	page      route

	appInfo  map[string]string
	port     string
	captcha  string
	language string

	files []services.FileInfo
	share bool
	list  ui.ListState

	uploadMenu  bool
	applied     string
	devices     []handlers.DiscoverDevice
	deviceList  ui.ListState
	devicesOnce sync.Once

	captchaOpen   bool
	captchaDigits []string
	captchaInputs [captchaLength]ui.Handle
	focusReset    bool

	noteOpen  bool
	noteID    string
	noteDraft string

	textOpen  bool
	textID    string
	textDraft string
}

// New 读取已持久化的状态，视图第一帧就能显示正确的内容。
func New() *App {
	app := &App{appInfo: services.System().GetAppInfo(), captchaDigits: make([]string, captchaLength)}

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

func (a *App) Attach(win *mygo.Window) {
	a.win = win
	a.language, a.applied = mygo.App.Locale(), mygo.App.Locale()
	a.watchDevices()
}

func (a *App) View(c *ui.Context) {
	if !a.themeInit {
		a.dark, a.themeInit = c.Theme().Dark, true
	}
	a.syncLanguage()
	c.SetTheme(Theme(PaletteFor(a.dark)))

	switch a.page {
	case settingsRoute:
		a.settings(c)
	case discoverRoute:
		a.discover(c)
	case allURLRoute:
		a.allURLs(c)
	default:
		a.home(c)
	}
}

// syncLanguage 把选择器里的语言应用到 App，只在真的换了以后切。
func (a *App) syncLanguage() {
	if a.language == "" || a.language == a.applied {
		return
	}
	a.applied = a.language
	mygo.App.SetLocale(a.language)
}

// languages 是 locales/ 里提供的语言，与目录文件保持一致。
var languages = []string{"en", "zh-Hans", "zh-Hant", "ja"}

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
