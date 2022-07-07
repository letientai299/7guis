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

func enable(p tview.Primitive) {
	var b *tview.Box
	switch x := p.(type) {
	case *tview.InputField:
		x.SetLabelColor(tview.Styles.PrimaryTextColor)
		b = x.Box
	case *tview.Button:
		x.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
		b = x.Box
	}

	if b != nil {
		b.SetFocusFunc(nil)
	}
}

func disable(p tview.Primitive) {
	var b *tview.Box
	switch x := p.(type) {
	case *tview.InputField:
		x.SetLabelColor(ColorDisabled)
		b = x.Box
	case *tview.Button:
		x.SetBackgroundColor(ColorDisabled)
		b = x.Box
	}

	if b != nil {
		b.SetFocusFunc(func() {
			b.Blur()
		})
	}
}
