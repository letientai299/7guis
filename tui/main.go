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

func main() {
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

func createApp() *tview.Application {
	tasks := tasks()
	pages := createPages(tasks)
	menu := createSidebar(tasks, pages)

	flex := tview.NewFlex()
	flex.AddItem(menu, len(url7GUIs)+3, 1, true)
	flex.AddItem(pages, 0, 1, false)

	app := tview.NewApplication().SetRoot(flex, true)
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
	for i, task := range tasks {
		var view tview.Primitive
		if task.widget != nil {
			view = task.widget
		} else {
			view = tview.NewTextView().SetText(task.name)
		}
		pages.AddPage(task.name, view, true, i == 0)
	}

	pages.SetBorder(true)
	return pages
}

func createSidebar(tasks []Task, pages *tview.Pages) tview.Primitive {
	menu := createMenu(tasks, pages)
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
	for i, task := range tasks {
		if len(task.name) > menuWidth {
			menuWidth = len(task.name)
		}
		menu.AddItem(task.name, "", rune(i+'0'), nil)
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

func tasks() []Task {
	tasks := []Task{
		{name: "Counter", widget: task.Counter()},
		{name: "Temperature Converter", widget: nil},
		{name: "Flight Booker", widget: nil},
		{name: "Timer", widget: nil},
		{name: "CRUD", widget: nil},
		{name: "Circle Drawer", widget: nil},
		{name: "Cells", widget: nil},
	}
	return tasks
}
