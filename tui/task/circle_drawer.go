package task

import (
	"github.com/rivo/tview"
)

func CircleDrawer() tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	controls := tview.NewFlex().
		AddItem(tview.NewButton("Undo"), 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(tview.NewButton("Redo"), 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(tview.NewButton("Clear"), 0, 1, false)

	container.AddItem(controls, 1, 0, true)
	cc := NewCircleCanvas()
	container.AddItem(cc, 0, 1, true)
	return container
}
