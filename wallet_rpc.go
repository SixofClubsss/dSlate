package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

var passHash [32]byte

const (
	WALLET_MAINNET_DEFAULT   = "http://127.0.0.1:10103/json_rpc"
	WALLET_TESTNET_DEFAULT   = "http://127.0.0.1:40403/json_rpc"
	WALLET_SIMULATOR_DEFAULT = "http://127.0.0.1:30000/json_rpc"
)

var rpcClientW = jsonrpc.NewClient(WALLET_MAINNET_DEFAULT)
var walletConnectBool bool

func GetAddress() error { /// get address with user:pass
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW = jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})
	var result *rpc.GetAddress_Result
	err := rpcClientW.CallFor(ctx, &result, "GetAddress")

	if err != nil {
		walletConnectBool = false
		walletCheckBox.SetChecked(false)
		fmt.Println(err)
		return nil
	}

	address := len(result.Address)
	if address == 66 {
		walletConnectBool = true
		walletCheckBox.SetChecked(true)
		fmt.Println("Wallet Connected")
		fmt.Println("Dero Address:" + result.Address)
		data := []byte(rpcLoginInput.Text)
		passHash = sha256.Sum256(data)
	}

	return err
}

func GetBalance() error { /// get wallet balance
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW = jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})
	var result *rpc.GetBalance_Result
	err := rpcClientW.CallFor(ctx, &result, "GetBalance")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	atomic := float64(result.Unlocked_Balance) /// unlocked balance in atomic units
	div := atomic / 100000
	str := strconv.FormatFloat(div, 'f', 5, 64)
	walletBalance.SetText("Balance: " + str)

	return err
}
