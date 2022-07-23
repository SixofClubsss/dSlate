package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var myApp = app.New()
var myWindow fyne.Window
var quit chan struct{}
var daemonAddress string
var walletAddress string

const (
	MIN_WIDTH  = 380
	MIN_HEIGHT = 800
)

func main() {
	myWindow = myApp.NewWindow("dSlate") /// start main app
	myWindow.SetMaster()                 /// if main closes, all windows close
	myWindow.Resize(fyne.NewSize(MIN_WIDTH, MIN_HEIGHT))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(resourceDReamTablesIconPng)
	myWindow.SetCloseIntercept(func() { /// do when close
		stopLoop()
		myWindow.Close()
	})
	/// organize content
	settingsContent := container.NewWithoutLayout(rpcWalletEdit(), rpcLoginEdit(), rpcConnectButton(), daemonSelectOption(), daemonConnectBox(), walletConnectBox(), heightDisplay(), balanceDisplay(), builtOnImage())
	searchContent := container.NewWithoutLayout(searchButton(), contractEdit())
	scroll := container.NewScroll(searchContent)
	imageContent := container.NewWithoutLayout(cardImage())
	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", settingsContent),
		container.NewTabItem("Search", scroll),
		container.NewTabItem("Blank Slate", blankWidget()),
		container.NewTabItem("?", imageContent),
	)
	/// start looking and set content
	fetchLoop()
	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()

}

func fetchLoop() { /// ping daemon and get height loop
	var ticker = time.NewTicker(6 * time.Second)
	quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C: /// do on interval
				Ping()
				isDaemonConnected()
				isWalletConnected()
				GetHeight()
			case <-quit: /// exit loop
				fmt.Println("Exiting...")
				ticker.Stop()
				return
			}
		}
	}()
}

func stopLoop() { /// to exit loop
	quit <- struct{}{}
}
