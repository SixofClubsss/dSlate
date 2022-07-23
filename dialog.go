package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var customDaemonInput = widget.NewEntry()

func confirmPopUp() { /// pop up for entering custom daemon address
	cw := fyne.CurrentApp().NewWindow("Enter Custom Address")
	cw.Resize(fyne.NewSize(400, 70))
	cw.SetFixedSize(true)
	content := container.NewWithoutLayout(confirmButton(), customeDaemonEdit())
	cw.SetContent(content)
	cw.Show()

}

func customeDaemonEdit() fyne.Widget { /// entry for custom daemon address
	customDaemonInput.SetPlaceHolder("Daemon Address")
	customDaemonInput.Resize(fyne.NewSize(270, 45))
	customDaemonInput.Move(fyne.NewPos(10, 10))
	return customDaemonInput
}

func confirmButton() fyne.Widget { /// confirming custom daemon address
	confirmButton := widget.NewButton("Enter", func() {
		log.Println("confirm tapped")
		suff := "/json_rpc"
		daemonAddress = customDaemonInput.Text + suff
		log.Println("Daemon Set To: CUSTOM")
	})
	confirmButton.Resize(fyne.NewSize(100, 42))
	confirmButton.Move(fyne.NewPos(290, 11))
	return confirmButton
}

func searchPopUp(b, s, u, c string) { /// pop up display for sc search results
	sw := fyne.CurrentApp().NewWindow("Search Results")
	sw.Resize(fyne.NewSize(680, 800))
	sw.SetFixedSize(true)
	content := container.NewVScroll(searchDisplay())
	sw.SetContent(content)
	display.SetText("SC ID: \n" + contractInput.Text + "\n\n" + "Balances: \n" + b + "\n\nString Keys: \n" + s + "\n\nUint Keys: \n" + u + "\n\nCode: " + c)
	sw.Show()

}
