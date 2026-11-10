package desktop

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

// SVG 取自 tabler icons（MIT），与 Web 端的图标同名同源。
//
//go:embed icons/*.svg
var iconFS embed.FS

var iconCache = map[string]*ui.SVG{}

// icon 解析并缓存图标，同一个图标只解析一次；名字写错时退回通用文件图标。
func icon(name string) *ui.SVG {
	if svg, ok := iconCache[name]; ok {
		return svg
	}
	data, err := iconFS.ReadFile("icons/" + name + ".svg")
	if err != nil {
		data, _ = iconFS.ReadFile("icons/file.svg")
	}
	svg := ui.MustParseSVG(data)
	iconCache[name] = svg
	return svg
}
