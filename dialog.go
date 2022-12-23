package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func confirmPopUp() { /// pop up for entering custom daemon address
	cw := fyne.CurrentApp().NewWindow("Enter Custom Address")
	cw.Resize(fyne.NewSize(380, 50))
	cw.SetFixedSize(true)

	custom_input := widget.NewEntry()
	custom_input.SetPlaceHolder("Custom Daemon Address:")

	confirm := widget.NewButton("Enter", func() {
		log.Println("Confirm tapped")
		daemonAddress = custom_input.Text
		log.Println("Daemon Set To: CUSTOM")
	})

	content := container.NewBorder(
		nil,
		nil,
		nil,
		confirm,
		custom_input)

	cw.SetContent(content)
	cw.Show()
}

func searchPopUp(b, s, u, c string) { /// pop up display for sc search results
	sw := fyne.CurrentApp().NewWindow("Search Results")
	sw.Resize(fyne.NewSize(680, 800))
	sw.SetFixedSize(true)

	strings := widget.NewLabel(s)
	strings.Wrapping = fyne.TextWrapWord

	uints := widget.NewLabel(u)
	uints.Wrapping = fyne.TextWrapWord

	code := widget.NewLabel(c)
	code.Wrapping = fyne.TextWrapWord

	tabs := container.NewAppTabs(
		container.NewTabItem("String Keys", container.NewVScroll(strings)),
		container.NewTabItem("Uint Keys", container.NewVScroll(uints)),
		container.NewTabItem("Code", container.NewVScroll(code)),
	)
	content := container.NewVScroll(tabs)
	sw.SetContent(content)
	sw.Show()
}
