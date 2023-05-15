package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/SixofClubsss/dReams/bundle"
	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/rpc"
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
				if menu.Gnomes.Init {
					menu.Gnomes.Indexer.Endpoint = rpc.Daemon.Rpc
				}
			case <-quit: /// exit loop
				log.Println("[dSlate] Exiting...")
				ticker.Stop()
				return
			}
		}
	}()
}

func stopLoop() { /// to exit loop
	quit <- struct{}{}
}
