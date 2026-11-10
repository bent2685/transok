package desktop

import (
	"sync"
	"time"

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
	noteInput ui.Handle

	textOpen    bool
	textNew     bool
	textID      string
	textName    string
	textDraft   string
	textTitle   ui.Handle
	textContent ui.Handle
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

	ui.Column(c).Fill().Children(func() {
		// 顶部原本是系统标题栏的区域：整条都能拖动窗口，双击则缩放/最小化。
		ui.Box(c).FillWidth().Height(a.titleBarStrip(c)).DragWindow()
		ui.Column(c).Grow(1).Children(func() {
			// 关键值按页区分：切页时新页淡入上移，旧页淡出，像共享轴转场。
			ui.Column(c).Key(a.page).Fill().Transition(pageMotion).Children(func() {
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
			})
		})
	})
}

// pageMotion 是页面之间的转场：进来时从下方淡入，离开时向上淡出。
var pageMotion = ui.ElementTransition{
	Duration:     200 * time.Millisecond,
	ExitDuration: 140 * time.Millisecond,
	Enter:        &ui.Motion{Y: 10},
	Exit:         &ui.Motion{Y: -8},
}

// titleBarStrip 是留给窗口按钮的那条空白带，再加上一点呼吸。
func (a *App) titleBarStrip(c *ui.Context) float32 {
	h := c.TitleBar().Height + 6
	if h < 34 {
		h = 34
	}
	return h
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
