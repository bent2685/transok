package desktop

import (
	"image"
	"image/color"

	"transok/locales"

	"github.com/egoist/mygo/ui"
	qrcode "github.com/skip2/go-qrcode"
)

const qrSize = 220

// qrBitmap 把地址编成二维码位图（含静区），按像素风渲染以保持边缘锐利。
func qrBitmap(url string) (*ui.Bitmap, error) {
	code, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	modules := code.Bitmap()
	img := image.NewGray(image.Rect(0, 0, len(modules), len(modules)))
	for y, row := range modules {
		for x, dark := range row {
			if dark {
				img.SetGray(x, y, color.Gray{Y: 0})
				continue
			}
			img.SetGray(x, y, color.Gray{Y: 255})
		}
	}
	return ui.NewBitmap(img).Pixelated(), nil
}

// showQR 为某个地址生成二维码并打开弹窗；分享抽屉先收起，免得两层叠着。
func (a *App) showQR(url string) {
	bmp, err := qrBitmap(url)
	if err != nil {
		return
	}
	a.qrImage = bmp
	a.qrURL = url
	a.shareOpen = false
	a.qrOpen = true
}

// qrDrawer 是扫码弹窗：白底的二维码下面写着对应的地址。
func (a *App) qrDrawer(c *ui.Context, p Palette) {
	a.drawer(c, p, &a.qrOpen, drawerSpec{
		icon:        "qrcode",
		title:       locales.ShareQRTitle(),
		description: locales.ShareQRDesc(),
		center:      true,
		body: func() {
			if a.qrImage == nil {
				return
			}
			ui.Column(c).FillWidth().Center().Gap(12).Children(func() {
				ui.Column(c).Radius(12).Background(ui.RGB(255, 255, 255)).Children(func() {
					ui.Image(c, a.qrImage).Size(qrSize, qrSize).Radius(12)
				})
				ui.Text(c, a.qrURL).FontSize(12).Font("monospace").TextColor(p.Muted).SingleLine()
			})
		},
	})
}

// qrButton 是地址行末尾的二维码按钮。
func (a *App) qrButton(c *ui.Context, p Palette, url string) {
	b := ui.ButtonBase(c.Key("qr"+url)).Size(40, 40).Padding(0).Radius(10).
		Background(p.Surface).Border(1, p.Border).Cursor(ui.CursorPointer).
		OnClick(func() { a.showQR(url) })
	tint := p.Muted
	if b.Hovered() {
		b.Border(1, ui.RGBA(124, 126, 44, 0.45))
		tint = p.Accent
	}
	b.Children(func() {
		ui.Icon(c, icon("qrcode")).FontSize(16).TextColor(tint)
	})
}
