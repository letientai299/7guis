package task

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/maja42/goval"
	"github.com/rivo/tview"
)

const (
	charDots = '…'
)

var regexCellName = regexp.MustCompile(`\b[A-Z]+\d+\b`)

func Cells() tview.Primitive {
	sh := NewSheet()
	return sh
}

type cellData struct {
	raw     string
	display string
	value   float64
	affects map[string]struct{}
	needs   map[string]struct{}
	err     error
}

func newCellData() *cellData {
	return &cellData{
		affects: make(map[string]struct{}),
		needs:   make(map[string]struct{}),
	}
}

// Sheet supports up to 26 columns since I'm too lazy to write the code for
// column name.
type Sheet struct {
	*tview.Box
	hint      string
	hintView  *tview.TextView
	cellWidth int // amount of character per cell

	txtStyle    tcell.Style
	headerStyle tcell.Style
	focusStyle  tcell.Style
	hoverStyle  tcell.Style

	offset  [2]int // coordinate of the top-left visible cell
	focused [2]int
	hovered [2]int

	editing bool
	input   *tview.InputField
	mem     map[string]*cellData
}

func NewSheet() *Sheet {
	hint := strings.TrimSpace(`
Navigate: ←↑→↓, hklj, Home ^
Enter to turn on edit mode, then Enter to commit change
`)

	const cellWidth = 15

	st := tview.Styles
	s := &Sheet{
		Box: tview.NewBox(),

		input: tview.NewInputField().
			SetLabelWidth(0).
			SetFieldWidth(cellWidth),

		mem: make(map[string]*cellData),

		hint:      hint,
		hintView:  tview.NewTextView().SetText(hint),
		offset:    [2]int{1, 1}, // unlike array, cell start with 1
		focused:   [2]int{1, 1},
		cellWidth: cellWidth,

		txtStyle:    tcell.Style{}.Foreground(st.PrimaryTextColor).Background(st.PrimitiveBackgroundColor),
		headerStyle: tcell.Style{}.Foreground(st.TertiaryTextColor).Background(st.PrimitiveBackgroundColor),
		focusStyle:  tcell.Style{}.Foreground(st.TertiaryTextColor).Background(st.PrimitiveBackgroundColor),
		hoverStyle:  tcell.Style{}.Foreground(st.SecondaryTextColor).Background(st.PrimitiveBackgroundColor),
	}

	return s
}

func (sh *Sheet) Draw(screen tcell.Screen) {
	sh.Box.DrawForSubclass(screen, sh)
	nRows, nCols, c1w := sh.viewport()
	sh.adjustViewport(nRows, nCols)
	sh.drawSheet(screen, nRows, nCols, c1w)
	sh.drawHovered(screen, nRows, nCols, c1w)
	sh.drawFocusedCell(screen, c1w)
	sh.drawHint(screen)
}

func (sh *Sheet) viewport() (nRows int, nCols int, c1w int) {
	_, _, w, h := sh.GetInnerRect()
	nRows = (h -
		2 /*hint*/ -
		2 /*last and incomplete row*/ -
		1 /*focused cell detail*/) / 2
	offset := sh.offset
	// 2 vertical line and 1 reserved space, in case, e.g. move from 99 to 100
	c1w = len(strconv.Itoa(offset[0]+nRows)) + 3
	nCols = (w - c1w) / (sh.cellWidth + 1)
	return
}

func (sh *Sheet) adjustViewport(nRows int, nCols int) {
	offset := sh.offset
	if sh.focused[0] < offset[0] {
		offset[0] = sh.focused[0]
	}

	if sh.focused[1] < offset[1] {
		offset[1] = sh.focused[1]
	}

	if sh.focused[0]-offset[0] > nRows-1 {
		offset[0] = sh.focused[0] - nRows + 1
	}

	if sh.focused[1]-offset[1] > nCols-1 {
		offset[1] = sh.focused[1] - nCols + 1
	}

	sh.offset = offset
}

