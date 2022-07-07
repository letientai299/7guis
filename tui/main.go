package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/letientai299/7guis/tui/task"
)

const url7GUIs = "https://eugenkiss.github.io/7guis/tasks/"

func tasks(app *tview.Application) []Task {
	tasks := []Task{
		{name: "Counter", widget: task.Counter()},
		{name: "Temperature Converter", widget: task.TemperatureConverter()},
		{name: "Flight Booker", widget: task.FlightBooker(app)},
		{name: "Timer", widget: nil},
		{name: "CRUD", widget: nil},
		{name: "Circle Drawer", widget: nil},
		{name: "Cells", widget: nil},
	}
	return tasks
}

func main() {
	configStyles()
	app := createApp()
	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(
			signals,
			syscall.SIGILL, syscall.SIGINT, syscall.SIGTERM,
		)
		<-signals
		app.Stop()
	}()

	<-done
}

func configStyles() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,
		ContrastBackgroundColor:     tcell.ColorDarkBlue,
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorGhostWhite,
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorDeepSkyBlue,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
	}
}

func createApp() *tview.Application {
	app := tview.NewApplication()
	tasks := tasks(app)

	pages := createPages(tasks)
	menu := createSidebar(tasks, pages)

	flex := tview.NewFlex()
	flex.AddItem(menu, len(url7GUIs)+3, 1, false)
	flex.AddItem(pages, 0, 1, true)

	app.SetRoot(flex, true)
	app.EnableMouse(true)

	focusingMenu := true
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key != tcell.KeyTAB && key != tcell.KeyBacktab {
			return event
		}

		if focusingMenu {
			app.SetFocus(pages)
		} else {
			app.SetFocus(menu)
		}

		focusingMenu = !focusingMenu
		return nil
	})

	return app
}

func createPages(tasks []Task) *tview.Pages {
	pages := tview.NewPages()
	for i, t := range tasks {
		var view tview.Primitive
		if t.widget != nil {
			view = t.widget
		} else {
			view = tview.NewTextView().SetText(t.name)
		}
		pages.AddPage(t.name, view, true, i == 0)
	}

	pages.SetBorder(true)
	return pages
}

func createSidebar(tasks []Task, pages *tview.Pages) tview.Primitive {
	menu := createMenu(tasks, pages)
	menu.SetCurrentItem(3)
	frame := tview.NewFrame(menu)
	frame.SetBorder(true)
	frame.SetBorders(0, 0, 1, 1, 1, 1)
	frame.AddText("7GUIs", true, tview.AlignLeft, tcell.ColorWhite)
	frame.AddText(url7GUIs, true, tview.AlignLeft, tcell.ColorBlue)
	return frame
}

func createMenu(tasks []Task, pages *tview.Pages) *tview.List {
	menu := tview.NewList()
	menuWidth := 0
	for i, t := range tasks {
		if len(t.name) > menuWidth {
			menuWidth = len(t.name)
		}
		menu.AddItem(t.name, "", rune(i+'0'), nil)
	}

	menu.ShowSecondaryText(false)
	menu.SetChangedFunc(func(_ int, task string, _ string, _ rune) {
		pages.SwitchToPage(task)
	})

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyRune {
			return event
		}

		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, event.Rune(), event.Modifiers())
		}

		if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, event.Rune(), event.Modifiers())
		}

		return event
	})

	return menu
}

type Task struct {
	name   string
	widget tview.Primitive
}
