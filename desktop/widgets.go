package desktop

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// 这是几个页面共用的视觉零件，尺寸对齐迁移前 Web 端的写法。
// 需要自己控制底色边框的地方一律用 ui.ButtonBase：ui.Button 会每帧按主题
// 重刷底色，压不住。

// glassGroup 是液态玻璃风格的按钮组：整块玻璃胶囊，里面并排放按钮，
// 悬停时只有被指到的那一格浮出一层薄光。
func glassGroup(c *ui.Context, p Palette) ui.Element {
	return ui.Row(c).Padding(3).Gap(2).Radius(14).Material(glass.Glass{})
}

// glassSegment 是按钮组里的一格：只放图标。
func glassSegment(c *ui.Context, p Palette, name string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Size(32, 32).Padding(0).Radius(11).
		Cursor(ui.CursorPointer).OnClick(onClick)
	glassHover(b, p)
	b.Children(func() {
		ui.Icon(c, icon(name)).FontSize(16).TextColor(p.Text)
	})
	return b
}

// glassChip 是按钮组里带文字的一格，用于操作型入口。
func glassChip(c *ui.Context, p Palette, iconName, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(32).Padding(0, 10).Radius(11).
		Cursor(ui.CursorPointer).OnClick(onClick)
	glassHover(b, p)
	b.Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			if label != "" {
				ui.Text(c, label).FontSize(11).FontWeight(700).LetterSpacing(0.6).
					TextColor(p.Text).SingleLine()
			}
			ui.Icon(c, icon(iconName)).FontSize(15).TextColor(p.Text)
		})
	})
	return b
}

// glassHover 给玻璃上的按钮加上按下/悬停时那层薄光。
func glassHover(b ui.Element, p Palette) {
	switch {
	case b.Pressed():
		b.Background(glassLift(p, 0.16))
	case b.Hovered():
		b.Background(glassLift(p, 0.10))
	}
}

// glassLift 是玻璃上叠的那层光：深色模式提亮，浅色模式压暗。
func glassLift(p Palette, a float32) ui.Color {
	if p.Dark {
		return ui.RGBA(255, 255, 255, a)
	}
	return ui.RGBA(0, 0, 0, a*0.7)
}

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

// closeButton 是抽屉右上角的关闭键：圆形底色，悬停提亮。
func closeButton(c *ui.Context, p Palette, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Size(28, 28).Padding(0).Radius(999).
		Background(p.SurfaceHover).Cursor(ui.CursorPointer).OnClick(onClick)
	tint := p.Muted
	if b.Hovered() {
		b.Background(p.Pressed)
		tint = p.Text
	}
	b.Children(func() {
		ui.Icon(c, icon("x")).FontSize(14).TextColor(tint)
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

// withAlpha 给颜色换一个透明度，做渐隐的分段。
func withAlpha(c ui.Color, a float32) ui.Color {
	return ui.RGBA(c.R, c.G, c.B, a)
}

// spacer 只是撑开一段固定高度：ui.Spacer 会吃掉剩余空间，这里不要它。
func spacer(c *ui.Context, h float32) {
	ui.Box(c).Height(h).FillWidth()
}

// labeledField 是带背景的输入字段：标签落在左上角，输入区透明贴在它下面，
// 获得焦点时外框转成橄榄色。
func labeledField(c *ui.Context, p Palette, label string, focused bool, input func()) ui.Element {
	field := ui.Column(c).FillWidth().Padding(9, 12, 10, 12).Radius(10).
		Background(p.Canvas).Border(1, p.Border)
	if focused {
		field.Border(1, p.Accent)
	}
	field.Children(func() {
		ui.Text(c, label).FontSize(10.5).FontWeight(700).LetterSpacing(0.5).TextColor(p.Muted)
		spacer(c, 5)
		input()
	})
	return field
}

// dialogButton 是抽屉里的无底色文字按钮。
func dialogButton(c *ui.Context, p Palette, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(38).Padding(0, 14).Radius(10).
		Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(13.5).TextColor(p.Text)
	})
	return b
}

// dialogIconButton 是抽屉里带图标的描边按钮，图标可以单独着色。
func dialogIconButton(c *ui.Context, p Palette, iconName string, iconColor ui.Color, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(38).Padding(0, 14).Radius(10).
		Background(p.Elevated).Border(1, p.Border).Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon(iconName)).FontSize(15).TextColor(iconColor)
			ui.Text(c, label).FontSize(13.5).TextColor(p.Text)
		})
	})
	return b
}

// dialogGhostButton 是抽屉里带图标的纯文字按钮，样式同无底色的取消键。
func dialogGhostButton(c *ui.Context, p Palette, iconName string, iconColor ui.Color, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(38).Padding(0, 14).Radius(10).
		Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(p.SurfaceHover)
	}
	b.Children(func() {
		ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon(iconName)).FontSize(15).TextColor(iconColor)
			ui.Text(c, label).FontSize(13.5).TextColor(p.Text)
		})
	})
	return b
}

// primaryDialogButton 是对话框里的主按钮。
func primaryDialogButton(c *ui.Context, p Palette, label string, onClick func()) ui.Element {
	b := ui.ButtonBase(c).Height(38).Padding(0, 16).Radius(10).
		Background(p.Accent).Cursor(ui.CursorPointer).OnClick(onClick)
	if b.Hovered() {
		b.Background(ui.Hex("#8b8d38"))
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(13.5).TextColor(p.AccentText)
	})
	return b
}
