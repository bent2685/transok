package desktop

import "github.com/egoist/mygo/ui"

// Palette 是原 Web 端 CSS 变量（base.scss）的 Go 版，深浅两套共用同一个橄榄主色。
type Palette struct {
	Dark         bool
	Canvas       ui.Color
	Surface      ui.Color
	SurfaceHover ui.Color
	Pressed      ui.Color
	Border       ui.Color
	Text         ui.Color
	Subtle       ui.Color
	Muted        ui.Color
	Accent       ui.Color
	AccentText   ui.Color
	Danger       ui.Color
	// Elevated 是浮层里按钮这类浮在外壳之上的面。
	Elevated ui.Color
	// Fade 是页面底部渐隐的落点：深色落纯黑，浅色落纯白。
	Fade ui.Color
}

func PaletteFor(dark bool) Palette {
	if dark {
		return Palette{
			Dark:         true,
			Canvas:       ui.Hex("#141414"),
			Surface:      ui.Hex("#1d1d1d"),
			SurfaceHover: ui.Hex("#272727"),
			Pressed:      ui.Hex("#313131"),
			Border:       ui.Hex("#2e2e2e"),
			Text:         ui.Hex("#ffffff"),
			Subtle:       ui.Hex("#cccccc"),
			Muted:        ui.Hex("#8a8a8a"),
			Accent:       ui.Hex("#7c7e2c"),
			AccentText:   ui.Hex("#ffffff"),
			Danger:       ui.Hex("#ef4444"),
			Elevated:     ui.Hex("#2b2b2b"),
			Fade:         ui.Hex("#000000"),
		}
	}
	return Palette{
		Canvas:       ui.Hex("#f2f2f3"),
		Surface:      ui.Hex("#ffffff"),
		SurfaceHover: ui.Hex("#ececec"),
		Pressed:      ui.Hex("#e2e2e2"),
		Border:       ui.Hex("#e4e4e4"),
		Text:         ui.Hex("#0a0a0a"),
		Subtle:       ui.Hex("#3d3d3d"),
		Muted:        ui.Hex("#767676"),
		Accent:       ui.Hex("#7c7e2c"),
		AccentText:   ui.Hex("#ffffff"),
		Danger:       ui.Hex("#dc2626"),
		Elevated:     ui.Hex("#ececec"),
		Fade:         ui.Hex("#ffffff"),
	}
}

func (a *App) palette() Palette { return PaletteFor(a.dark) }
