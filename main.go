package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

func init() { /// Handle ctrl-c close
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		kill_process = true
		stopLoop()
		fmt.Println()
		os.Exit(1)
	}()
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
				if Gnomes.Init {
					Gnomes.Indexer.Endpoint = daemonAddress
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