func (sh *Sheet) drawSheet(
	screen tcell.Screen,
	nRows int, // num visible rows, including the header row
	nCols int, // num visible columns, including the header col
	c1w int, // width of header col
) {
	x, y, w, _ := sh.GetInnerRect()
	y++
	for r := 0; r < nRows+1; r++ {
		if r != 0 {
			screen.SetContent(x, 2*r+y, tview.Borders.LeftT, nil, sh.txtStyle)
		} else {
			screen.SetContent(x, 2*r+y, tview.Borders.TopLeft, nil, sh.txtStyle)
		}

		screen.SetContent(x, 2*r+y+1, tview.Borders.Vertical, nil, sh.txtStyle)
		if r != 0 {
			rowName := strconv.Itoa(r - 1 + sh.offset[0])
			sh.drawTxtAlignRight(screen, rowName, sh.headerStyle, x, 2*r+y+1, c1w)
		}
		screen.SetContent(x+c1w, 2*r+y+1, tview.Borders.Vertical, nil, sh.txtStyle)

		// line
		for i := 1; i < w; i++ {
			c := tview.Borders.Horizontal
			if (i-c1w)%(sh.cellWidth+1) == 0 {
				if r != 0 {
					c = tview.Borders.Cross
				} else {
					c = tview.Borders.TopT
				}
			}
			screen.SetContent(x+i, 2*r+y, c, nil, sh.txtStyle)
		}

		// content
		for c := 1; c < nCols+1; c++ {
			pos := x + c1w + (c-1)*(sh.cellWidth+1)
			if r == 0 {
				colName := sh.colName(c - 1 + sh.offset[1])
				sh.drawTxtAlignRight(screen, colName, sh.headerStyle, pos+1, 2*r+y+1, sh.cellWidth)
			} else {
				cell := [2]int{r - 1 + sh.offset[0], c - 1 + sh.offset[1]}
				display := sh.getCellDisplay(cell)
				if display != "" {
					_, err := strconv.ParseFloat(display, 64)
					if err == nil {
						// draw number align right
						sh.drawTxtAlignRight(screen, display, sh.txtStyle, pos+1, 2*r+y+1, sh.cellWidth)
					} else {
						sh.drawTxt(screen, display, sh.txtStyle, pos+1, 2*r+y+1, sh.cellWidth)
					}
				}
			}

			screen.SetContent(x+c1w+c*(sh.cellWidth+1), 2*r+y+1, tview.Borders.Vertical, nil, sh.txtStyle)
		}

		screen.SetContent(x+w-1, 2*r+y+1, charDots, nil, sh.txtStyle)
	}

	// last line
	screen.SetContent(x, 2*nRows+y+2, tview.Borders.LeftT, nil, sh.txtStyle)
	for i := 1; i < w; i++ {
		c := tview.Borders.Horizontal
		if (i-c1w)%(sh.cellWidth+1) == 0 {
			c = tview.Borders.Cross
		}
		screen.SetContent(x+i, 2*nRows+y+2, c, nil, sh.txtStyle)
	}
}

func (sh *Sheet) drawFocusedCell(screen tcell.Screen, c1w int) {
	sh.highlightCell(screen, sh.focused, c1w, sh.focusStyle)

	x, y, _, _ := sh.GetInnerRect()
	name := sh.cellName(sh.focused)
	if len(name) < c1w {
		name = strings.Repeat(" ", c1w-len(name)) + name
	}
	name += " = "
	var detail string
	if sh.editing {
		detail = sh.getCell(sh.focused).raw
	} else {
		detail = sh.getCellDisplay(sh.focused)
	}

	for i, c := range name {
		if c != ' ' {
			screen.SetContent(x+i, y, c, nil, sh.focusStyle)
		}
	}

	for i, c := range detail {
		if c != ' ' {
			screen.SetContent(x+i+len(name), y, c, nil, sh.txtStyle)
		}
	}

	if !sh.editing {
		screen.HideCursor()
		return
	}

	y++

	ix := x + (sh.focused[1]-sh.offset[1])*(sh.cellWidth+1) + c1w + 1
	iy := y + (sh.focused[0]-sh.offset[0]+1)*2 + 1
	sh.input.SetRect(ix, iy, sh.cellWidth, 1)
	sh.input.Draw(screen)
}

func (sh *Sheet) drawHovered(screen tcell.Screen, rows int, cols int, c1w int) {
	if !sh.HasFocus() ||
		sh.hovered[0] < sh.offset[0] ||
		sh.hovered[1] < sh.offset[1] ||
		sh.hovered[0] > sh.offset[0]+rows ||
		sh.hovered[1] > sh.offset[1]+cols {
		return
	}

	sh.highlightCell(screen, sh.hovered, c1w, sh.hoverStyle)
}

