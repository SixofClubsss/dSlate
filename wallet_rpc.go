package main

import (
	"crypto/sha256"
	"log"
	"strconv"

	"github.com/SixofClubsss/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
)

const (
	WALLET_MAINNET_DEFAULT   = "127.0.0.1:10103"
	WALLET_TESTNET_DEFAULT   = "127.0.0.1:40403"
	WALLET_SIMULATOR_DEFAULT = "127.0.0.1:30000"
)

var (
	walletAddress string
	walletConnect bool
	passHash      [32]byte
)

func GetAddress() error { /// get address with user:pass
	rpcClientW, ctx, cancel := rpc.SetWalletClient(walletAddress, rpcLoginInput.Text)
	defer cancel()

	var result *dero.GetAddress_Result
	err := rpcClientW.CallFor(ctx, &result, "GetAddress")

	if err != nil {
		walletConnect = false
		walletCheckBox.SetChecked(false)
		log.Println(err)
		return nil
	}

	address := len(result.Address)
	if address == 66 {
		walletConnect = true
		walletCheckBox.SetChecked(true)
		log.Println("Wallet Connected")
		log.Println("Dero Address:" + result.Address)
		data := []byte(rpcLoginInput.Text)
		passHash = sha256.Sum256(data)
	}

	return err
}

func GetBalance() error { /// get wallet balance
	rpcClientW, ctx, cancel := rpc.SetWalletClient(walletAddress, rpcLoginInput.Text)
	defer cancel()

	var result *dero.GetBalance_Result
	err := rpcClientW.CallFor(ctx, &result, "GetBalance")

	if err != nil {
		log.Println(err)
		return nil
	}

	atomic := float64(result.Unlocked_Balance) /// unlocked balance in atomic units
	div := atomic / 100000
	str := strconv.FormatFloat(div, 'f', 5, 64)
	walletBalance.SetText("Balance: " + str)

	return err
}
