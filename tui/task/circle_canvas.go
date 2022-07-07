package task

import "C"
import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const characterBlockRatio = 0.45 // ratio between height and width of a character
const (
	minRadius = 3
	maxRadius = 60
)

type point [2]int

type circle struct {
	x, y, r int
}

func (c circle) center() point { return [2]int{c.x, c.y} }

// pointsOnCircle lists integer points on the circle
func (c *circle) pointsOnCircle() []point {
	var res []point
	for py := -c.r; py <= c.r; py++ {
		px := int(math.Round(math.Sqrt(float64(c.r*c.r - py*py))))
		y := int(math.Round(float64(py) * characterBlockRatio))
		res = append(res, [2]int{c.x + px, c.y + y})
		res = append(res, [2]int{c.x - px, c.y - y})
	}

	return res
}

func NewCircleCanvas() *CircleCanvas {
	return &CircleCanvas{
		Box:       tview.NewBox(),
		mCircles:  make(map[point]*circle),
		centerDot: '⬤',
		dot:       '🞄',
	}
}

// CircleCanvas limitations:
// - There can only be a single circle for any center point, no matter what the
// radius is
type CircleCanvas struct {
	*tview.Box
	mouseX, mouseY int
	circles        []*circle
	mCircles       map[point]*circle
	selected       point
	centerDot      rune
	dot            rune
	slider         *Slider
	onAddCircle    func(circle2 circle)
}

func (cc *CircleCanvas) Draw(screen tcell.Screen) {
	cc.Box.DrawForSubclass(screen, cc)
	if cc.HasFocus() {
		if cc.mouseX > 0 && cc.mouseY > 0 {
			x, y, w, _ := cc.GetInnerRect()
			mouseInfo := fmt.Sprintf("(%d,%d)", cc.mouseX, cc.mouseY)
			tview.Print(screen, mouseInfo, x, y, w, tview.AlignLeft, tview.Styles.PrimaryTextColor)
			screen.SetContent(cc.mouseX+x, cc.mouseY+y, cc.centerDot, nil, tcell.Style{}.Foreground(tview.Styles.SecondaryTextColor))
		}
	}

	var special *circle
	for _, circle := range cc.circles {
		if circle.center() == cc.selected {
			special = circle
		} else {
			cc.draw(screen, circle)
		}
	}

	// draw selected as last for full visibility
	if special != nil {
		cc.draw(screen, special)
	}

	if cc.slider != nil {
		cc.slider.Draw(screen)
	}
}

func (cc *CircleCanvas) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return cc.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		mx, my := event.Position()
		if !cc.InRect(mx, my) {
			return false, nil
		}

		if cc.slider != nil && cc.slider.InRect(mx, my) {
			cc.mouseX = -1
			cc.mouseY = -1
			return cc.slider.MouseHandler()(action, event, setFocus)
		}

		rectX, rectY, _, _ := cc.GetInnerRect()
		cc.mouseX = mx - rectX
		cc.mouseY = my - rectY
		if action == tview.MouseLeftClick {
			if cc.slider != nil {
				cc.removeSelected()
				return true, cc
			}

			if !cc.HasFocus() {
				setFocus(cc)
			} else {
				cc.checkAndAddCircle(cc.mouseX, cc.mouseY)
			}
		}

		return true, cc
	})
}

func (cc *CircleCanvas) checkAndAddCircle(x int, y int) {
	c := &circle{x: x, y: y, r: 10}

	center := [2]int{x, y}
	if cc.mCircles[center] == nil {
		cc.circles = append(cc.circles, c)
		cc.mCircles[center] = c
		cc.removeSelected()
		return
	}

	// has a circle at the selected center already
	cc.selected = center
	cc.slider = NewSlider("Radius", minRadius, int64(c.r), maxRadius)
	rx, ry, _, _ := cc.GetInnerRect()

	sliderW := maxRadius - minRadius
	cc.slider.SetRect(rx+center[0]-sliderW/2, ry+center[1]+3, sliderW, 3)
	cc.slider.SetChangedFunc(func(v int64) {
		cc.mCircles[center].r = int(v)
	})
}

func (cc *CircleCanvas) removeSelected() {
	cc.selected = [2]int{-1, -1} // clear selected
	cc.slider = nil
}

func (cc *CircleCanvas) draw(screen tcell.Screen, c *circle) {
	x, y, w, h := cc.GetInnerRect()
	color := tview.Styles.PrimaryTextColor
	if c.center() == cc.selected {
		color = tview.Styles.TertiaryTextColor
	}

	style := tcell.Style{}.Foreground(color)
	screen.SetContent(c.x+x, c.y+y, cc.centerDot, nil, style)

	points := c.pointsOnCircle()
	for _, p := range points {
		px, py := p[0], p[1] // relative coordinate within canvas
		if px > w-1 || py > h-1 || px < 0 || py < 0 {
			continue
		}

		screen.SetContent(px+x, py+y, cc.dot, nil, style)
	}
}

func (cc *CircleCanvas) Focus(delegate func(p tview.Primitive)) {
	if cc.slider != nil {
		delegate(cc.slider)
	} else {
		delegate(cc.Box)
	}
}

func (cc *CircleCanvas) HasFocus() bool {
	if cc.slider != nil {
		return cc.slider.HasFocus()
	}

	return cc.Box.HasFocus()
}

func (cc *CircleCanvas) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cc.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if cc.slider != nil {
			cc.slider.InputHandler()(event, setFocus)
		}
	})
}
