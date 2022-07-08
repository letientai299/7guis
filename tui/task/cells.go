package task

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// https://en.wikipedia.org/wiki/Box-drawing_character
// todo: use tview.Borders
const (
	charDots = '…'
)

func Cells() tview.Primitive {
	sh := NewSheet()
	return sh
}

type cellData struct {
	raw     string
	display string
	affects []string
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
	data    map[string]*cellData
}

func NewSheet() *Sheet {
	hint := strings.TrimSpace(`
Navigate: ←↑→↓, hklj, Home ^
Enter to turn on edit mode, then Enter to commit change
`)

	const cellWidth = 10

	s := &Sheet{
		Box: tview.NewBox(),

		input: tview.NewInputField().
			SetLabelWidth(0).
			SetFieldWidth(cellWidth),

		data: make(map[string]*cellData),

		hint:      hint,
		hintView:  tview.NewTextView().SetText(hint),
		offset:    [2]int{1, 1}, // unlike array, cell start with 1
		focused:   [2]int{1, 1},
		cellWidth: cellWidth,

		txtStyle:    tcell.Style{}.Foreground(tview.Styles.PrimaryTextColor),
		headerStyle: tcell.Style{}.Foreground(tview.Styles.TertiaryTextColor),
		focusStyle:  tcell.Style{}.Foreground(tview.Styles.TertiaryTextColor),
		hoverStyle:  tcell.Style{}.Foreground(tview.Styles.SecondaryTextColor),
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
	nRows = (h - 2 /*hint*/ - 2 /*last and incomplete row*/) / 2
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
				display := sh.getDisplay([2]int{r, c})
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
	if !sh.editing {
		screen.HideCursor()
		return
	}

	x, y, _, _ := sh.GetInnerRect()
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

	// calculate focused loc top left location
	dy := (loc[0]-sh.offset[0])*2 + 2
	dx := (loc[1]-sh.offset[1])*(sh.cellWidth+1) + c1w

	// draw the border using focus style
	for i := 0; i < sh.cellWidth; i++ {
		screen.SetContent(x+dx+1+i, y+dy, tview.Borders.HorizontalFocus, nil, style)
		screen.SetContent(x+dx+1+i, y+dy+2, tview.Borders.HorizontalFocus, nil, style)
	}

	screen.SetContent(x+dx, y+dy+1, tview.Borders.VerticalFocus, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy+1, tview.Borders.VerticalFocus, nil, style)

	screen.SetContent(x+dx, y+dy, tview.Borders.TopLeftFocus, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy, tview.Borders.TopRightFocus, nil, style)
	screen.SetContent(x+dx+sh.cellWidth+1, y+dy+2, tview.Borders.BottomRightFocus, nil, style)
	screen.SetContent(x+dx, y+dy+2, tview.Borders.BottomLeftFocus, nil, style)
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
				sh.update(sh.focused, raw)
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
			sh.input.SetText(sh.getRaw(sh.focused))

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
			sh.update(sh.focused, sh.input.GetText())
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

func (sh *Sheet) getRaw(cell [2]int) string {
	data, ok := sh.data[sh.cellName(cell)]
	if !ok {
		return ""
	}
	return data.raw
}

func (sh *Sheet) update(cell [2]int, raw string) {
	raw = strings.TrimSpace(raw)

	name := sh.cellName(cell)
	data, ok := sh.data[name]
	if !ok {
		data = &cellData{
			raw:     raw,
			display: raw, // TODO: calculate
		}
	} else {
		data.raw = raw
		data.display = raw
	}

	if data.raw == "" {
		delete(sh.data, name)
	} else {
		sh.data[name] = data
	}
}

func (sh *Sheet) getDisplay(cell [2]int) string {
	data, ok := sh.data[sh.cellName(cell)]
	if !ok {
		return ""
	}
	return data.display
}
