package main

import (
	"crypto/sha256"

	"github.com/dReam-dApps/dReams/rpc"
)

// Select menu changes daemon address
func whichDaemon(s string) {
	switch s {
	case "TESTNET":
		rpc.Daemon.Rpc = "127.0.0.1:40102"
	case "SIMULATOR":
		rpc.Daemon.Rpc = "127.0.0.1:20000"
	case "CUSTOM":
		confirmPopUp() /// enter custom address in new window
	default:
		rpc.Daemon.Rpc = "127.0.0.1:10102"
	}
}

// Check if daemon is connected
func isDaemonConnected() {
	if rpc.Daemon.Connect {
		if !daemonCheckBox.Checked {
			logger.Println("[secret-slate] Daemon RPC Connected")
		}
		daemonCheckBox.SetChecked(true)
	} else {
		if debug {
			logger.Println("[secret-slate] Daemon RPC Not Connected")
		}
		currentHeight.SetText("Height:")
		if daemonCheckBox.Checked {
			daemonCheckBox.SetChecked(false)
		}

	}
}

// Check if wallet is connected
func isWalletConnected() {
	if rpc.Wallet.Connect {
		if !walletCheckBox.Checked {
			logger.Println("[secret-slate] Wallet RPC Connected")
			walletCheckBox.SetChecked(true)
		}
		GetBalance()

	} else {
		if debug {
			logger.Println("[secret-slate] Wallet RPC Not Connected")
		}
		if walletCheckBox.Checked {
			walletCheckBox.SetChecked(false)
			rpc.Wallet.Connect = false
		}
	}

	if walletCheckBox.Checked { /// if wallet is connected and any changes to inputs, show disconnected
		checkPass()
		if rpcWalletInput.Text != rpc.Wallet.Rpc {
			walletBalance.SetText("Balance: ")
			rpc.Wallet.Address = ""
			walletCheckBox.SetChecked(false)
			rpc.Wallet.Connect = false
		}
	}
}

// Check if user:pass has changed
func checkPass() {
	data := []byte(rpcLoginInput.Text)
	hash := sha256.Sum256(data)

	if hash != passHash {
		walletBalance.SetText("Balance: ")
		walletCheckBox.SetChecked(false)
		rpc.Wallet.Connect = false
	}
}
