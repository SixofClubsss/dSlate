package main

import (
	"crypto/sha256"
	"fmt"
)

const (
	pre  = "http://"
	suff = "/json_rpc"
)

func whichDaemon(s string) { /// select menu changes dameon address
	switch s {
	case "TESTNET":
		daemonAddress = pre + "127.0.0.1:40102" + suff
	case "SIMULATOR":
		daemonAddress = pre + "127.0.0.1:20000" + suff
	case "CUSTOM":
		confirmPopUp() /// enter custom address in new window
	default:
		daemonAddress = pre + "127.0.0.1:10102" + suff
	}
}

func isDaemonConnected() { /// check if daemon is connected
	if daemonConnectBool {
		if !daemonCheckBox.Checked {
			fmt.Println("Daemon Connected")
		}
		daemonCheckBox.SetChecked(true)
	} else {
		fmt.Println("Daemon Not Connected")
		currentHeight.SetText("Height:")
		if daemonCheckBox.Checked {
			daemonCheckBox.SetChecked(false)
		}

	}

}

func isWalletConnected() { /// check if wallet is connected
	if walletConnectBool {
		if !walletCheckBox.Checked {
			fmt.Println("Wallet Connected")
			walletCheckBox.SetChecked(true)
		}
		GetBalance()

	} else {
		fmt.Println("Wallet Not Connected")
		if walletCheckBox.Checked {
			walletCheckBox.SetChecked(false)
			walletConnectBool = false
		}
	}

	if walletCheckBox.Checked { /// if wallet is connected and any changes to inputs, show disconnected
		checkPass()
		if pre+rpcWalletInput.Text+suff != walletAddress {
			walletBalance.SetText("Balance: ")
			walletAddress = ""
			walletCheckBox.SetChecked(false)
			walletConnectBool = false
		}
	}

}

func checkPass() { /// check if user:pass has changed
	data := []byte(rpcLoginInput.Text)
	hash := sha256.Sum256(data)

	if hash != passHash {
		walletBalance.SetText("Balance: ")
		walletCheckBox.SetChecked(false)
		walletConnectBool = false
	}
}
