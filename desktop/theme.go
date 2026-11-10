package desktop

import "github.com/egoist/mygo/ui"

// Theme 把调色板映射成 native UI 的主题，控件（按钮、输入框、滚动条）
// 跟着同一套颜色走。
func Theme(p Palette) *ui.Theme {
	t := ui.LightTheme()
	if p.Dark {
		t = ui.DarkTheme()
	}

	t.Dark = p.Dark
	t.Background = p.Canvas
	t.Surface, t.SurfaceHover, t.SurfacePressed = p.Surface, p.SurfaceHover, p.Pressed
	t.Border = p.Border
	t.Text, t.TextMuted = p.Text, p.Muted
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText = p.Accent, ui.Hex("#8b8d38"), ui.Hex("#686a24"), p.AccentText
	t.InactiveAccent = p.Accent
	t.Danger, t.Warning = p.Danger, p.Danger
	t.Selection = ui.RGBA(124, 126, 44, 0.35)
	t.Focus = p.Accent
	t.Radius, t.Spacing, t.FontSize = 8, 4, 13
	t.Scrollbar, t.ScrollbarWidth = p.Border, 5
	return t
}
