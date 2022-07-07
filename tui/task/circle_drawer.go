package task

import (
	"github.com/rivo/tview"
)

func CircleDrawer() tview.Primitive {
	cd := &circleDrawer{}
	return cd.render()
}

type action struct {
	up   func()
	down func()
}

type circleDrawer struct {
	circles     []*circle
	actions     []action
	actionIndex int
	canvas      *CircleCanvas
}

func (cd *circleDrawer) render() tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	controls := tview.NewFlex().
		AddItem(tview.NewButton("Undo").SetSelectedFunc(cd.undo), 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(tview.NewButton("Redo").SetSelectedFunc(cd.redo), 0, 1, false)

	container.AddItem(controls, 3, 0, true)
	cd.canvas = NewCircleCanvas().
		SetOnUpdateCircle(cd.onUpdate).
		SetOnAddCircle(cd.onAdd)

	container.AddItem(cd.canvas, 0, 1, true)
	return container
}

func (cd *circleDrawer) undo() {
	if cd.actionIndex == 0 {
		return // can't undo
	}

	cd.actionIndex--
	a := cd.actions[cd.actionIndex]
	a.down()
	cd.canvas.setCircles(cd.circles)
}

func (cd *circleDrawer) redo() {
	if cd.actionIndex == len(cd.actions) {
		return // can't redo
	}

	a := cd.actions[cd.actionIndex]
	cd.actionIndex++
	a.up()
	cd.canvas.setCircles(cd.circles)
}

func (cd *circleDrawer) onUpdate(c *circle, oldR int, newR int) {
	a := action{
		up:   func() { c.r = newR },
		down: func() { c.r = oldR },
	}
	a.up()
	cd.addAction(a)
}

func (cd *circleDrawer) onAdd(c *circle) {
	a := action{
		up:   func() { cd.circles = append(cd.circles, c) },
		down: func() { cd.circles = cd.circles[:len(cd.circles)-1] },
	}
	a.up()
	cd.addAction(a)
}

func (cd *circleDrawer) addAction(a action) {
	if cd.actionIndex < len(cd.actions) {
		cd.actions[cd.actionIndex] = a
		cd.actionIndex++
		return
	}

	cd.actions = append(cd.actions, a)
	cd.actionIndex++
}
