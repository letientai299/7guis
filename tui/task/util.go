package task

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ColorInvalid  = tcell.ColorRed
	ColorDisabled = tcell.ColorDarkGray
)

func centerScreen(widget tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(widget, width, 0, true).
				AddItem(nil, 0, 1, false),
			height, 0, true).
		AddItem(nil, 0, 1, false)
}
