package main

import (
	"crypto/sha256"
	"log"
)

func whichDaemon(s string) { /// select menu changes dameon address
	switch s {
	case "TESTNET":
		daemonAddress = "127.0.0.1:40102"
	case "SIMULATOR":
		daemonAddress = "127.0.0.1:20000"
	case "CUSTOM":
		confirmPopUp() /// enter custom address in new window
	default:
		daemonAddress = "127.0.0.1:10102"
	}
}

func isDaemonConnected() { /// check if daemon is connected
	if daemonConnect {
		if !daemonCheckBox.Checked {
			log.Println("[dSlate] Daemon RPC Connected")
		}
		daemonCheckBox.SetChecked(true)
	} else {
		if debug {
			log.Println("[dSlate] Daemon RPC Not Connected")
		}
		currentHeight.SetText("Height:")
		if daemonCheckBox.Checked {
			daemonCheckBox.SetChecked(false)
		}

	}
}

func isWalletConnected() { /// check if wallet is connected
	if walletConnect {
		if !walletCheckBox.Checked {
			log.Println("[dSlate] Wallet RPC Connected")
			walletCheckBox.SetChecked(true)
		}
		GetBalance()

	} else {
		if debug {
			log.Println("[dSlate] Wallet RPC Not Connected")
		}
		if walletCheckBox.Checked {
			walletCheckBox.SetChecked(false)
			walletConnect = false
		}
	}

	if walletCheckBox.Checked { /// if wallet is connected and any changes to inputs, show disconnected
		checkPass()
		if rpcWalletInput.Text != walletAddress {
			walletBalance.SetText("Balance: ")
			walletAddress = ""
			walletCheckBox.SetChecked(false)
			walletConnect = false
		}
	}
}

func checkPass() { /// check if user:pass has changed
	data := []byte(rpcLoginInput.Text)
	hash := sha256.Sum256(data)

	if hash != passHash {
		walletBalance.SetText("Balance: ")
		walletCheckBox.SetChecked(false)
		walletConnect = false
	}
}
