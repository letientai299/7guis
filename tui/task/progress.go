package task

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewProgress(label string, max float32) *Progress {
	return &Progress{
		Mutex: &sync.Mutex{},
		Box:   tview.NewBox(),
		max:   max,
		label: label,
		formatValue: func(value, max float32) string {
			percent := value * 100 / max
			return fmt.Sprintf("%0.2f%%", percent)
		},
	}
}

type Progress struct {
	*tview.Box
	*sync.Mutex

	value       float32
	max         float32
	label       string
	formatValue func(value, max float32) string
}

func (p *Progress) SetMax(max float32) *Progress {
	p.Lock()
	defer p.Unlock()
	p.max = max
	return p
}

func (p *Progress) SetFormatter(fn func(value, max float32) string) *Progress {
	p.Lock()
	defer p.Unlock()
	p.formatValue = fn
	return p
}

func (p *Progress) SetValue(value float32) *Progress {
	p.Lock()
	defer p.Unlock()
	p.value = value
	return p
}

func (p *Progress) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)
	x, y, width, _ := p.GetInnerRect()
	labelLine := fmt.Sprintf(`%s (%s)`, p.label, p.formatValue(p.value, p.max))
	tview.Print(
		screen, labelLine, x, y, width, tview.AlignLeft,
		tview.Styles.PrimaryTextColor,
	)

	done := int(math.Floor(float64(p.value * float32(width) / p.max)))
	if done > width {
		done = width
	}

	progressLine := strings.Repeat("▆", done) + strings.Repeat("▁", width-done)
	tview.Print(
		screen, progressLine, x, y+1, width, tview.AlignLeft,
		tview.Styles.PrimaryTextColor,
	)
}

func (p *Progress) GetMax() float32 {
	p.Lock()
	defer p.Unlock()
	return p.max
}

func (p *Progress) GetValue() float32 {
	p.Lock()
	defer p.Unlock()
	return p.value
}
