package task

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type flightBooker struct {
	app            *tview.Application
	typeInput      *tview.DropDown
	departureInput *tview.InputField
	arrivalInput   *tview.InputField
	bookButton     *tview.Button
	focusedIndex   int
}

func FlightBooker(app *tview.Application) tview.Primitive {
	fb := &flightBooker{app: app}
	return fb.render()
}

func (fb *flightBooker) render() tview.Primitive {
	const labelWidth = 20
	const optionWidth = 24
	const datePlaceHolder = "dd.mm.yyyy"

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	fb.typeInput = tview.NewDropDown().
		SetLabel("Flight type").
		SetLabelWidth(labelWidth).
		SetFieldWidth(optionWidth).
		SetOptions([]string{
			"One way",
			"Return",
		}, nil).
		SetCurrentOption(0)
	flex.AddItem(fb.typeInput, 2, 1, true)

	fb.departureInput = tview.NewInputField().
		SetLabel("Departure").
		SetPlaceholder(datePlaceHolder).
		SetLabelWidth(labelWidth)
	flex.AddItem(fb.departureInput, 2, 1, true)

	fb.arrivalInput = tview.NewInputField().
		SetLabel("Arrival").
		SetPlaceholder(datePlaceHolder).
		SetLabelWidth(labelWidth)
	flex.AddItem(fb.arrivalInput, 2, 1, false)

	fb.bookButton = tview.NewButton("Book")
	flex.AddItem(fb.bookButton, 1, 1, true)
	flex.SetInputCapture(fb.inputSwitch)
	fb.focusedIndex = 0 // flight type field

	wrap := tview.NewFlex().SetDirection(tview.FlexRow)
	wrap.SetBorder(true)
	wrap.AddItem(flex, 0, 1, true)
	hint := tview.NewTextView().SetText("Use Ctrl-N/Ctrl-P to switch input fields")
	wrap.AddItem(hint, 1, 0, false)

	fb.typeInput.SetSelectedFunc(fb.onFlightTypeChange)
	fb.typeInput.SetCurrentOption(0)

	fb.departureInput.SetChangedFunc(func(s string) {
		_, err := fb.parseDate(s)
		if err != nil {
			fb.departureInput.SetFieldTextColor(ColorInvalid)
		} else {
			fb.departureInput.SetFieldTextColor(tview.Styles.PrimaryTextColor)
		}
		fb.adjustBookButton()
	})

	fb.arrivalInput.SetChangedFunc(func(s string) {
		_, err := fb.parseDate(s)
		if err != nil {
			fb.arrivalInput.SetFieldTextColor(ColorInvalid)
		} else {
			fb.arrivalInput.SetFieldTextColor(tview.Styles.PrimaryTextColor)
		}
		fb.adjustBookButton()
	})

	fb.bookButton.SetSelectedFunc(fb.book)

	fb.adjustBookButton()
	return centerScreen(wrap, labelWidth+optionWidth, 11)
}

func (fb *flightBooker) inputSwitch(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		fb.focusNext()
	case tcell.KeyCtrlP:
		fb.focusBack()
	default:
		return event
	}

	return nil
}

func (fb *flightBooker) focusNext() {
	fs := fb.enabledFields()
	if fb.focusedIndex < len(fs)-1 {
		fb.focusedIndex++
	} else {
		fb.focusedIndex = 0
	}
	fb.app.SetFocus(fs[fb.focusedIndex])
}

func (fb *flightBooker) focusBack() {
	fs := fb.enabledFields()
	if fb.focusedIndex > 0 {
		fb.focusedIndex--
	} else {
		fb.focusedIndex = len(fs) - 1
	}
	fb.app.SetFocus(fs[fb.focusedIndex])
}

func (fb *flightBooker) enabledFields() []tview.Primitive {
	enables := []tview.Primitive{
		fb.typeInput, fb.departureInput,
	}
	if index, _ := fb.typeInput.GetCurrentOption(); index == 1 {
		enables = append(enables, fb.arrivalInput)
	}

	if fb.canBook() {
		enables = append(enables, fb.bookButton)
	}
	return enables
}

func (fb *flightBooker) adjustBookButton() bool {
	v := fb.canBook()
	if v {
		enable(fb.bookButton)
	} else {
		disable(fb.bookButton)
	}
	return v
}

func (fb *flightBooker) canBook() bool {
	depart := fb.departureInput.GetText()
	if i, _ := fb.typeInput.GetCurrentOption(); i == 0 { // one way
		_, err := fb.parseDate(depart)
		return err == nil
	}

	arrival := fb.arrivalInput.GetText()
	d, errD := fb.parseDate(depart)
	a, errA := fb.parseDate(arrival)
	if errD != nil || errA != nil {
		return false
	}

	return d.Before(a)
}

func (fb *flightBooker) onFlightTypeChange(_ string, index int) {
	if index == 0 { // one-way
		disable(fb.arrivalInput)
		return
	}

	enable(fb.arrivalInput)
	fb.adjustBookButton()
}

func (fb *flightBooker) parseDate(v string) (time.Time, error) {
	return time.Parse("02.01.2006", v)
}

func (fb *flightBooker) book() {
	fb.app.Suspend(func() {
		var msg string
		if i, _ := fb.typeInput.GetCurrentOption(); i == 0 {
			msg = fmt.Sprintf("Booked a one way flight for %s", fb.departureInput.GetText())
		} else {
			msg = fmt.Sprintf("Booked a returned flight for %s and %s",
				fb.departureInput.GetText(),
				fb.arrivalInput.GetText(),
			)
		}

		log.Println(msg)
		log.Println("Press Enter to continue")
		_, _ = fmt.Scanln()
	})
}
