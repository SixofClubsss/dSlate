package main

import (
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
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
	gnomes.InitLogrusLog(logrus.InfoLevel)
	myApp.Settings().SetTheme(bundle.DeroTheme(color.Black))
	myWindow = myApp.NewWindow("dSlate") /// start main app
	myWindow.SetMaster()                 /// if main closes, all windows close
	myWindow.Resize(fyne.NewSize(MIN_WIDTH, MIN_HEIGHT))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(resourceDReamTablesIconPng)
	myWindow.SetCloseIntercept(func() { /// do when close
		kill_process = true
		exit()
		myWindow.Close()
	})

	/// start looking and set content
	fetch()
	myWindow.SetContent(placeWith())
	myWindow.ShowAndRun()

}

// Handle ctrl-c close
func init() {
	gnomon.SetFastsync(true, true, 10000)
	gnomon.SetDBStorageType("boltdb")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		kill_process = true
		exit()
		fmt.Println()
		os.Exit(1)
	}()
}

// Main loop, ping daemon, get height check connection status and update Gnomon endpoint
func fetch() {
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
				gnomes.EndPoint()

			case <-quit: /// exit loop
				logger.Println("[dSlate] Exiting...")
				ticker.Stop()
				return
			}
		}
	}()
}

// To exit the routine running from fetch()
func exit() {
	quit <- struct{}{}
}
