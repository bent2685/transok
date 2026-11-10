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
	if !tt.HasText(locales.HomeNoFiles()) {
		t.Errorf("empty state missing, texts %q", tt.Texts())
	}
	if !tt.HasText(locales.HomeUploadButton()) {
		t.Errorf("upload button missing, texts %q", tt.Texts())
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
