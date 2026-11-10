package desktop

import (
	"math"
	"time"

	"transok/locales"

	"github.com/egoist/mygo/ui"
)

// dropGlow 是文件悬在窗口上时的浮层：整窗压暗，几团色斑缓慢荡漾，
// 像 Apple Music 的流动渐变；中间提示松手即可添加。
func (a *App) dropGlow(c *ui.Context, open bool) {
	p := a.palette()

	// 浮层常驻构建，拖拽结束靠 Exit 淡出，不会一下子消失。
	ui.Overlay(c, func() {
		if !open {
			return
		}

		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Background(ui.RGBA(0, 0, 0, 0.55)).PassThrough().
			Transition(ui.ElementTransition{
				Duration: 160 * time.Millisecond,
				Enter:    &ui.Motion{Opacity: 0},
				Exit:     &ui.Motion{Opacity: 0},
			})

		back.Children(func() {
			// 色斑铺得比窗口大一圈，模糊之后边缘才不会露白。
			blobs := ui.Box(c).Absolute().Left(-80).Top(-80).Right(-80).Bottom(-80).
				PassThrough()
			blobs.Blur(70)
			blobs.Children(func() {
				flow := blobs.Loop("drop-flow", 9000*time.Millisecond, ui.Linear)
				for i, col := range dropColors {
					phase := float64(i)*1.7 + 0.6
					t := 2 * math.Pi * (float64(flow) + float64(i)/float64(len(dropColors)))
					blob := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
						PassThrough()
					blob.RadialGradient(ui.RadialGradient{
						Center: ui.Point{
							X: float32(0.5 + 0.34*math.Sin(t+phase)),
							Y: float32(0.5 + 0.30*math.Cos(t*0.8+phase)),
						},
						Radius: 210,
						Oklab:  true,
						Stops: []ui.GradientStop{
							{Color: withAlpha(col, 0.55), At: 0},
							{Color: withAlpha(col, 0.22), At: 0.55},
							{Color: withAlpha(col, 0), At: 1},
						},
					})
				}
			})

			ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().
				Gap(12).PassThrough().Children(func() {
				tile := ui.Column(c).Size(56, 56).Radius(999).Center().Background(p.Accent)
				tile.Shadow(0, 10, 30, 0, ui.RGBA(124, 126, 44, 0.55))
				tile.Children(func() {
					ui.Icon(c, icon("upload")).FontSize(26).TextColor(p.AccentText)
				})
				ui.Text(c, locales.HomeDropRelease()).FontSize(15).FontWeight(700).
					TextColor(ui.Hex("#ffffff"))
			})
		})
	})
}

// dropColors 是荡漾的几团颜色：橄榄主色打底，配冷色与暖色各一，
// 混在一起才有流动感。
var dropColors = []ui.Color{
	ui.Hex("#7c7e2c"),
	ui.Hex("#4f7d5a"),
	ui.Hex("#8a6a3a"),
	ui.Hex("#6b6f2f"),
}
