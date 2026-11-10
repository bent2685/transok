package desktop

import (
	"time"

	"transok/backend/app"
	"transok/backend/services"
	"transok/locales"

	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
)

// enterSelect 进入多选模式，顺手清空上一轮的勾选。
func (a *App) enterSelect() {
	a.selectMode = true
	a.selected = map[string]bool{}
}

func (a *App) leaveSelect() {
	a.selectMode = false
	a.selected = map[string]bool{}
}

func (a *App) selectedCount() int {
	n := 0
	for _, f := range a.files {
		if a.selected[f.Id] {
			n++
		}
	}
	return n
}

func (a *App) allSelected() bool {
	if len(a.files) == 0 {
		return false
	}
	for _, f := range a.files {
		if !a.selected[f.Id] {
			return false
		}
	}
	return true
}

func (a *App) toggleSelected(id string) {
	if a.selected[id] {
		delete(a.selected, id)
		return
	}
	a.selected[id] = true
}

func (a *App) toggleAllSelected() {
	if a.allSelected() {
		a.selected = map[string]bool{}
		return
	}
	for _, f := range a.files {
		a.selected[f.Id] = true
	}
}

// deleteSelected 把勾选的条目从分享清单里去掉，磁盘上的文件不动。
func (a *App) deleteSelected() {
	kept := make([]services.FileInfo, 0, len(a.files))
	for _, f := range a.files {
		if !a.selected[f.Id] {
			kept = append(kept, f)
		}
	}
	a.files = kept
	a.deleteAsk = false
	a.leaveSelect()
	services.Storage().Set("share-list", a.files)

	if len(kept) == 0 {
		a.share = false
		services.Storage().Set("is-share", false)
		app.Gin().Stop()
	}
}

// selectBar 是多选模式下底部那排操作：删除、全选或取消全选，都是玻璃按钮。
func (a *App) selectBar(c *ui.Context, p Palette) {
	count := a.selectedCount()

	ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
		// 没勾任何一项时删除键是不可用的。
		del := ui.ButtonBase(c).Height(44).Padding(0, 18).Radius(999).
			Cursor(ui.CursorPointer)
		if count == 0 {
			del.Material(glass.Glass{})
		} else {
			del.Material(glass.Glass{Tint: p.Danger, Interactive: true}).
				OnClick(func() { a.deleteAsk = true })
		}
		del.Children(func() {
			ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
				tint := p.AccentText
				if count == 0 {
					tint = p.Muted
				}
				ui.Icon(c, icon("trash")).FontSize(15).TextColor(tint)
				ui.Textf(c, "%s (%d)", locales.SelectDelete(), count).FontSize(13).
					FontWeight(700).TextColor(tint)
			})
		})

		label := locales.SelectAll()
		if a.allSelected() {
			label = locales.SelectNone()
		}
		toggle := ui.ButtonBase(c).Height(44).Padding(0, 18).Radius(999).
			Material(glass.Glass{Interactive: true}).Cursor(ui.CursorPointer).
			OnClick(a.toggleAllSelected)
		toggle.Children(func() {
			ui.Row(c).Gap(7).AlignItems(ui.Center).Children(func() {
				name := "check"
				if a.allSelected() {
					name = "x"
				}
				ui.Icon(c, icon(name)).FontSize(15).TextColor(p.Text)
				ui.Text(c, label).FontSize(13).FontWeight(700).TextColor(p.Text)
			})
		})
	})
}

// selectBox 是每行右侧的圆形选框：未选中是空心描边，选中后填成橄榄色，
// 中间的勾随弹簧放大浮出。
func selectBox(c *ui.Context, p Palette, checked bool) {
	box := ui.Column(c).Size(20, 20).Radius(999).Center().PassThrough()

	t := clamp01(box.AnimateTo("select", boolFloat(checked), ui.Animation{
		Spring: ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3},
	}))
	box.Background(p.Canvas.Mix(p.Accent, t))
	box.BorderWidth(1.5 * (1 - t)).BorderColor(p.Muted.Mix(p.Accent, t))

	if t > 0.01 {
		box.Children(func() {
			mark := ui.Column(c).Size(20, 20).Center().PassThrough()
			mark.Opacity(t).Scale(0.55+0.45*t, 0.55+0.45*t)
			mark.Children(func() {
				ui.Icon(c, icon("check")).FontSize(12).TextColor(p.AccentText)
			})
		})
	}
}

// deleteDrawer 是删除前的确认，跟其他弹层一样走抽屉。
func (a *App) deleteDrawer(c *ui.Context, p Palette) {
	count := a.selectedCount()
	a.drawer(c, p, &a.deleteAsk, drawerSpec{
		icon:        "trash",
		title:       locales.DeleteTitle(count),
		description: locales.DeleteDesc(),
		footer: func() {
			ui.Spacer(c).Grow(1)
			dialogButton(c, p, locales.DialogCancel(), func() { a.deleteAsk = false })
			dialogIconButton(c, p, "trash", p.Danger, locales.SelectDelete(), a.deleteSelected)
		},
	})
}
