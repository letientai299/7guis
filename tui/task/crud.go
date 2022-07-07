package task

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func CRUD() tview.Primitive {
	data := []*person{
		{name: "Barry", surname: "Allen"},
		{name: "Bruce", surname: "Wayne"},
		{name: "Clark", surname: "Kent"},
		{name: "Diana", surname: "Prince"},
	}

	for i := 0; i < 40; i++ {
		data = append(data, &person{
			name:    fmt.Sprintf("%d name", i),
			surname: fmt.Sprintf("%d surname", i),
		})
	}

	c := &crud{
		persons:  data,
		filtered: append([]*person{}, data...),
	}

	return c.render()
}

type person struct {
	name    string
	surname string
}

func parsePerson(text string) *person {
	ss := strings.Split(text, ",")
	p := &person{}
	if len(ss) == 0 {
		return p
	}

	p.surname = strings.TrimSpace(ss[0])
	if len(ss) < 2 {
		return p
	}

	p.name = strings.TrimSpace(ss[1])
	return p
}

type persons []*person

func (p persons) GetCell(row, _ int) *tview.TableCell {
	cell := tview.NewTableCell(p[row].String())
	return cell
}

func (p persons) GetRowCount() int                          { return len(p) }
func (p persons) GetColumnCount() int                       { return 1 }
func (p persons) SetCell(row, _ int, cell *tview.TableCell) { p[row] = parsePerson(cell.Text) }
func (p persons) RemoveRow(_ int)                           {}
func (p persons) RemoveColumn(_ int)                        {}
func (p persons) InsertRow(_ int)                           {}
func (p persons) InsertColumn(_ int)                        {}
func (p persons) Clear()                                    {}

func (p person) String() string {
	if p.surname == "" {
		return p.name
	}

	if p.name == "" {
		return p.surname
	}

	return p.surname + ", " + p.name
}

type crud struct {
	persons      persons
	filtered     persons
	selected     int
	table        *tview.Table
	filterText   string
	inputName    string
	inputSurname string
}

func (c *crud) render() tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true)

	filterInput := tview.NewInputField().
		SetLabel("Filter prefix").
		SetLabelWidth(16).
		SetChangedFunc(c.setFilter)
	container.AddItem(filterInput, 1, 0, true)

	c.table = tview.NewTable().
		SetContent(c.filtered).
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(row, _ int) {
			c.selected = row
		})

	c.table.SetBorder(true)

	container.AddItem(
		tview.NewFlex().
			AddItem(
				c.table,
				0, 1, false,
			).
			AddItem(
				tview.NewForm().
					AddInputField("Name", "", 20, nil, c.inputNameChange).
					AddInputField("Surname", "", 20, nil, c.inputSurnameChange),
				0, 1, true,
			),
		0, 1, true,
	)

	container.AddItem(
		tview.NewFlex().
			AddItem(tview.NewButton("Create").SetSelectedFunc(c.onCreate), 0, 1, true).
			AddItem(tview.NewBox(), 1, 0, false).
			AddItem(tview.NewButton("Update").SetSelectedFunc(c.onUpdate), 0, 1, true).
			AddItem(tview.NewBox(), 1, 0, false).
			AddItem(tview.NewButton("Delete").SetSelectedFunc(c.onDelete), 0, 1, true),
		1, 0, false,
	)

	return centerScreen(container, 50, 16)
}

func (c *crud) setFilter(text string) {
	c.filterText = text
	c.onFilter()
}

func (c *crud) onFilter() {
	text := strings.ToLower(c.filterText)
	c.filtered = c.filtered[:0]
	for i := 0; i < len(c.persons); i++ {
		p := c.persons[i]
		if strings.HasPrefix(strings.ToLower(p.name), text) ||
			strings.HasPrefix(strings.ToLower(p.surname), text) {
			c.filtered = append(c.filtered, p)
		}
	}

	c.table.SetContent(c.filtered)
}

func (c *crud) inputNameChange(text string)    { c.inputName = text }
func (c *crud) inputSurnameChange(text string) { c.inputSurname = text }

func (c *crud) onCreate() {
	if c.inputName == "" && c.inputSurname == "" {
		return
	}

	c.persons = append(c.persons, &person{
		name:    c.inputName,
		surname: c.inputSurname,
	})

	c.onFilter()
}

func (c *crud) onUpdate() {
	if c.inputName == "" && c.inputSurname == "" {
		return
	}

	p := c.filtered[c.selected]
	p.name = c.inputName
	p.surname = c.inputSurname
}

func (c *crud) onDelete() {
	toBeDeleted := c.filtered[c.selected]
	n := 0
	for _, p := range c.persons {
		if p == toBeDeleted {
			continue
		}
		c.persons[n] = p
		n++
	}
	c.persons = c.persons[:n]
	c.onFilter()
}
