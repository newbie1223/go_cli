// Demo code for the Grid primitive.
// package main

// import (
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// func main() {
// current, err := os.Getwd()
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println(current)
// sc := bufio.NewScanner(os.Stdin)
// for sc.Scan() {
// 	input := sc.Text()
// 	inputs := strings.Split(input, " ")
// 	switch inputs[0] {
// 	case "exit":
// 		return
// 	case "ls":
// 		f, err := os.Open(current)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		files, err := f.Readdir(-1)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		for _, file := range files {
// 			fmt.Println(file.Name())
// 		}
// 	case "cd":
// 		// 実際にディレクトリが存在するかどうかはチェックする
// 		f, err := os.Open(current)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		files, err := f.Readdir(-1)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		exist := false
// 		for _, file := range files {
// 			if file.Name() == inputs[1] {
// 				if file.IsDir() {
// 					current = current + "/" + inputs[1]
// 				} else {
// 					fmt.Println("Not a directory")
// 				}
// 				exist = true
// 				break
// 			}
// 		}
// 		if !exist {
// 			fmt.Println("No such file or directory")
// 		}
// 	default:
// 		fmt.Println("Unknown command")
// 	}
// 	fmt.Print("current> ", current, "\n")
// }
// 	newPrimitive := func(text string) tview.Primitive {
// 		return tview.NewTextView().
// 			SetTextAlign(tview.AlignCenter).
// 			SetText(text)
// 	}
// 	menu := newPrimitive("Menu")
// 	main := newPrimitive("Main content")
// 	sideBar := newPrimitive("Side Bar")

// 	grid := tview.NewGrid().
// 		SetRows(-3, -1, -3).
// 		SetColumns(-33, -33, -34).
// 		SetBorders(true).
// 		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
// 		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

// 	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
// 	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
// 		AddItem(main, 1, 0, 1, 3, 0, 0, false).
// 		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

// 	// Layout for screens wider than 100 cells.
// 	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
// 		AddItem(main, 1, 1, 1, 1, 0, 100, false).
// 		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

// 	if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
// 		panic(err)
// 	}
// }

// package main

// import (
// 	"fmt"

// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// )

// func main() {
// 	app := tview.NewApplication()

// 	// チャット履歴を表示するためのTextViewを作成します。
// 	historyView := tview.NewTextView().
// 		SetDynamicColors(true).
// 		SetRegions(true).
// 		SetChangedFunc(func() {
// 			app.Draw()
// 		})

// 	// ユーザーがメッセージを入力するためのInputFieldを作成します。
// 	inputField := tview.NewInputField()
// 	inputField.SetLabel("メッセージ: ")
// 	inputField.SetFieldWidth(0)
// 	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		if event.Key() == tcell.KeyEnter {
// 			// メッセージを取得してクリアします。
// 			message := inputField.GetText()
// 			inputField.SetText("")

// 			// メッセージを履歴に追加します。
// 			fmt.Fprintf(historyView, "あなた: %s\n", message)
// 		}
// 		return event
// 	})

// 	// レイアウトを作成します。
// 	flex := tview.NewFlex().
// 		SetDirection(tview.FlexRow).
// 		AddItem(historyView, 0, 1, false).
// 		AddItem(inputField, 1, 1, true)

// 	if err := app.SetRoot(flex, true).Run(); err != nil {
// 		panic(err)
// 	}
// // }
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	inputField := tview.NewInputField()
	chatView := tview.NewTextView()

	go func() {
		listener, _ := net.Listen("tcp", ":8080")
		for {
			conn, _ := listener.Accept()
			go handleConnection(conn, chatView, app)
		}
	}()

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := inputField.GetText()
			if text != "" {
				fmt.Fprintf(chatView, "You: %s\n", text)
				conn, _ := net.Dial("tcp", "localhost:8080")
				fmt.Fprintf(conn, text+"\n")
				conn.Close()
				inputField.SetText("")
			}
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 1, 1, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func handleConnection(conn net.Conn, chatView *tview.TextView, app *tview.Application) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		if message == "" {
			break
		}
		message = strings.TrimSpace(message)
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(chatView, "%s: %s\n", conn.RemoteAddr(), message)
		})
		time.Sleep(100 * time.Millisecond)
	}
}
