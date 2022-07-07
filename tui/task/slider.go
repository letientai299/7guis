package task

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Slider struct {
	*tview.Box
	label       string
	min         int64
	max         int64
	val         int64
	formatValue func(value int64) string
	step        int64
	changed     func(int64)
}

func NewSlider(label string, min, val, max int64) *Slider {
	unit := (max - min) / 100
	if unit == 0 {
		unit = 1
	}

	return &Slider{
		Box:         tview.NewBox(),
		min:         min,
		val:         val,
		max:         max,
		step:        unit,
		label:       label,
		formatValue: func(value int64) string { return strconv.FormatInt(value, 10) },
		changed:     func(int64) {},
	}
}

func (p *Slider) SetChangedFunc(fn func(v int64)) *Slider {
	p.changed = fn
	return p
}

func (p *Slider) SetValue(v int64) *Slider {
	p.val = v
	if p.changed != nil {
		p.changed(v)
	}

	return p
}

func (p *Slider) SetFormatter(fn func(v int64) string) *Slider {
	p.formatValue = fn
	return p
}

func (p *Slider) Draw(screen tcell.Screen) {
	var lineChar = "═"
	if p.HasFocus() {
		lineChar = "▬"
	}

	const btn = "█"
	p.Box.DrawForSubclass(screen, p)
	x, y, width, _ := p.GetInnerRect()
	labelLine := fmt.Sprintf(`%s (%s)`, p.label, p.formatValue(p.val))
	tview.Print(
		screen, labelLine, x, y, width, tview.AlignLeft,
		tview.Styles.PrimaryTextColor,
	)

	left := int(
		math.Floor(float64(float32(p.val-p.min)/float32(p.max-p.min)*
			float32(width))),
	) + 1
	if left > width {
		left = width
	}

	leftLine := strings.Repeat(lineChar, left-1) + btn
	tview.Print(
		screen, leftLine, x, y+1, width, tview.AlignLeft,
		tview.Styles.TertiaryTextColor,
	)

	rightLine := strings.Repeat(lineChar, width-left)
	tview.Print(
		screen, rightLine, x, y+1, width, tview.AlignRight,
		tview.Styles.PrimaryTextColor,
	)
}

func (p *Slider) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		oldVal := p.val
		switch event.Key() {
		case tcell.KeyLeft:
			p.val -= p.step
			if p.val < p.min {
				p.val = p.min
			}
		case tcell.KeyRight:
			p.val += p.step
			if p.val > p.max {
				p.val = p.max
			}
		}

		if oldVal != p.val {
			p.changed(p.val)
		}
	})
}

func (p *Slider) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return p.WrapMouseHandler(
		func(
			action tview.MouseAction,
			event *tcell.EventMouse,
			setFocus func(p tview.Primitive),
		) (consumed bool, capture tview.Primitive) {
			x, y := event.Position()
			if !p.InRect(x, y) {
				return false, nil
			}
			setFocus(p)
			rx, _, w, _ := p.GetInnerRect()

			if (action == tview.MouseMove && event.Buttons() == tcell.ButtonPrimary) ||
				action == tview.MouseLeftClick {
				left := x - rx + 1
				f := math.Round(float64(left) / float64(w) * float64(p.max-p.min))
				p.val = int64(f) + p.min
				p.changed(p.val)
				return true, p
			}

			return false, nil
		})
}
