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

// Sheet supports up to 26 columns since I'm too lazy to write the code for
// column name.
type Sheet struct {
	*tview.Box
	hint        string
	hintView    *tview.TextView
	cellWidth   int
	txtStyle    tcell.Style
	headerStyle tcell.Style
	focusStyle  tcell.Style
	offsetCell  [2]int // coordinate of the top-left visible cell
	focusedCell [2]int
}

func NewSheet() *Sheet {
	hint := strings.TrimSpace(`
Navigate: ←↑→↓, hklj, Home ^
Enter to turn on edit mode, then Enter to commit change
`)
	s := &Sheet{
		Box:         tview.NewBox(),
		hint:        hint,
		hintView:    tview.NewTextView().SetText(hint),
		offsetCell:  [2]int{1, 1}, // unlike array, cell start with 1
		focusedCell: [2]int{1, 1},
		cellWidth:   10,
		txtStyle:    tcell.Style{}.Foreground(tview.Styles.PrimaryTextColor),
		headerStyle: tcell.Style{}.Foreground(tview.Styles.SecondaryTextColor),
		focusStyle:  tcell.Style{}.Foreground(tview.Styles.TertiaryTextColor),
	}
	return s
}

func (sh *Sheet) Draw(screen tcell.Screen) {
	sh.Box.DrawForSubclass(screen, sh)
	nRows, nCols, c1w := sh.viewportInfo()
	sh.adjustViewport(nRows, nCols)
	sh.drawSheet(screen, nRows, nCols, c1w)
	sh.drawFocusedCell(screen, c1w)
	sh.drawHint(screen)
}

func (sh *Sheet) viewportInfo() (nRows int, nCols int, c1w int) {
	_, _, w, h := sh.GetInnerRect()
	bot := strings.Count(sh.hint, "\n")
	nRows = (h - bot /*hint*/ - 2 /*last and incomplete row*/) / 2
	offset := sh.offsetCell
	// 2 vertical line and 1 reserved space, in case, e.g. move from 99 to 100
	c1w = len(strconv.Itoa(offset[0]+nRows)) + 3
	nCols = (w - c1w) / (sh.cellWidth)
	return
}

func (sh *Sheet) adjustViewport(nRows int, nCols int) {
	offset := sh.offsetCell
	if sh.focusedCell[0] < offset[0] {
		offset[0] = sh.focusedCell[0]
	}

	if sh.focusedCell[1] < offset[1] {
		offset[1] = sh.focusedCell[1]
	}

	if sh.focusedCell[0]-offset[0] > nRows-1 {
		offset[0] = sh.focusedCell[0] - nRows + 1
	}

	if sh.focusedCell[1]-offset[1] > nCols-1 {
		offset[1] = sh.focusedCell[1] - nCols + 1
	}

	sh.offsetCell = offset
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
			rowName := strconv.Itoa(r - 1 + sh.offsetCell[0])
			sh.drawTxtAlignRight(screen, rowName, sh.headerStyle, x, 2*r+y+1, c1w)
		}
		screen.SetContent(x+c1w, 2*r+y+1, tview.Borders.Vertical, nil, sh.txtStyle)

		// line
		for i := 1; i < w; i++ {
			c := tview.Borders.Horizontal
			if (i-c1w)%(sh.cellWidth) == 0 {
				if r != 0 {
					c = tview.Borders.Cross
				} else {
					c = tview.Borders.TopT
				}
			}
			screen.SetContent(x+i, 2*r+y, c, nil, sh.txtStyle)
		}

		// content
		for c := 0; c < nCols; c++ {
			pos := x + c1w + c*sh.cellWidth
			if r == 0 {
				colName := sh.colName(c + sh.offsetCell[1])
				sh.drawTxtAlignRight(screen, colName, sh.headerStyle, pos, 2*r+y+1, sh.cellWidth)
			}
			screen.SetContent(x+c1w+(c+1)*sh.cellWidth, 2*r+y+1, tview.Borders.Vertical, nil, sh.txtStyle)
		}

		screen.SetContent(x+w-1, 2*r+y+1, charDots, nil, sh.txtStyle)
	}

	// last line
	screen.SetContent(x, 2*nRows+y+2, tview.Borders.LeftT, nil, sh.txtStyle)
	for i := 1; i < w; i++ {
		c := tview.Borders.Horizontal
		if (i-c1w)%(sh.cellWidth) == 0 {
			c = tview.Borders.Cross
		}
		screen.SetContent(x+i, 2*nRows+y+2, c, nil, sh.txtStyle)
	}
}

func (sh *Sheet) drawFocusedCell(screen tcell.Screen, c1w int) {
	x, y, _, _ := sh.GetInnerRect()

	// calculate focused cell top left location
	dy := (sh.focusedCell[0]-sh.offsetCell[0])*2 + 2
	dx := (sh.focusedCell[1]-sh.offsetCell[1])*sh.cellWidth + c1w

	// draw the border using focus style
	for i := 0; i < sh.cellWidth-1; i++ {
		screen.SetContent(x+dx+1+i, y+dy, tview.Borders.HorizontalFocus, nil, sh.focusStyle)
		screen.SetContent(x+dx+1+i, y+dy+2, tview.Borders.HorizontalFocus, nil, sh.focusStyle)
	}

	screen.SetContent(x+dx, y+dy+1, tview.Borders.VerticalFocus, nil, sh.focusStyle)
	screen.SetContent(x+dx+sh.cellWidth, y+dy+1, tview.Borders.VerticalFocus, nil, sh.focusStyle)

	screen.SetContent(x+dx, y+dy, tview.Borders.TopLeftFocus, nil, sh.focusStyle)
	screen.SetContent(x+dx+sh.cellWidth, y+dy, tview.Borders.TopRightFocus, nil, sh.focusStyle)
	screen.SetContent(x+dx+sh.cellWidth, y+dy+2, tview.Borders.BottomRightFocus, nil, sh.focusStyle)
	screen.SetContent(x+dx, y+dy+2, tview.Borders.BottomLeftFocus, nil, sh.focusStyle)
}

func (sh Sheet) drawHint(screen tcell.Screen) {
	x, y, w, h := sh.GetInnerRect()
	sh.hintView.SetRect(x, y+h-2, w, 2)
	sh.hintView.Draw(screen)
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
		k := e.Key()
		ru := e.Rune()

		rows, cols, _ := sh.viewportInfo()

		switch {
		case k == tcell.KeyLeft || (k == tcell.KeyRune && ru == 'h'):
			if sh.focusedCell[1] > 1 {
				sh.focusedCell[1]--
			}
		case k == tcell.KeyRight || (k == tcell.KeyRune && ru == 'l'):
			sh.focusedCell[1]++

		case k == tcell.KeyUp || (k == tcell.KeyRune && ru == 'k'):
			if sh.focusedCell[0] > 1 {
				sh.focusedCell[0]--
			}

		case k == tcell.KeyDown || (k == tcell.KeyRune && ru == 'j'):
			sh.focusedCell[0]++

		// Move to top row
		case (k == tcell.KeyHome && e.Modifiers() == tcell.ModShift) ||
			(k == tcell.KeyRune && ru == 'g'):
			sh.focusedCell[0] = 1

		// Move to left most column
		case k == tcell.KeyHome || (k == tcell.KeyRune && ru == '^'):
			sh.focusedCell[1] = 1
		}

		sh.adjustViewport(rows, cols)
		setFocus(sh)
	})
}

func (sh *Sheet) MouseHandler() func(ac tview.MouseAction, e *tcell.EventMouse, focus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return sh.Box.MouseHandler()
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
