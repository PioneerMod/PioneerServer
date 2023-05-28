package main

import (
	"fmt"
	"os"
	"pioneer-server/io"
	"pioneer-server/net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const version = "0.0.1"

func main() {
	os.Setenv("PIONEER_SRV_VERS", version)

	if io.IsCLIMode(os.Args) {
		io.ExecuteCLI(os.Args)
		return
	}

	config := io.GetConfig()
	app := tview.NewApplication()

	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	var inputField *tview.InputField = nil
	inputField = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				fmt.Fprintf(logView, "%s\n", inputField.GetText())
				inputField.SetText("")
			}
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(logView, 0, 1, false).
		AddItem(inputField, 1, 0, true)

	io.InitLogger(config.Logger, logView)
	io.PrintSplashScreen(logView)
	io.Log(io.Info, "Initializing Pioneer Server...")

	server := net.CreateTcpServer(config.Server.IsLocal, config.Server.Port)
	server.Start()

	io.Log(io.Info, "Pioneer Server started!")

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