func (sh *Sheet) highlightCell(screen tcell.Screen, loc [2]int, c1w int, style tcell.Style) {
	x, y, _, _ := sh.GetInnerRect()
	y++

	// calculate focused loc top left location
	dy := (loc[0]-sh.offset[0])*2 + 2
	dx := (loc[1]-sh.offset[1])*(sh.cellWidth+1) + c1w

	// draw the border using focus style
	for i := 0; i < sh.cellWidth; i++ {

		screen.SetContent(x+dx+1+i, y+dy, tview.Borders.Horizontal, nil, style)
		screen.SetContent(x+dx+1+i, y+dy+2, tview.Borders.Horizontal, nil, style)
	}

	screen.SetContent(x+dx, y+dy+1, tview.Borders.Vertical, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy+1, tview.Borders.Vertical, nil, style)

	screen.SetContent(x+dx, y+dy, tview.Borders.TopLeft, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy, tview.Borders.TopRight, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy+2, tview.Borders.BottomRight, nil, style)
	screen.SetContent(x+dx, y+dy+2, tview.Borders.BottomLeft, nil, style)
}

func (sh Sheet) drawHint(screen tcell.Screen) {
	x, y, w, h := sh.GetInnerRect()
	sh.hintView.SetRect(x, y+h-2, w, 2)
	sh.hintView.Draw(screen)
}

func (sh *Sheet) drawTxt(
	screen tcell.Screen,
	s string,
	style tcell.Style,
	x int,
	y int,
	w int,
) {
	rs := []rune(s)
	if len(rs) > w { // offset for vertical line
		rs = append(rs[:w-1], charDots)
	}
	for i, c := range rs {
		screen.SetContent(x+i, y, c, nil, style)
	}
}

func (sh *Sheet) drawTxtAlignRight(
	screen tcell.Screen,
	s string,
	style tcell.Style,
	x int,
	y int,
	w int,
) {
	rs := []rune(s)
	if len(rs) > w { // offset for vertical line
		rs = append(rs[:w-1], charDots)
	}

	pad := w - len(rs)
	for i, c := range rs {
		screen.SetContent(x+pad+i, y, c, nil, style)
	}
}

func (sh *Sheet) InputHandler() func(e *tcell.EventKey, focus func(p tview.Primitive)) {
	return sh.WrapInputHandler(func(e *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if sh.editing {
			if e.Key() == tcell.KeyEnter {
				raw := sh.input.GetText()
				sh.updateCell(sh.focused, raw)
				sh.editing = false
				setFocus(sh)
				return
			}

			// let the input field handle it
			sh.input.InputHandler()(e, setFocus)
			setFocus(sh.input)
			return
		}

		k := e.Key()
		ru := e.Rune()

		rows, cols, _ := sh.viewport()

		switch {
		case k == tcell.KeyEnter:
			sh.editing = true
			sh.input.SetText(sh.getCell(sh.focused).raw)

		case k == tcell.KeyLeft || (k == tcell.KeyRune && ru == 'h'):
			if sh.focused[1] > 1 {
				sh.focused[1]--
			}
		case k == tcell.KeyRight || (k == tcell.KeyRune && ru == 'l'):
			sh.focused[1]++

		case k == tcell.KeyUp || (k == tcell.KeyRune && ru == 'k'):
			if sh.focused[0] > 1 {
				sh.focused[0]--
			}

		case k == tcell.KeyDown || (k == tcell.KeyRune && ru == 'j'):
			sh.focused[0]++

		// Move to top row
		case (k == tcell.KeyHome && e.Modifiers() == tcell.ModShift) ||
			(k == tcell.KeyRune && ru == 'g'):
			sh.focused[0] = 1

		// Move to left most column
		case k == tcell.KeyHome || (k == tcell.KeyRune && ru == '^'):
			sh.focused[1] = 1
		}

		sh.adjustViewport(rows, cols)
		setFocus(sh)
	})
}

func (sh *Sheet) MouseHandler() func(ac tview.MouseAction, e *tcell.EventMouse, focus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return sh.WrapMouseHandler(func(ac tview.MouseAction, e *tcell.EventMouse, f func(p tview.Primitive)) (bool, tview.Primitive) {
		mx, my := e.Position()
		if !sh.InRect(mx, my) {
			sh.hovered = [2]int{0, 0}
			return false, nil
		}

		rectX, rectY, _, _ := sh.GetInnerRect()
		rectY++

		mx -= rectX
		my -= rectY
		rows, cols, w := sh.viewport()
		r := ((my - 2) / 2) + sh.offset[0]
		c := ((mx - w) / sh.cellWidth) + sh.offset[1]

		if r == 0 || c == 0 || r > rows-1 || c > cols {
			sh.hovered = [2]int{0, 0}
			return false, nil
		}

		if e.Buttons() == tcell.ButtonPrimary {
			// commit the input text
			sh.updateCell(sh.focused, sh.input.GetText())
			sh.editing = false
			sh.focused = [2]int{r, c}
		} else {
			sh.hovered = [2]int{r, c}
		}

		return true, sh
	})
}

