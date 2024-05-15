package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	// "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type ui struct {
	app     *tview.Application
	curTime *tview.TableCell
}

const refreshInterval = 500 * time.Millisecond

func currentTimeString() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("15:04:05"))
}

func updateTime(u *ui) {
	for {
		time.Sleep(refreshInterval)
		u.app.QueueUpdateDraw(func() {
			u.curTime.SetText(currentTimeString())
		})
	}
}

func createCommandList() (commandList *tview.List) {
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Command")
	return commandList
}

func createInfoPanel(app *tview.Application, u *ui) (infoPanel *tview.Grid) {

	infoTable := tview.NewTable()
	infoTable.SetBorder(true).SetTitle("Information")

	cnt := 0
	infoTable.SetCellSimple(cnt, 0, "Time:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	u.curTime = tview.NewTableCell("0")
	infoTable.SetCell(cnt, 1, u.curTime)
	cnt++

	// infoPanel = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(infoTable, 0, 1, false)
	infoPanel = tview.NewGrid().SetRows(-1).SetColumns(-1).AddItem(infoTable, 0, 0, 1, 1, 0, 0, false)

	return infoPanel
}

func createLayout(cList tview.Primitive, recvPanel tview.Primitive, tPanel tview.Primitive) (layout *tview.Flex) {
	// bodyLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
	// 	AddItem(cList, 20, 1, true).
	// 	AddItem(recvPanel, 0, 1, false).
	// 	AddItem(tPanel, 0, 1, false)

	bodyLayout := tview.NewGrid().SetRows(3, 1, 3).SetColumns(3, 3, 4).
		SetBorders(true).
		AddItem(cList, 0, 0, 1, 3, 0, 0, false).
		AddItem(recvPanel, 2, 0, 1, 3, 0, 0, false).
		AddItem(tPanel, 1, 1, 3, 2, 0, 0, false)

	layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(bodyLayout, 0, 1, true)
	// layout = tview.NewGrid().SetRows(-1).SetColumns(-1).AddItem(bodyLayout, 0, 0, 1, 1, 0, 0, true)

	// layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(bodyLayout, 0, 1, true)
	// bodyLayout := tview.NewGrid().AddItem(cList, 0, 0, 0, 0, 0, 0, true).AddItem(recvPanel, 0, 0, 1, 2, 100, 0, true).AddItem(tPanel, 1, 1, 3, 2, 300, 0, true)
	// layout = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(bodyLayout, 0, 1, true)
	return layout
}

// func createModalForm(pages *tview.Pages, form tview.Primitive, height int, width int) tview.Primitive {
// 	modal := tview.NewFlex().SetDirection(tview.FlexColumn).
// 		AddItem(nil, 0, 1, false).
// 		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
// 			AddItem(nil, 0, 1, false).
// 			AddItem(form, height, 1, true).
// 			AddItem(nil, 0, 1, false), width, 1, true).
// 		AddItem(nil, 0, 1, false)
// 	return modal
// }

func createTextViewPanel(app *tview.Application, title string) (textPanel *tview.TextView) {
	textPanel = tview.NewTextView()
	textPanel.SetTitle(title)
	textPanel.SetBorder(true)
	return textPanel
}

func createTextInputPanel(app *tview.Application, title string) (inputPanel *tview.Grid) {

	textView := tview.NewTextView()
	textView.SetTitle("gwork")
	textView.SetBorder(true)

	inputField := tview.NewInputField()
	inputField.SetLabel("@user: ")
	inputField.SetTitle("chat").SetBorder(true)
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			textView.SetText(textView.GetText(true) + inputField.GetText() + "\n")
			inputField.SetText("")
			return nil
		}
		return event
	})

	// inputPanel = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(inputField, 3, 0, true).AddItem(textView, 0, 1, false)
	inputPanel = tview.NewGrid().SetRows(-1, -3).SetColumns(-1).AddItem(textView, 0, 0, 1, 1, 0, 0, false).AddItem(inputField, 1, 0, 1, 1, 0, 0, false)

	return inputPanel

}

// func testCommand(pages *tview.Pages) func() {
// 	return func() {
// 		cancelFunc := func() {
// 			pages.SwitchToPage("main")
// 			pages.RemovePage("modal")
// 		}

// 		onFunc := func() {
// 			pages.SwitchToPage("main")
// 			pages.RemovePage("modal")
// 		}

// 		form := tview.NewForm()
// 		form.AddButton("ON", onFunc)
// 		form.AddButton("Cancel", cancelFunc)
// 		form.SetCancelFunc(cancelFunc)
// 		form.SetButtonsAlign(tview.AlignCenter)
// 		form.SetBorder(true).SetTitle("Test")
// 		modal := createModalForm(pages, form, 13, 55)
// 		pages.AddPage("modal", modal, true, true)
// 	}
// }

func createApplication() (app *tview.Application) {
	app = tview.NewApplication()
	pages := tview.NewPages()

	ui := &ui{}
	ui.app = app
	infoPanel := createInfoPanel(app, ui)
	textPanel := createTextViewPanel(app, "gwork")
	// inputPanel := createTextViewPanel(app, "chat",textPanel)

	commandList := createCommandList()
	// commandList.AddItem("Test", "", 'p', testCommand(pages))
	commandList.AddItem("Test", "", 'p', func() {
		textPanel.SetText(textPanel.GetText(true) + "test\n")
	})
	commandList.AddItem("Clear", "", 'c', func() {
		textPanel.SetText("")
	})
	commandList.AddItem("Quit", "", 'q', func() {
		app.Stop()
	})
	commandList.AddItem("Input", "", 'i', func() {
		inputPanel := createTextInputPanel(app, "chat")
		// pages.AddPage("input", inputPanel, true, true)
		pages.AddPage("main", inputPanel, true, true)
	})
	layout := createLayout(commandList, infoPanel, textPanel)
	pages.AddPage("main", layout, true, true)

	go updateTime(ui)
	app.SetRoot(pages, true)
	return app
}

func main() {
	// runewidth.DefaultCondition = &runewidth.Condition{EastAsianWidth: false}

	app := createApplication()

	if err := app.Run(); err != nil {
		panic(err)
	}
	if err := app.EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
