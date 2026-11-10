package desktop

import (
	"testing"

	"transok/backend/services"

	"github.com/egoist/mygo/ui"
)

func newTestApp() *App {
	return &App{appInfo: map[string]string{"name": "Transok", "version": "0.5.2"}, port: "9482"}
}

func TestHomeShowsEmptyState(t *testing.T) {
	tt := ui.NewTester(newTestApp().View, 380, 620)

	if !tt.HasText("Transok") {
		t.Errorf("header missing, texts %q", tt.Texts())
	}
	if !tt.HasText("还没有内容") {
		t.Errorf("empty state missing, texts %q", tt.Texts())
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
	if !tt.HasText("2.0 KB") {
		t.Errorf("file size missing, texts %q", tt.Texts())
	}
	if !tt.HasText("(1)") {
		t.Errorf("file count missing, texts %q", tt.Texts())
	}
}
