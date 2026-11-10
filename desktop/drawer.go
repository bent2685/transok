package desktop

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// 抽屉的进出场：遮罩淡入淡出，面板从自身高度滑上来，带一点回弹。
var (
	drawerBackdrop = ui.ElementTransition{
		Duration: 240 * time.Millisecond,
		Enter:    &ui.Motion{Opacity: 0},
		Exit:     &ui.Motion{Opacity: 0},
	}
	drawerPanel = ui.ElementTransition{
		Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.16}},
		Enter:     &ui.Motion{YPercent: 100, Opacity: 1},
		Exit:      &ui.Motion{YPercent: 100, Opacity: 1},
	}
)

// drawerSpec 描述一个抽屉：顶部橄榄图标块、标题与说明，中间由 body 填，
// 底部由 footer 排按钮。二级面板都走这一个组件，观感才一致。
type drawerSpec struct {
	icon        string
	title       string
	description string
	// center 把默认头部改成一列居中，用于短文案的面板（如验证码）。
	center bool
	// header 取代默认的图标 + 标题行，用于标题本身可编辑的抽屉。
	header func()
	body   func()
	footer func()
}

// drawer 是全局的底部抽屉：点遮罩、按 Esc 或点右上角的叉都能收起。
func (a *App) drawer(c *ui.Context, p Palette, open *bool, spec drawerSpec) {
	close := func() { *open = false }

	ui.Overlay(c, func() {
		if !*open {
			return
		}

		// 遮罩与面板是兄弟：各自的进出场才互不牵连。
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Background(ui.RGBA(0, 0, 0, 0.62)).Modal().Transition(drawerBackdrop)
		back.OnClick(close)

		wrap := ui.Column(c).Absolute().Left(0).Right(0).Bottom(0).
			Justify(ui.End).Transition(drawerPanel)
		wrap.Children(func() {
			panel := ui.Column(c).FillWidth().Radius(18, 18, 0, 0).
				Padding(8, 20, 18, 20).Background(p.Surface)
			if panel.OverlayShortcut(0, ui.KeyEscape) {
				close()
			}

			panel.Children(func() {
				// 关闭键浮在右上角的圆底按钮上，不占排版空间。
				ui.Row(c).Absolute().Right(16).Top(14).Children(func() {
					closeButton(c, p, close)
				})

				spacer(c, 26)
				if spec.header != nil {
					spec.header()
				} else if spec.center {
					ui.Column(c).Center().Gap(0).Children(func() {
						tile := ui.Column(c).Size(52, 52).Radius(13).Center().Background(p.Accent)
						tile.Shadow(0, 6, 18, 0, ui.RGBA(124, 126, 44, 0.45))
						tile.Children(func() {
							ui.Icon(c, icon(spec.icon)).FontSize(26).TextColor(p.AccentText)
						})
						spacer(c, 14)
						ui.Text(c, spec.title).FontSize(22).FontWeight(800).LetterSpacing(-0.5)
						if spec.description != "" {
							spacer(c, 6)
							ui.Text(c, spec.description).FontSize(12.5).TextColor(p.Muted)
						}
					})
				} else {
					ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
						tile := ui.Column(c).Size(40, 40).Radius(11).Center().Background(p.Accent)
						tile.Shadow(0, 4, 14, 0, ui.RGBA(124, 126, 44, 0.35))
						tile.Children(func() {
							ui.Icon(c, icon(spec.icon)).FontSize(21).TextColor(p.AccentText)
						})
						ui.Column(c).Grow(1).Gap(3).Children(func() {
							ui.Text(c, spec.title).FontSize(20).FontWeight(800).LetterSpacing(-0.4)
							if spec.description != "" {
								ui.Text(c, spec.description).FontSize(12).TextColor(p.Muted)
							}
						})
					})
				}

				spacer(c, 24)
				if spec.body != nil {
					spec.body()
				}
				if spec.footer != nil {
					spacer(c, 26)
					ui.Row(c).Gap(14).AlignItems(ui.Center).Children(spec.footer)
				}
			})
		})
	})
}