func (sh *Sheet) colName(n int) string {
	// https://leetcode.com/problems/excel-sheet-column-title/description/
	// https://github.com/letientai299/leetcode/blob/master/go/168.excel-sheet-column-title.go

	bf := bytes.Buffer{}
	for n != 0 {
		x := (n - 1) % 26
		bf.WriteByte(byte('A' + x))
		n = (n - x) / 26
	}

	a := bf.Bytes()
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return string(a)
}

func (sh *Sheet) cellName(cell [2]int) string {
	return sh.colName(cell[1]) + strconv.Itoa(cell[0])
}

func (sh *Sheet) HasFocus() bool {
	if sh.editing {
		return sh.input.HasFocus()
	}
	return sh.Box.HasFocus()
}

func (sh *Sheet) Focus(delegate func(p tview.Primitive)) {
	if sh.editing {
		delegate(sh.input)
	} else {
		delegate(sh.Box)
	}
}

func (sh *Sheet) getCell(cell [2]int) *cellData {
	data, ok := sh.mem[sh.cellName(cell)]
	if !ok {
		return newCellData()
	}
	return data
}

func (sh *Sheet) updateCell(cell [2]int, raw string) {
	raw = strings.TrimSpace(raw)

	current := sh.cellName(cell)
	data, ok := sh.mem[current]
	if !ok {
		data = newCellData()
	}

	defer func() {
		if data.raw == "" {
			delete(sh.mem, current)
		} else {
			sh.mem[current] = data
		}

		for other := range data.needs {
			o, ok := sh.mem[other]
			if !ok {
				o = newCellData()
			}
			o.affects[current] = struct{}{}
			sh.mem[other] = o
		}

		for other := range data.affects {
			sh.recompute(other)
		}
	}()

	data.raw = raw
	if !strings.HasPrefix(data.raw, "=") {
		data.display = raw
		data.err = nil
		data.needs = nil
		for other := range data.needs {
			if o, ok := sh.mem[other]; ok {
				delete(o.affects, current)
			}
		}
		sh.mem[current] = data
		return
	}

	newNeeds := sh.needCellNames(data.raw)
	for other := range data.needs {
		if _, stillAffect := newNeeds[other]; !stillAffect {
			o, ok := sh.mem[other]
			if ok {
				delete(o.affects, current)
			}
		}
	}

	data.needs = newNeeds
	sh.compute(current, data)
}

func (sh *Sheet) compute(name string, data *cellData) {
	// TODO (tai): this won't work with cell ranges
	replacements := make([]string, 0, len(data.needs)*2)
	data.err = nil

	for need := range data.needs {
		other, ok := sh.mem[need]

		var f float64
		if !ok || need == name {
			other = newCellData()
		}

		if other.value != 0 {
			f = other.value
		} else if other.display != "" {
			var err error
			f, err = strconv.ParseFloat(other.display, 64)
			if err != nil {
				data.err = fmt.Errorf("%s(%s) is NaN", need, other.display)
				break
			}
		}

		replacements = append(replacements, need, strconv.FormatFloat(f, 'g', -1, 64))
	}

	if data.err != nil {
		return
	}

	expStr := strings.NewReplacer(replacements...).
		Replace(data.raw)[1:] // remove prefix equal sign (=)
	result, err := goval.NewEvaluator().Evaluate(expStr, nil, nil)
	if err != nil {
		data.err = err
		return
	}

	switch v := result.(type) {
	case float64:
		data.value = v
		data.display = fmt.Sprintf("%g", v)
		data.err = nil
	case int:
		data.value = float64(v)
		data.display = strconv.Itoa(v)
		data.err = nil
	default:
		data.err = fmt.Errorf("%v is NaN", result)
	}
}

func (sh *Sheet) needCellNames(raw string) map[string]struct{} {
	rawNeeds := regexCellName.FindAllString(raw, -1)
	needs := make(map[string]struct{})
	for _, need := range rawNeeds {
		needs[need] = struct{}{}
	}

	return needs
}

func (sh *Sheet) recompute(cell string) {
	data, ok := sh.mem[cell]
	if !ok {
		data = newCellData()
	}
	sh.compute(cell, data)
	sh.mem[cell] = data
}

func (sh *Sheet) getCellDisplay(cell [2]int) string {
	c := sh.getCell(cell)
	if c.err != nil {
		return c.err.Error()
	}

	return c.display
}
