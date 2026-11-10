package desktop

import (
	"strings"
	"time"

	"transok/locales"

	"github.com/egoist/mygo/ui"
)

const captchaLength = 5

// 抽屉：背景淡入淡出，面板从自己的高度滑上来，回弹一点。
var (
	sheetBackdrop = ui.ElementTransition{
		Duration: 240 * time.Millisecond,
		Enter:    &ui.Motion{Opacity: 0},
		Exit:     &ui.Motion{Opacity: 0},
	}
	sheetPanel = ui.ElementTransition{
		Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.16}},
		Enter:     &ui.Motion{YPercent: 100, Opacity: 1},
		Exit:      &ui.Motion{YPercent: 100, Opacity: 1},
	}
)

// dialogPop 是居中弹窗的进出场：淡入并轻微放大。
var dialogPop = ui.ElementTransition{
	Animation: &ui.Animation{Spring: ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.12}},
	Enter:     &ui.Motion{Scale: 0.94, Opacity: 0.6},
	Exit:      &ui.Motion{Scale: 0.96, Opacity: 0.3},
}

// modalCard 是居中弹窗：遮罩淡入，卡片带进出场，点遮罩或 Esc 关闭。
func modalCard(c *ui.Context, p Palette, open *bool, build func()) {
	ui.Overlay(c, func() {
		if !*open {
			return
		}
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Background(ui.RGBA(0, 0, 0, 0.62)).Modal().Transition(sheetBackdrop)
		back.OnClick(func() { *open = false })

		wrap := ui.Row(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().
			Transition(dialogPop)
		wrap.Children(func() {
			card := ui.Column(c).Width(300).Padding(18).Gap(10).Radius(14).
				Background(p.Surface).Border(1, p.Border)
			card.Shadow(0, 12, 30, 0, ui.RGBA(0, 0, 0, 0.35))
			if card.OverlayShortcut(0, ui.KeyEscape) {
				*open = false
			}
			card.Children(build)
		})
	})
}

// captchaSheet 是设置口令的底部抽屉，5 个格子逐位输入。
func (a *App) captchaSheet(c *ui.Context, p Palette) {
	ui.Overlay(c, func() {
		if !a.captchaOpen {
			return
		}

		// 遮罩与面板是兄弟：各自的进出场才互不牵连。
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Background(ui.RGBA(0, 0, 0, 0.62)).Modal().Transition(sheetBackdrop)
		back.OnClick(func() { a.captchaOpen = false })

		wrap := ui.Column(c).Absolute().Left(0).Right(0).Bottom(0).
			Justify(ui.End).Transition(sheetPanel)
		wrap.Children(func() {
			panel := ui.Column(c).FillWidth().Radius(18, 18, 0, 0).
				Padding(8, 20, 18, 20).Background(p.Surface)
			if panel.OverlayShortcut(0, ui.KeyEscape) {
				a.captchaOpen = false
			}

			panel.Children(func() {
				// 托手居中对齐；关闭键单独浮在右上角，不参与居中计算。
				ui.Row(c).Height(28).Justify(ui.Center).AlignItems(ui.Center).PassThrough().Children(func() {
					ui.Box(c).Size(44, 4).Radius(999).Background(p.Handle)
				})
				ui.Row(c).Absolute().Right(14).Top(8).Children(func() {
					ghostButton(c, p, "x", func() { a.captchaOpen = false })
				})

				spacer(c, 14)
				ui.Row(c).Justify(ui.Center).Children(func() {
					tile := ui.Column(c).Size(50, 50).Radius(13).Center().Background(p.Accent)
					tile.Shadow(0, 6, 18, 0, ui.RGBA(124, 126, 44, 0.45))
					tile.Children(func() {
						ui.Icon(c, icon("shield-lock")).FontSize(26).TextColor(p.AccentText)
					})
				})

				spacer(c, 14)
				ui.Row(c).Justify(ui.Center).Children(func() {
					ui.Text(c, locales.CaptchaTitle()).FontSize(24).FontWeight(800).LetterSpacing(-0.6)
				})
				spacer(c, 7)
				ui.Row(c).Justify(ui.Center).Children(func() {
					ui.Text(c, locales.CaptchaDescription(captchaLength)).FontSize(13).TextColor(p.Muted)
				})

				spacer(c, 24)
				a.otpBoxes(c, p)

				spacer(c, 26)
				ui.Row(c).Gap(14).AlignItems(ui.Center).Children(func() {
					dialogButton(c, p, locales.DialogCancel(), func() { a.captchaOpen = false })
					dialogIconButton(c, p, "refresh", p.Text, locales.CaptchaReset(), a.resetCaptcha)
					ui.Spacer(c).Grow(1)
					dialogIconButton(c, p, "check", p.Accent, locales.DialogConfirm(), func() {
						a.applyCaptcha()
						a.captchaOpen = false
					})
				})
			})
		})
	})
}

func (a *App) otpBoxes(c *ui.Context, p Palette) {
	next := -1
	ui.Row(c).Gap(12).Justify(ui.Center).Children(func() {
		for i := 0; i < captchaLength; i++ {
			box := ui.TextInput(c.Key("otp"+string(rune('0'+i))), &a.captchaDigits[i]).
				Bind(&a.captchaInputs[i]).Width(46).Height(54).Radius(10).
				FontSize(18).FontWeight(700).TextAlign(ui.Center).MaxLength(1).
				Background(p.Surface).Border(1, p.Border)
			if a.captchaInputs[i].Focused(c) {
				box.Border(1, p.Accent)
			}
			if box.Changed() && a.captchaDigits[i] != "" {
				next = i + 1
			}
		}
	})
	if a.focusReset {
		a.focusReset = false
		a.captchaInputs[0].Focus(c)
	} else if next > 0 && next < captchaLength {
		a.captchaInputs[next].Focus(c)
	}
}

func (a *App) resetCaptcha() {
	for i := range a.captchaDigits {
		a.captchaDigits[i] = ""
	}
	a.focusReset = true
}

func joincode(digits []string) string {
	return strings.Join(digits, "")
}
