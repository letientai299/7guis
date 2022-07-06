package task

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

func TemperatureConverter() tview.Primitive {
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	errDisplay := tview.NewTextView().SetText("OK")
	errDisplay.SetBorderPadding(0, 0, 1, 1)
	errDisplay.SetDynamicColors(true)

	form := tview.NewForm()
	container.AddItem(form, 5, 1, true)
	container.AddItem(errDisplay, 2, 1, false)

	celsiusInput := tview.NewInputField().
		SetLabel("Celsius").
		SetText("0")

	form.AddFormItem(celsiusInput)

	fahrenheitInput := tview.NewInputField().
		SetLabel("Fahrenheit").
		SetText("32")

	form.AddFormItem(fahrenheitInput)
	form.SetItemPadding(1)
	form.SetBorderPadding(0, 0, 0, 0)
	form.SetBorder(true)

	setErr := func(s string) {
		if len(s) == 0 {
			s = "OK"
		}
		errDisplay.SetText(s)
	}

	var c2f, f2c func(string)

	c2f = func(s string) {
		if len(s) == 0 {
			setErr("")
			return
		}

		c, err := strconv.ParseFloat(s, 64)
		if err != nil {
			setErr("[red]Invalid C value: " + s)
			return
		}

		f := c*9.0/5.0 + 32

		fahrenheitInput.SetChangedFunc(nil)
		fahrenheitInput.SetText(fmt.Sprintf("%.2f", f))
		fahrenheitInput.SetChangedFunc(f2c)
	}
	celsiusInput.SetChangedFunc(c2f)

	f2c = func(s string) {
		if len(s) == 0 {
			setErr("")
			return
		}

		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			setErr("[red]Invalid F value: " + s)
			return
		}
		c := (f - 32) * 5.0 / 9.0
		celsiusInput.SetChangedFunc(nil)
		celsiusInput.SetText(fmt.Sprintf("%.2f", c))
		celsiusInput.SetChangedFunc(c2f)
	}
	fahrenheitInput.SetChangedFunc(f2c)

	container.SetBorder(true)
	return centerScreen(container, 40, 9)
}
