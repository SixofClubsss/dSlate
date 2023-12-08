package main

import (
	"fmt"         // default formatting
	"image/color" // default color
	"os"          // default os control
	"os/signal"   // default os process control
	"syscall"     // this package needs to be replaced
	"time"        // default time package

	// fyne application
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// dreams application
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	// logger
	"github.com/sirupsen/logrus"
)

const (
	MIN_WIDTH  = 100 //
	MIN_HEIGHT = 400 //
)

var (

	//define the app
	appName = "secret-slate" // name the app

	myApp    = app.New()   // the application is a new fyne app
	myWindow fyne.Window   // contact the fyne interace
	quit     chan struct{} // define the cloese of the app

	debug        bool // toggle for debugging
	process_on   bool // toggle when process is running
	kill_process bool // toggle process to kill
)

func main() {
	// set up your logger
	menu.InitLogrusLog(logrus.InfoLevel)

	// set you app theme
	myApp.Settings().SetTheme(bundle.DeroTheme(color.Black))

	// start the app
	myWindow = myApp.NewWindow(appName)

	//set master process
	/// if main closes, all windows close
	myWindow.SetMaster()

	// set the window resizing function
	// use the const
	myWindow.Resize(fyne.NewSize(MIN_WIDTH, MIN_HEIGHT))

	// set the ability to change the app size
	myWindow.SetFixedSize(false)

	// set the app incon
	myWindow.SetIcon(resourceDerobotJpg)

	// set the closing procedures
	myWindow.SetCloseIntercept(func() { /// do when close
		kill_process = true // use your handy bool
		stopLoop()          // defined below
		myWindow.Close()    //close the window
	})

	/// start looking and set content
	fetchLoop()

	// set the app content
	myWindow.SetContent(placeWith()) // layout.go

	// now execute
	myWindow.ShowAndRun()

}

// Handle ctrl-c close
func init() {

	//
	menu.Gnomes.Fast = true
	//
	menu.Gnomes.DBType = "boltdb"

	// make signal 1
	c := make(chan os.Signal, 1)

	// tell the system that the app is closing
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// now send signal 1 into the channel
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
	var timer = time.NewTicker(6 * time.Second) // set a timer
	quit = make(chan struct{})                  // ?

	// now make an anon func
	go func() {
		for { // for loop
			select { // filter
			case <-timer.C: // do on timer
				rpc.Ping()            // // daemon_rpc.go
				isDaemonConnected()   // functions.go
				isWalletConnected()   // functions.go
				GetHeight()           // daemon_rpc.go
				menu.GnomonEndPoint() //gnomon.go

			case <-quit: /// exit loop
				logger.Println("[secret-slate] Exiting...")
				timer.Stop()
				return
			}
		}
	}()
}

func stopLoop() { /// to exit loop
	quit <- struct{}{}
}
