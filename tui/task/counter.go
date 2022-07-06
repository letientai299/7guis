package task

import (
	"fmt"

	"github.com/rivo/tview"
)

func Counter() tview.Primitive {
	counter := tview.NewFlex()

	count := 0
	txt := tview.NewTextView().
		SetText(fmt.Sprintf("Count is %d", count))

	label := "Count"
	button := tview.NewButton(label).SetSelectedFunc(func() {
		count++
		txt.SetText(fmt.Sprintf("Count is %d", count))
	})

	counter.AddItem(txt, 0, 1, false)
	counter.AddItem(button, 0, 1, true)
	counter.SetBorder(true)

	return centerScreen(counter, 30, 3)
}
