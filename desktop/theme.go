package desktop

import "github.com/egoist/mygo/ui"

// 颜色取自 Web 端的语义 token（frontend/src/assets/styles/base.scss），
// 深浅两套共用同一个橄榄主色 #7c7e2c。
var (
	olive        = ui.Hex("#7c7e2c")
	oliveHover   = ui.Hex("#8b8d38")
	olivePressed = ui.Hex("#686a24")
	danger       = ui.Hex("#ef4444")
)

// Theme 把 Web 端的 CSS 变量映射成 native UI 的主题。
func Theme(dark bool) *ui.Theme {
	t := ui.LightTheme()
	background, surface, surfaceHover, surfacePressed := ui.Hex("#ffffff"), ui.Hex("#f5f5f5"), ui.Hex("#ececec"), ui.Hex("#e2e2e2")
	border, text, textMuted := ui.Hex("#e0e0e0"), ui.Hex("#0a0a0a"), ui.Hex("#666666")
	scrollbar := ui.Hex("#c8c8c8")
	if dark {
		t = ui.DarkTheme()
		background, surface, surfaceHover, surfacePressed = ui.Hex("#0a0a0a"), ui.Hex("#1a1a1a"), ui.Hex("#242424"), ui.Hex("#2f2f2f")
		border, text, textMuted = ui.Hex("#2a2a2a"), ui.Hex("#ffffff"), ui.Hex("#cccccc")
		scrollbar = ui.Hex("#3a3a3a")
	}

	t.Dark = dark
	t.Background, t.Surface, t.SurfaceHover, t.SurfacePressed = background, surface, surfaceHover, surfacePressed
	t.Border, t.Text, t.TextMuted = border, text, textMuted
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText = olive, oliveHover, olivePressed, ui.Hex("#ffffff")
	t.InactiveAccent = olive
	t.Danger, t.Warning = danger, danger
	t.Selection = ui.RGBA(124, 126, 44, 0.35)
	t.Focus = olive
	t.Radius, t.Spacing, t.FontSize = 8, 4, 14
	t.Scrollbar = scrollbar
	return t
}
