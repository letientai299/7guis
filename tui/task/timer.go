package task

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Timer(app *tview.Application) tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true)

	total := time.Second * 10
	p := NewProgress("Elapsed Time", float32(total))
	p.SetFormatter(func(value, max float32) string {
		sec := value / float32(time.Second)
		return fmt.Sprintf("%.1fs", sec)
	})
	container.AddItem(p, 3, 0, false)

	go func() {
		timerUpdate(app, p)
	}()

	slider := NewSlider("Duration",
		int64(time.Second),
		int64(total),
		int64(20*time.Second),
	).SetFormatter(func(v int64) string {
		return time.Duration(v).String()
	}).SetChangedFunc(func(v int64) {
		p.SetMax(float32(v))
	})

	container.AddItem(slider, 2, 0, true)
	container.AddItem(
		tview.NewTextView().SetText("Adjust duration with Left/Right or Mouse click"),
		0, 2, false,
	)

	btn := tview.NewButton("Press Enter to reset timer").
		SetSelectedFunc(func() {
			p.SetValue(0)
		})

	container.AddItem(btn, 1, 0, true)
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			p.SetValue(0)
			return nil
		}
		return event
	})

	return centerScreen(container, 50, 10)
}

func timerUpdate(app *tview.Application, p *Progress) {
	dur := time.Millisecond * 100
	for range time.Tick(dur) {
		v := time.Duration(p.GetValue())
		v += dur
		max := p.GetMax()
		total := time.Duration(max)
		if v > total {
			v = total
		}

		app.QueueUpdateDraw(func() { p.SetValue(float32(v)) })
	}
}
