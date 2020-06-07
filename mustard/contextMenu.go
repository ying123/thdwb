package mustard

import (
	"image"
	assets "thdwb/assets"
	gg "thdwb/gg"

	"github.com/goki/freetype/truetype"
)

func (window *Window) EnableContextMenus() {
	window.contextMenu = &contextMenu{
		overlay: &Overlay{
			ref: "contextMenu",
		},
	}
}

func (window *Window) AddContextMenuEntry(entryText string, action func()) {
	window.contextMenu.entries = append(
		window.contextMenu.entries,
		&menuEntry{
			entryText: entryText,
			action:    action,
		},
	)
}

func (window *Window) DestroyContextMenu() {
	window.RemoveOverlay(window.contextMenu.overlay)
	window.contextMenu.entries = nil
}

func prepEntry(ctx *gg.Context, entry string, width float64) string {
	w, _ := ctx.MeasureString(entry)

	if w < width {
		return entry
	}

	for i := 0; i < len(entry); i++ {
		nW, _ := ctx.MeasureString(entry[:len(entry)-i] + "...")

		if nW <= width {
			return entry[:len(entry)-i] + "..."
		}
	}

	return entry
}

func (window *Window) DrawContextMenu() {
	menuWidth := float64(200)
	menuHeight := float64(len(window.contextMenu.entries) * 20)

	menuTop := float64(window.cursorY)
	menuLeft := float64(window.cursorX)

	if menuLeft+menuWidth > float64(window.width) {
		menuLeft = float64(window.width) - menuWidth
	}

	if menuTop+menuHeight > float64(window.height) {
		menuTop = float64(window.height) - menuHeight
	}

	ctx := gg.NewContext(int(menuWidth), int(menuHeight))
	ctx.DrawRectangle(0, 0, menuWidth, menuHeight)
	ctx.SetHexColor("#eee")
	ctx.Fill()

	font, _ := truetype.Parse(assets.OpenSans(400))
	ctx.SetHexColor("#222")
	ctx.SetFont(font, 16)

	for idx, entry := range window.contextMenu.entries {
		ctx.DrawString(prepEntry(ctx, entry.entryText, menuWidth), 0, 16+float64(idx*20))
		ctx.Fill()
	}

	ctx.DrawRectangle(0, 0, menuWidth, menuHeight)
	ctx.SetHexColor("#ddd")
	ctx.Stroke()

	overlay := extractOverlay(
		ctx.Image().(*image.RGBA),
		image.Point{
			int(menuLeft),
			int(menuTop),
		})

	window.SetContextMenuOverlay(overlay)
}

func (window *Window) SetContextMenuOverlay(overlay *Overlay) {
	window.contextMenu.overlay = overlay
	window.AddOverlay(overlay)
}

func extractOverlay(buffer *image.RGBA, postion image.Point) *Overlay {
	return &Overlay{
		ref:    "contextMenu",
		active: true,

		top:  float64(postion.Y),
		left: float64(postion.X),

		position: postion,
		buffer:   buffer,
	}
}