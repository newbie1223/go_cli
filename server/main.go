package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

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

	// inputField.SetChangedFunc(func(text string) { // inputFieldの入力内容が変更される度に実行
	// 	textView.SetText(textView.GetText(true) + inputField.GetText() + "\n")
	// 	inputField.SetText("")
	// 	return
	// })

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 0, true).
		AddItem(textView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
