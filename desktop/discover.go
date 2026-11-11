package desktop

import (
	"fmt"
	"time"

	"transok/backend/services"
	"transok/backend/utils/mdns/handlers"
	"transok/locales"

	"github.com/egoist/mygo/ui"
)

func (a *App) discover(c *ui.Context) {
	p := a.palette()

	ui.Column(c).Fill().Padding(0, 16, 16, 16).Children(func() {
		header := ui.Row(c).Gap(12).AlignItems(ui.Center).DragWindow()
		header.Children(func() {
			circleButton(c, p, "arrow-left", func() { a.page = homeRoute })
			ui.Text(c, locales.DiscoverTitle()).FontSize(27).FontWeight(800).LetterSpacing(-1)
			ui.Spacer(c).Grow(1)
			spin := ui.Icon(c, icon("loader-2")).FontSize(18).TextColor(p.Accent)
			spin.Rotate(spin.Loop("discover-spin", 2*time.Second, ui.Linear) * 360)
		})

		spacer(c, 10)

		if len(a.devices) == 0 {
			ui.Column(c).Grow(1).Center().Gap(12).Children(func() {
				ui.Icon(c, icon("files-off")).FontSize(34).TextColor(p.Muted)
				ui.Text(c, locales.DiscoverNoDevices()).FontSize(13.5).TextColor(p.Muted)
			})
			return
		}

		ui.List(c, &a.deviceList, len(a.devices), func(i int) {
			d := a.devices[i]
			row := ui.Row(c).Padding(9, 10).Gap(10).Radius(10).AlignItems(ui.Center).
				Background(p.Surface).Border(1, p.Border)
			if row.Hovered() {
				row.Border(1, ui.RGBA(124, 126, 44, 0.35))
			}
			row.Children(func() {
				ui.Column(c).Size(36, 36).Radius(8).Center().
					Background(p.SurfaceHover).Children(func() {
					ui.Icon(c, icon("device-laptop")).FontSize(18).TextColor(p.Text)
				})
				ui.Column(c).Grow(1).Gap(3).Children(func() {
					ui.Text(c, d.Uname).FontSize(13.5).FontWeight(700).SingleLine()
					ui.Text(c, fmt.Sprintf("%s:%s", d.IP, d.Port)).FontSize(11.5).
						TextColor(p.Muted).SingleLine()
				})
			})
		}).Grow(1).Gap(8)
	})
}

// allURLs 列出本机可用的分享地址。
func (a *App) allURLs(c *ui.Context) {
	p := a.palette()

	ui.Column(c).Fill().Padding(0, 16, 16, 16).Children(func() {
		dragHeader(c, p, locales.AllURLTitle(), func() { a.page = homeRoute })
		spacer(c, 14)

		urls := a.shareURLs()
		if len(urls) == 0 {
			ui.Column(c).Grow(1).Center().Children(func() {
				ui.Text(c, locales.AllURLEmpty()).FontSize(13).TextColor(p.Muted)
			})
			return
		}

		ui.Column(c).Gap(8).Children(func() {
			for _, url := range urls {
				copied := a.copiedRecently(c, url)
				row := ui.ButtonBase(c).FillWidth().Padding(11, 12).Gap(8).Radius(10).
					AlignItems(ui.Center).Justify(ui.Start).Background(p.Surface).
					Border(1, p.Border).Cursor(ui.CursorPointer).
					OnClick(func() { a.copyURL(url) })
				if row.Hovered() || copied {
					row.Border(1, ui.RGBA(124, 126, 44, 0.45))
				}
				if copied {
					row.Background(ui.RGBA(124, 126, 44, 0.12))
				}
				row.Children(func() {
					ui.Text(c, url).Grow(1).FontSize(12.5).Font("monospace").SingleLine()
					if copied {
						ui.Row(c).Gap(4).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icon("check")).FontSize(14).TextColor(p.Accent)
							ui.Text(c, locales.ShareCopied()).FontSize(11.5).TextColor(p.Accent)
						})
						return
					}
					ui.Icon(c, icon("copy")).FontSize(15).TextColor(p.Muted)
				})
			}
		})
	})
}

func (a *App) shareURLs() []string {
	urls := make([]string, 0, 3)
	for _, ip := range services.System().GetLocalIps() {
		urls = append(urls, fmt.Sprintf("http://%s:%s", ip, a.port))
	}
	return urls
}

// watchDevices 每 3 秒在后台取一次局域网设备，结果回投到视图。
func (a *App) watchDevices() {
	a.devicesOnce.Do(func() {
		go func() {
			for range time.Tick(3 * time.Second) {
				devices := handlers.GetDiscoverHandler().GetDiscoverList()
				a.win.Update(func() { a.devices = devices })
			}
		}()
	})
}
