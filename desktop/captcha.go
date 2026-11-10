package desktop

import (
	"strings"

	"transok/locales"

	"github.com/egoist/mygo/ui"
)

const captchaLength = 5

// captchaDrawer 是设置口令的抽屉，5 个格子逐位输入。
func (a *App) captchaDrawer(c *ui.Context, p Palette) {
	a.drawer(c, p, &a.captchaOpen, drawerSpec{
		icon:        "shield-lock",
		title:       locales.CaptchaTitle(),
		description: locales.CaptchaDescription(captchaLength),
		center:      true,
		body: func() {
			a.otpBoxes(c, p)
			// 格子里下面留一段空，按钮不至于贴着。
			spacer(c, 10)
		},
		footer: func() {
			dialogGhostButton(c, p, "refresh", p.Text, locales.CaptchaReset(), a.resetCaptcha)
			ui.Spacer(c).Grow(1)
			dialogIconButton(c, p, "check", p.Accent, locales.DialogConfirm(), func() {
				a.applyCaptcha()
				a.captchaOpen = false
			})
		},
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
