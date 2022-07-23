package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/deroproject/derohe/rpc"
)

var primes = []string{"MAINNET", "TESTNET", "SIMULATOR", "CUSTOM"} /// set select menu
var dropDown = widget.NewSelect(primes, func(s string) {           /// do when select changes
	whichDaemon(s)
	log.Println("Daemon Set To:", s)
})

var rpcLoginInput = widget.NewPasswordEntry() /// declare some widgets
var rpcWalletInput = widget.NewEntry()
var contractInput = widget.NewEntry()
var display = widget.NewLabel("")

var daemonCheckBox = widget.NewCheck("Daemon Connected", func(value bool) {
	/// do something on change
})

var walletCheckBox = widget.NewCheck("Wallet Connected", func(value bool) {
	/// do something on change
})

var currentHeight = widget.NewEntry()
var walletBalance = widget.NewEntry()

func rpcLoginEdit() fyne.Widget { /// user:pass password entry
	rpcLoginInput.SetPlaceHolder("Enter user:pass")
	rpcLoginInput.Resize(fyne.NewSize(360, 45))
	rpcLoginInput.Move(fyne.NewPos(10, 650))

	return rpcLoginInput
}

func rpcWalletEdit() fyne.Widget { /// wallet rpc address entry
	rpcWalletInput.SetPlaceHolder("Wallet RPC Address")
	rpcWalletInput.Resize(fyne.NewSize(250, 45))
	rpcWalletInput.Move(fyne.NewPos(10, 700))

	return rpcWalletInput
}

func rpcConnectButton() fyne.Widget { /// wallet connect button
	button := widget.NewButton("Connect", func() { /// do on pressed
		log.Println("Connect Pressed")
		pre := "http://"
		suff := "/json_rpc"
		walletAddress = pre + rpcWalletInput.Text + suff
		GetAddress()
	})
	button.Resize(fyne.NewSize(100, 42))
	button.Move(fyne.NewPos(270, 702))

	return button
}

func daemonSelectOption() fyne.Widget { /// daemon select menu
	dropDown.SetSelectedIndex(0)
	dropDown.Resize(fyne.NewSize(180, 45))
	dropDown.Move(fyne.NewPos(10, 550))

	return dropDown
}

func daemonConnectBox() fyne.Widget { /// daemon check box
	daemonCheckBox.Resize(fyne.NewSize(30, 30))
	daemonCheckBox.Move(fyne.NewPos(3, 595))
	daemonCheckBox.Disable()

	return daemonCheckBox
}

func walletConnectBox() fyne.Widget { /// wallet check box
	walletCheckBox.Resize(fyne.NewSize(30, 30))
	walletCheckBox.Move(fyne.NewPos(3, 620))
	walletCheckBox.Disable()

	return walletCheckBox
}

func heightDisplay() fyne.Widget { /// height display entry is read only
	currentHeight.SetText("Height:")
	currentHeight.Disable()
	currentHeight.Resize(fyne.NewSize(170, 45))
	currentHeight.Move(fyne.NewPos(200, 550))

	return currentHeight

}

func balanceDisplay() fyne.Widget {
	walletBalance.SetText("Balance:")
	walletBalance.Disable()
	walletBalance.Resize(fyne.NewSize(170, 45))
	walletBalance.Move(fyne.NewPos(200, 600))

	return walletBalance

}

func searchDisplay() fyne.Widget { /// label for search results
	display.Resize(fyne.NewSize(360, 780))
	display.Move(fyne.NewPos(5, 10))
	display.Wrapping = fyne.TextWrapWord

	return display
}

func contractEdit() fyne.Widget { /// contract entry
	contractInput.SetPlaceHolder("Enter Contract Id:")
	contractInput.Resize(fyne.NewSize(360, 45))
	contractInput.Move(fyne.NewPos(10, 15))

	return contractInput
}

func searchButton() fyne.Widget { /// SC search button
	button := widget.NewButton("Search", func() {
		log.Println("Searching for: " + contractInput.Text)
		p := &rpc.GetSC_Params{
			SCID:      contractInput.Text,
			Code:      true,
			Variables: true,
		}
		getSC(p)
	})
	button.Resize(fyne.NewSize(360, 42))
	button.Move(fyne.NewPos(10, 63))
	return button
}

func builtOnImage() fyne.CanvasObject { ///  main image
	img := canvas.NewImageFromResource(resourceBuiltOnDeroPng)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(360, 430))
	img.Move(fyne.NewPos(10, 210))

	return img
}

func cardImage() fyne.CanvasObject { /// card image
	img := canvas.NewImageFromResource(resourceDero1Png)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(450, 330))
	img.Move(fyne.NewPos(-33, 200))

	return img
}

func blankWidget() fyne.Widget { /// slate label
	blank := widget.NewLabel("Something goes here...")
	return blank
}
