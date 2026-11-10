package desktop

import (
	"strings"

	"transok/locales"

	"github.com/egoist/mygo/ui"
)

const captchaLength = 5

// captchaSheet 是设置口令的底部弹层，5 个格子逐位输入。
func (a *App) captchaSheet(c *ui.Context, p Palette) {
	ui.Overlay(c, func() {
		back := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			Background(ui.RGBA(0, 0, 0, 0.62)).Modal().Justify(ui.End)
		back.Children(func() {
			panel := ui.Column(c).FillWidth().Radius(18, 18, 0, 0).
				Padding(10, 20, 18, 20).Gap(0).Background(p.Surface)
			if panel.PressedOutside() || panel.OverlayShortcut(0, ui.KeyEscape) {
				a.captchaOpen = false
			}

			panel.Children(func() {
				// 顶部手柄 + 关闭
				ui.Row(c).AlignItems(ui.Center).Children(func() {
					ui.Row(c).Grow(1).Justify(ui.Center).Children(func() {
						ui.Box(c).Size(36, 4).Radius(999).Background(p.Border)
					})
					ghostButton(c, p, "x", func() { a.captchaOpen = false })
				})

				spacer(c, 10)
				ui.Row(c).Justify(ui.Center).Children(func() {
					tile := ui.Column(c).Size(52, 52).Radius(14).Center().Background(p.Accent)
					tile.Shadow(0, 6, 18, 0, ui.RGBA(124, 126, 44, 0.45))
					tile.Children(func() {
						ui.Icon(c, icon("shield-lock")).FontSize(26).TextColor(p.AccentText)
					})
				})

				spacer(c, 12)
				ui.Row(c).Justify(ui.Center).Children(func() {
					ui.Text(c, locales.CaptchaTitle()).FontSize(22).FontWeight(800).LetterSpacing(-0.4)
				})
				spacer(c, 6)
				ui.Row(c).Justify(ui.Center).Children(func() {
					ui.Text(c, locales.CaptchaDescription(captchaLength)).FontSize(12.5).TextColor(p.Muted)
				})

				spacer(c, 18)
				a.otpBoxes(c, p)

				spacer(c, 18)
				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					dialogButton(c, p, locales.DialogCancel(), func() { a.captchaOpen = false })
					dialogIconButton(c, p, "refresh", locales.CaptchaReset(), a.resetCaptcha)
					ui.Spacer(c).Grow(1)
					dialogIconButton(c, p, "check", locales.DialogConfirm(), func() {
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
