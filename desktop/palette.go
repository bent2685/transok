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
}

func PaletteFor(dark bool) Palette {
	if dark {
		return Palette{
			Dark:         true,
			Canvas:       ui.Hex("#0a0a0a"),
			Surface:      ui.Hex("#1a1a1a"),
			SurfaceHover: ui.Hex("#242424"),
			Pressed:      ui.Hex("#2f2f2f"),
			Border:       ui.Hex("#2a2a2a"),
			Text:         ui.Hex("#ffffff"),
			Subtle:       ui.Hex("#cccccc"),
			Muted:        ui.Hex("#8a8a8a"),
			Accent:       ui.Hex("#7c7e2c"),
			AccentText:   ui.Hex("#ffffff"),
			Danger:       ui.Hex("#ef4444"),
			Elevated:     ui.Hex("#2b2b2b"),
		}
	}
	return Palette{
		Canvas:       ui.Hex("#ffffff"),
		Surface:      ui.Hex("#f6f6f6"),
		SurfaceHover: ui.Hex("#efefef"),
		Pressed:      ui.Hex("#e6e6e6"),
		Border:       ui.Hex("#e4e4e4"),
		Text:         ui.Hex("#0a0a0a"),
		Subtle:       ui.Hex("#3d3d3d"),
		Muted:        ui.Hex("#767676"),
		Accent:       ui.Hex("#7c7e2c"),
		AccentText:   ui.Hex("#ffffff"),
		Danger:       ui.Hex("#dc2626"),
		Elevated:     ui.Hex("#ececec"),
	}
}

func (a *App) palette() Palette { return PaletteFor(a.dark) }
