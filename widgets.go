package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
)

// / declare some widgets
var (
	primes   = []string{"MAINNET", "TESTNET", "SIMULATOR", "CUSTOM"} /// set select menu
	dropDown = widget.NewSelect(primes, func(s string) {             /// do when select changes
		whichDaemon(s)
		logger.Println("[secret-slate] Daemon Set To:", s)
	})

	rpcLoginInput  = widget.NewPasswordEntry()
	rpcWalletInput = widget.NewSelectEntry([]string{WALLET_SIMULATOR_DEFAULT, WALLET_TESTNET_DEFAULT, WALLET_MAINNET_DEFAULT})
	contractInput  = widget.NewMultiLineEntry()

	daemonCheckBox = widget.NewCheck("Daemon Connected", func(value bool) {
		menu.Gnomes.Stop("dSlate")
	})

	walletCheckBox = widget.NewCheck("Wallet Connected", func(value bool) {
		/// do something on change
	})

	currentHeight = widget.NewEntry()
	walletBalance = widget.NewEntry()

	gnomonEnabled = widget.NewRadioGroup([]string{}, func(s string) {})

	debugEnabled = widget.NewCheck("Debug", func(b bool) {
		if b {
			debug = true
		} else {
			debug = false
		}
	})
)

// User:pass password entry
func rpcLoginEdit() fyne.Widget {
	rpcLoginInput.SetPlaceHolder("RPC: user:pass")
	rpcLoginInput.Resize(fyne.NewSize(360, 45))
	rpcLoginInput.Move(fyne.NewPos(10, 650))
	this_pass := binding.BindString(&rpc.Wallet.UserPass)
	rpcLoginInput.Bind(this_pass)

	return rpcLoginInput
}

// Wallet rpc address entry
func rpcWalletEdit() fyne.Widget {
	rpcWalletInput.SetPlaceHolder("Wallet IP:PORT")
	rpcWalletInput.Resize(fyne.NewSize(250, 45))
	rpcWalletInput.Move(fyne.NewPos(10, 700))

	return rpcWalletInput
}

// Wallet connect button
func rpcConnectButton() fyne.Widget {
	button := widget.NewButton("Connect", func() { /// do on pressed
		rpc.Wallet.Rpc = rpcWalletInput.Text
		GetAddress()
	})
	button.Resize(fyne.NewSize(100, 42))
	button.Move(fyne.NewPos(270, 702))

	return button
}

// Daemon select menu
func daemonSelectOption() fyne.Widget {
	dropDown.SetSelectedIndex(0)
	dropDown.Resize(fyne.NewSize(180, 45))
	dropDown.Move(fyne.NewPos(10, 550))

	return dropDown
}

// Daemon check box
func daemonConnectBox() fyne.Widget {
	daemonCheckBox.Resize(fyne.NewSize(30, 30))
	daemonCheckBox.Move(fyne.NewPos(3, 595))
	daemonCheckBox.Disable()

	return daemonCheckBox
}

// Wallet check box
func walletConnectBox() fyne.Widget {
	walletCheckBox.Resize(fyne.NewSize(30, 30))
	walletCheckBox.Move(fyne.NewPos(3, 620))
	walletCheckBox.Disable()

	return walletCheckBox
}

// Height display entry is read only
func heightDisplay() fyne.Widget {
	currentHeight.SetText("Height:")
	currentHeight.Disable()
	currentHeight.Resize(fyne.NewSize(170, 45))
	currentHeight.Move(fyne.NewPos(200, 550))

	return currentHeight

}

//

// Balance display entry is read only
func balanceDisplay() fyne.Widget {
	walletBalance.SetText("Balance:")
	walletBalance.Disable()
	walletBalance.Resize(fyne.NewSize(170, 45))
	walletBalance.Move(fyne.NewPos(200, 600))

	return walletBalance

}

// SC entry
func contractEdit() fyne.Widget {
	contractInput.SetPlaceHolder("Enter Contract Id:")
	contractInput.Resize(fyne.NewSize(360, 45))
	contractInput.Move(fyne.NewPos(10, 15))
	contractInput.Wrapping = fyne.TextWrapWord

	return contractInput
}

// Print SC code button
func contractCode() fyne.Widget {
	button := widget.NewButton("SC Code", func() {
		if len(contractInput.Text) == 64 {
			getSCcode(contractInput.Text)
		}
	})

	return button
}

