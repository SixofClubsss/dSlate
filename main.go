package main

import (
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

const (
	MIN_WIDTH  = 380
	MIN_HEIGHT = 800
)

var (
	myApp    = app.New()
	myWindow fyne.Window
	quit     chan struct{}

	debug        bool
	process_on   bool
	kill_process bool
)

func main() {
	menu.InitLogrusLog(runtime.GOOS == "windows")
	myApp.Settings().SetTheme(bundle.DeroTheme(color.Black))
	myWindow = myApp.NewWindow("dSlate") /// start main app
	myWindow.SetMaster()                 /// if main closes, all windows close
	myWindow.Resize(fyne.NewSize(MIN_WIDTH, MIN_HEIGHT))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(resourceDReamTablesIconPng)
	myWindow.SetCloseIntercept(func() { /// do when close
		kill_process = true
		stopLoop()
		myWindow.Close()
	})

	/// start looking and set content
	fetchLoop()
	myWindow.SetContent(placeWith())
	myWindow.ShowAndRun()

}

// Handle ctrl-c close
func init() {
	menu.Gnomes.Fast = true
	menu.Gnomes.DBType = "boltdb"
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		kill_process = true
		stopLoop()
		fmt.Println()
		os.Exit(1)
	}()
}

// Main loop, ping daemon, get height check connection status and update Gnomon endpoint
func fetchLoop() {
	var ticker = time.NewTicker(6 * time.Second)
	quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C: /// do on interval
				rpc.Ping()
				isDaemonConnected()
				isWalletConnected()
				GetHeight()
				menu.GnomonEndPoint()

			case <-quit: /// exit loop
				logger.Println("[dSlate] Exiting...")
				ticker.Stop()
				return
			}
		}
	}()
}

func stopLoop() { /// to exit loop
	quit <- struct{}{}
}
