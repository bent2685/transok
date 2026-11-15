package desktop

import (
	"testing"

	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo/ui"
)

func newTestApp() *App {
	return &App{
		appInfo:       map[string]string{"name": "Transok", "version": "0.5.2"},
		port:          "9482",
		captchaDigits: make([]string, captchaLength),
	}
}

func TestHomeShowsEmptyState(t *testing.T) {
	tt := ui.NewTester(newTestApp().View, 380, 620)

	if !tt.HasText("Transok") {
		t.Errorf("header missing, texts %q", tt.Texts())
	}
	for _, want := range []string{locales.HomeEmptyTitle(), locales.HomeEmptyHint()} {
		if !tt.HasText(want) {
			t.Errorf("empty state missing %q, texts %q", want, tt.Texts())
		}
	}
	if !tt.HasText(locales.HomeDropHint()) {
		t.Errorf("drop hint missing, texts %q", tt.Texts())
	}
}

func TestHomeListsSharedFiles(t *testing.T) {
	app := newTestApp()
	app.files = []services.FileInfo{
		{Id: "1", Name: "photo.jpg", Path: "/tmp/photo.jpg", Size: 2048, Type: "jpg"},
	}
	tt := ui.NewTester(app.View, 380, 620)

	if !tt.HasText("photo.jpg") {
		t.Errorf("file name missing, texts %q", tt.Texts())
	}
	if !tt.HasText("2.00KB") {
		t.Errorf("file size missing, texts %q", tt.Texts())
	}
	if !tt.HasText("(1)") {
		t.Errorf("file count missing, texts %q", tt.Texts())
	}
}

func TestSettingsShowsSections(t *testing.T) {
	app := newTestApp()
	app.page = settingsRoute
	tt := ui.NewTester(app.View, 380, 620)

	for _, want := range []string{
		locales.SettingsTitle(),
		locales.SettingsCommon(),
		locales.SettingsAbout(),
		locales.SettingsPort(),
	} {
		if !tt.HasText(want) {
			t.Errorf("%q missing, texts %q", want, tt.Texts())
		}
	}
}

func TestCaptchaSheetShows(t *testing.T) {
	app := newTestApp()
	app.captchaOpen = true
	tt := ui.NewTester(app.View, 380, 620)

	for _, want := range []string{
		locales.CaptchaTitle(),
		locales.CaptchaDescription(captchaLength),
		locales.DialogConfirm(),
	} {
		if !tt.HasText(want) {
			t.Errorf("%q missing, texts %q", want, tt.Texts())
		}
	}
}

func TestJoinCode(t *testing.T) {
	if got := joincode([]string{"1", "2", "", "4", "5"}); got != "1245" {
		t.Errorf("code %q, want 1245", got)
	}
}

func TestPureTextCreatedOnlyOnConfirm(t *testing.T) {
	app := newTestApp()

	app.addPureText()
	if len(app.files) != 0 {
		t.Fatalf("files %d after opening the drawer, want 0", len(app.files))
	}

	// 取消：什么都不留
	app.textOpen = false
	app.textNew = false
	if len(app.files) != 0 {
		t.Fatalf("files %d after canceling, want 0", len(app.files))
	}

	// 确认：落成一条
	app.addPureText()
	app.textName, app.textDraft = "标题", "正文"
	app.saveText()
	if len(app.files) != 1 {
		t.Fatalf("files %d after confirming, want 1", len(app.files))
	}
	if got := app.files[0]; got.Name != "标题" || got.Text != "正文" || got.Type != "pure-text" {
		t.Errorf("item %+v", got)
	}
	if app.files[0].Id == "" {
		t.Error("id is empty")
	}
}

func TestCaptchaBackspaceClearsAllBoxes(t *testing.T) {
	app := newTestApp()
	app.captchaOpen = true
	copy(app.captchaDigits, []string{"1", "2", "3", "4", "5"})
	focused := false
	tt := ui.NewTester(func(c *ui.Context) {
		app.View(c)
		if !focused {
			focused = true
			app.captchaInputs[captchaLength-1].Focus(c)
		}
	}, 380, 620)
	tt.Frame()
	for i := 0; i < captchaLength; i++ {
		tt.Key(0, ui.KeyBackspace)
		tt.Frame()
	}
	if got := joincode(app.captchaDigits); got != "" {
		t.Errorf("digits %q after %d backspaces, want empty", app.captchaDigits, captchaLength)
	}
}

func TestWindowAcceptsFilesOverIcons(t *testing.T) {
	for _, overlay := range []bool{false, true} {
		name := "首页"
		if overlay {
			name = "拖拽浮层"
		}
		t.Run(name, func(t *testing.T) {
			app := newTestApp()
			tt := ui.NewTester(func(c *ui.Context) {
				app.View(c)
				if overlay {
					app.dropGlow(c, true)
				}
			}, 380, 620)
			// 不存在的路径不会写入分享清单，只验证拖放事件是否被窗口接住。
			path := t.TempDir() + "/missing"
			for _, point := range []ui.Point{
				{X: 30, Y: 230},
				{X: 190, Y: 294}, // 浮层中央的上传图标。
				{X: 190, Y: 566}, // 底部的主要按钮图标。
				{X: 190, Y: 130}, // 上传区。
				{X: 190, Y: 350}, // 空列表。
			} {
				if !tt.DropFiles(point.X, point.Y, path) {
					t.Errorf("窗口未接住 (%g, %g) 的文件拖放", point.X, point.Y)
				}
				tt.Frame()
			}
		})
	}
}