// SC search button
func searchButton() fyne.Widget {
	button := widget.NewButton("Search", func() {
		logger.Println("[secret-slate] Searching for: " + contractInput.Text)
		p := &dero.GetSC_Params{
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

// Main image
func builtOnImage() fyne.CanvasObject {
	img := canvas.NewImageFromResource(resource8BiVILihPng)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(38, 54))
	img.Move(fyne.NewPos(10, 210))

	return img
}

// Card image
func cardImage() fyne.CanvasObject {
	img := canvas.NewImageFromResource(resourceQrcodePng)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(450, 330))
	img.Move(fyne.NewPos(-3, 200))

	return img
}

// Blank slate label
func blankWidget() fyne.Widget {
	blank := widget.NewLabel("Something goes here...")
	return blank
}

// Blank slate label
func walletAddress() fyne.Widget {
	stringTxt := fmt.Sprintf(rpc.Wallet.Address)
	blank := widget.NewLabel(stringTxt)
	return blank
}

// Enable/Disable Gnomon radio group
func enableGnomon() fyne.CanvasObject {
	label := widget.NewLabel("Gnomon")
	label.Alignment = fyne.TextAlignCenter
	gnomonEnabled = widget.NewRadioGroup([]string{"On", "Off"}, func(s string) {
		switch s {
		case "On":
			if rpc.Daemon.Connect {
				go menu.StartGnomon("dSlate", "boltdb", []string{}, 0, 0, nil)
			} else {
				gnomonEnabled.SetSelected("Off")
			}
		case "Off":
			menu.Gnomes.Stop("dSlate")
		default:
		}
	})
	gnomonEnabled.Horizontal = true

	cont := container.NewVBox(
		label,
		container.NewCenter(gnomonEnabled))

	return cont
}

// Gnomon search objects
func gnomonOpts() fyne.CanvasObject {
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	kv_entry := widget.NewEntry()
	kv_entry.SetPlaceHolder("Key:")

	korv := widget.NewRadioGroup([]string{"Key", "Value"}, func(s string) {})
	korv.Horizontal = true

	soru := widget.NewRadioGroup([]string{"String", "Uint64"}, func(s string) {})
	soru.Horizontal = true

	search := widget.NewButton("Search", func() {
		if menu.Gnomes.Init {
			switch korv.Selected {
			case "Key":
				switch soru.Selected {
				case "String":
					logger.Println("[secret-slate] Search results for string key "+kv_entry.Text+" on SCID "+contractInput.Text, searchByKey(contractInput.Text, kv_entry.Text, true))
					label.SetText(searchByKey(contractInput.Text, kv_entry.Text, true))
				case "Uint64":
					logger.Println("[secret-slate] Search results for uint64 key "+kv_entry.Text+" on SCID "+contractInput.Text, searchByKey(contractInput.Text, kv_entry.Text, false))
					label.SetText(searchByKey(contractInput.Text, kv_entry.Text, false))
				default:
					logger.Println("[secret-slate] Select string or uint64")
				}
			case "Value":
				switch soru.Selected {
				case "String":
					logger.Println("[secret-slate] Search results for string value "+kv_entry.Text+" on SCID "+contractInput.Text, searchByValue(contractInput.Text, kv_entry.Text, true))
					label.SetText(searchByValue(contractInput.Text, kv_entry.Text, true))
				case "Uint64":
					logger.Println("[secret-slate] Search results for uint64 value "+kv_entry.Text+" on SCID "+contractInput.Text, searchByValue(contractInput.Text, kv_entry.Text, false))
					label.SetText(searchByValue(contractInput.Text, kv_entry.Text, false))
				default:
					logger.Println("[secret-slate] Select string or uint64")
				}
			default:
				logger.Println("[secret-slate] Select key or value")
			}
		} else {
			logger.Println("[secret-slate] Gnomon not initialized")
		}

	})

	cont := container.NewVBox(
		label,
		container.NewCenter(korv),
		container.NewCenter(soru),
		container.NewAdaptiveGrid(2, kv_entry, search))

	return cont

}

// NFA objects
func nfaOpts() fyne.CanvasObject {
	label := canvas.NewText("", color.White)
	label.TextSize = 18

	file_name := widget.NewEntry()
	file_name.SetPlaceHolder("File Name:")

	start := dwidget.NewDeroEntry("", 1, 0)
	start.SetPlaceHolder("Starting at #:")
	start.Validator = validation.NewRegexp(`^\d{1,}`, "Format Not Valid")

	limit := dwidget.NewDeroEntry("", 1, 0)
	limit.ExtendBaseWidget(limit)
	limit.SetPlaceHolder("Ending at #:")
	limit.Validator = validation.NewRegexp(`^\d{1,}`, "Format Not Valid")

	fee := dwidget.NewDeroEntry("", 1, 0)
	fee.ExtendBaseWidget(fee)
	fee.SetPlaceHolder("Fee:")
	fee.Validator = validation.NewRegexp(`^\d{1,}`, "Format Not Valid")

	stop := widget.NewButton("Stop Loop", func() {
		logger.Println("[secret-slate] Stopping install loop")
		label.Text = "Stopping install loop..."
		label.Refresh()
		kill_process = true
	})

	extension := widget.NewSelect([]string{".jpg", ".png", ".gif", ".mp3", ".mp4", ".pdf", ".zip", ".7z", ".avi", ".mov", ".ogg"}, func(s string) {})
	extension.PlaceHolder = "ext"

	var install fyne.Widget
	install = widget.NewButton("Install NFAs", func() {
		if rpc.Wallet.Connect && rpc.Daemon.Connect {
			if rpc.Wallet.Address == "" || len(rpc.Wallet.Address) != 66 || rpc.Wallet.Address[0:4] != "deto" {
				dialog.NewInformation("Error", "Check wallet connection\n\ndSlate installs to testnet only", myWindow).Show()
				return
			}

			go func() {
				if fee.Validate() == nil && limit.Validate() == nil && start.Validate() == nil {
					if !process_on {
						process_on = true
						install.Hide()
						stop.Show()
						file_name.Disable()
						start.Disable()
						limit.Disable()
						fee.Disable()
						extension.Disable()

						name := file_name.Text
						lim := rpc.StringToInt(limit.Text)
						inc := rpc.StringToInt(start.Text)

						logger.Println("[secret-slate] Starting install loop for", name+strconv.Itoa(inc)+".bas", "to", name+strconv.Itoa(lim)+".bas")

						for i := 10; i > 0; i-- {
							if kill_process {
								break
							}

							label.Text = "Starting install loop in " + strconv.Itoa(i)
							label.Refresh()
							time.Sleep(1 * time.Second)
						}

						label.Text = ""
						label.Refresh()

						for {
							if kill_process {
								break
							}

							path := name + strconv.Itoa(inc) + ".bas"
							if _, err := os.Stat(path); err == nil {
								logger.Println("[secret-slate] Installing", path)
								label.Text = "Installing " + path
								label.Refresh()
							} else if errors.Is(err, os.ErrNotExist) {
								logger.Errorln("[secret-slate]", path, "Not Found")
								break
							}

							file, err := os.ReadFile(path)
							if err != nil {
								logger.Errorln("[secret-slate]", err)
								break
							}

							if tx := uploadNFA(string(file)); tx == "" {
								label.Text = ("Error installing " + path)
								break
							} else {
								logger.Println("[secret-slate] Confirming install TX")
								rpc.ConfirmTx(tx, "dSlate", 45)
							}

							inc++
							if inc > lim {
								break
							}

							logger.Println("[secret-slate] Waiting for block")
							time.Sleep(6 * time.Second)
						}

						label.Text = ""
						label.Refresh()
						install.Show()
						stop.Hide()
						file_name.Enable()
						start.Enable()
						limit.Enable()
						fee.Enable()
						extension.Enable()
						process_on = false
						kill_process = false
						logger.Println("[secret-slate] Install loop complete")

					} else {
						logger.Println("[secret-slate] Install already running")
					}
				} else {
					stop.Hide()
					logger.Println("[secret-slate] Install entries not valid")
				}
			}()
		}
	})

	update := widget.NewButton("Update Contract", func() {
		path := file_name.Text
		if _, err := os.Stat(path); err == nil {
			logger.Println("[secret-slate] Update Path", path)
			file, err := os.ReadFile(path)

			if err != nil {
				logger.Errorln("[secret-slate]", err)
				return
			}
			code := string(file)
			if code != "" {
				fe := rpc.StringToInt(fee.Text)
				updateContract(contractInput.Text, string(file), uint64(fe))
			} else {
				logger.Errorln("[secret-slate] Failed to update, code is empty string")
			}

		} else if errors.Is(err, os.ErrNotExist) {
			logger.Errorln("[secret-slate]", path, "Not Found")
		}

	})

	stop.Hide()

	wf_entry := widget.NewEntry()
	wf_entry.SetPlaceHolder("Wallet file name:")
	wf_pass := widget.NewPasswordEntry()
	wf_pass.SetPlaceHolder("Wallet file password:")

	verify_button := widget.NewButton("Verify Sign File", func() {
		if wf, err := walletapi.Open_Encrypted_Wallet(wf_entry.Text, wf_pass.Text); err == nil {
			input_file := file_name.Text
			output_file := strings.TrimSuffix(input_file, ".sign")

			if data, err := os.ReadFile(input_file); err != nil {
				logger.Errorln("[secret-slate] Cannot read input file", err)
			} else if signer, message, err := wf.CheckSignature(data); err != nil {
				logger.Errorln("[secret-slate] Signature verify failed", input_file, err)
			} else {
				logger.Println("[secret-slate] Signed by", "address", signer.String())

				if os.WriteFile(output_file, message, 0600); err != nil {
					logger.Errorln("[secret-slate] Cannot write output file", output_file, err)
				}
				logger.Println("[secret-slate] Successfully wrote message to file. please check", "file", output_file)
			}

			wf.Close_Encrypted_Wallet()
		} else {
			logger.Errorln("[secret-slate] Wallet", err)
		}
	})

	sign_button := widget.NewButton("Sign File", func() {
		if limit.Validate() == nil && start.Validate() == nil {
			if wf, err := walletapi.Open_Encrypted_Wallet(wf_entry.Text, wf_pass.Text); err == nil {
				go func() {
					if !process_on {
						process_on = true
						install.Hide()
						stop.Show()
						file_name.Disable()
						start.Disable()
						limit.Disable()
						fee.Disable()
						extension.Disable()

						ext := extension.Selected
						input := file_name.Text
						lim := rpc.StringToInt(limit.Text)
						inc := rpc.StringToInt(start.Text)

						logger.Println("[secret-slate] Starting sign loop for", input+strconv.Itoa(inc)+ext, "to", input+strconv.Itoa(lim)+ext)

						for i := 10; i > 0; i-- {
							if kill_process {
								break
							}

							label.Text = "Starting sign loop in " + strconv.Itoa(i)
							label.Refresh()
							time.Sleep(1 * time.Second)
						}

						label.Text = ""
						label.Refresh()

						for {
							if kill_process {
								break
							}

							input_file := input + strconv.Itoa(inc) + ext
							output_file := input_file + ".sign"

							if data, err := os.ReadFile(input_file); err != nil {
								logger.Errorln("[secret-slate] Cannot read input file", err)
								break
							} else if err := os.WriteFile(output_file, wf.SignData(data), 0600); err != nil {
								logger.Errorln("[secret-slate] Cannot write output file", output_file)
								break
							} else {
								logger.Println("[secret-slate] Successfully signed file. please check", output_file)
							}

							inc++

							if inc > lim {
								break
							}

							time.Sleep(6 * time.Second)
						}

						label.Text = ""
						label.Refresh()
						install.Show()
						stop.Hide()
						file_name.Enable()
						start.Enable()
						limit.Enable()
						fee.Enable()
						extension.Enable()
						process_on = false
						kill_process = false
						logger.Println("[secret-slate] Sign loop complete")

					} else {
						logger.Println("[secret-slate] Loop already running")
					}

					wf.Close_Encrypted_Wallet()
				}()
			} else {
				logger.Errorln("[secret-slate] Wallet", err)
			}
		} else {
			logger.Println("[secret-slate] Sign entries not valid")
		}
	})

	return container.NewVBox(
		container.NewBorder(nil, nil, nil, extension, file_name),
		layout.NewSpacer(),
		verify_button,
		layout.NewSpacer(),
		wf_entry,
		wf_pass,
		sign_button,
		layout.NewSpacer(),
		container.NewCenter(label),
		start,
		limit,
		layout.NewSpacer(),
		fee,
		install,
		stop,
		layout.NewSpacer(),
		update)
}
