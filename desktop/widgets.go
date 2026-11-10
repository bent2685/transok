package desktop

import "github.com/egoist/mygo/ui"

// 这是几个页面共用的视觉零件，尺寸对齐迁移前 Web 端的写法。
// 需要自己控制底色边框的地方一律用 ui.ButtonBase：ui.Button 会每帧按主题
// 重刷底色，压不住。

// iconButton 是标题栏右侧那种方块图标按钮。
func iconButton(c *ui.Context, p Palette, name string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Size(36, 36).Padding(0).Radius(8).
		Background(p.Surface).Border(1, p.Border).
		Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover).Border(1, ui.RGBA(124, 126, 44, 0.45))
	}
	b.Children(func() {
		ui.Icon(c, icon(name)).FontSize(17).TextColor(p.Text)
	})
	return b
}

// circleButton 是详情页左上角的返回键。
func circleButton(c *ui.Context, p Palette, name string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Size(38, 38).Padding(0).Radius(999).
		Background(p.Surface).Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Icon(c, icon(name)).FontSize(17).TextColor(p.Text)
	})
	return b
}

// chip 是「纯文本」这类带图标的小胶囊按钮。
func chip(c *ui.Context, p Palette, iconName, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(32).Padding(0, 12).Radius(8).
		Background(p.Surface).Border(1, p.Border).
		Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover).Border(1, ui.RGBA(124, 126, 44, 0.45))
	}
	b.Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			if label != "" {
				ui.Text(c, label).FontSize(11).FontWeight(700).LetterSpacing(0.6).SingleLine()
			}
			ui.Icon(c, icon(iconName)).FontSize(15).TextColor(p.Text)
		})
	})
	return b
}

// ghostButton 是行尾那种无底色的图标按钮。
func ghostButton(c *ui.Context, p Palette, name string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Size(28, 28).Padding(0).Radius(7).
		Cursor(ui.CursorPointer).OnClick(onClick)
	tint := p.Muted
	if b.Hovered() {
		b.Background(ui.RGBA(124, 126, 44, 0.18))
		tint = p.Accent
	}
	b.Children(func() {
		ui.Icon(c, icon(name)).FontSize(15).TextColor(tint)
	})
	return b
}

// card 是设置页那种分组卡片，行与行之间用分割线。
func card(c *ui.Context, p Palette) ui.Element {
	return ui.Column(c).Padding(2, 14).Radius(12).
		Background(p.Surface).Border(1, p.Border)
}

func divider(c *ui.Context, p Palette) {
	ui.Box(c).Height(1).FillWidth().Background(p.Border)
}

func sectionLabel(c *ui.Context, p Palette, text string) {
	ui.Text(c, text).FontSize(12).FontWeight(700).TextColor(p.Subtle)
}

// spacer 只是撑开一段固定高度：ui.Spacer 会吃掉剩余空间，这里不要它。
func spacer(c *ui.Context, h float32) {
	ui.Box(c).Height(h).FillWidth()
}

// dialogButton 是弹窗里的纯文字按钮。
func dialogButton(c *ui.Context, p Palette, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(32).Padding(0, 12).Radius(8).
		Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(12.5).TextColor(p.Text)
	})
	return b
}

// dialogIconButton 是弹窗里带图标的描边按钮。
func dialogIconButton(c *ui.Context, p Palette, iconName, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(32).Padding(0, 12).Radius(8).
		Background(p.Surface).Border(1, p.Border).Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon(iconName)).FontSize(14).TextColor(p.Text)
			ui.Text(c, label).FontSize(12.5).TextColor(p.Text)
		})
	})
	return b
}

// primaryDialogButton 是弹窗里的主按钮。
func primaryDialogButton(c *ui.Context, p Palette, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(32).Padding(0, 14).Radius(8).
		Background(p.Accent).Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(ui.Hex("#8b8d38"))
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(12.5).TextColor(p.AccentText)
	})
	return b
}
