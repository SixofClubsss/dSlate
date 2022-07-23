package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

const (
	DAEMON_MAINNET_DEFAULT   = "http://127.0.0.1:10102/json_rpc"
	DAEMON_TESTNET_DEFAULT   = "http://127.0.0.1:40102/json_rpc"
	DAEMON_SIMULATOR_DEFAULT = "http://127.0.0.1:20000/json_rpc"
)

var rpcClientD = jsonrpc.NewClient(DAEMON_MAINNET_DEFAULT) /// daemon default to mainnet

var daemonConnectBool bool

func Ping() error { /// ping blockchain for connection
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD = jsonrpc.NewClient(daemonAddress)
	var result string
	err := rpcClientD.CallFor(ctx, &result, "DERO.Ping")
	if err != nil {
		daemonConnectBool = false
		return nil
	}

	if result == "Pong " {
		daemonConnectBool = true
	} else {
		daemonConnectBool = false
	}

	return err
}

func GetHeight() error { /// get current height and displays
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD = jsonrpc.NewClient(daemonAddress)
	var result *rpc.Daemon_GetHeight_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetHeight")

	if err != nil {
		return nil
	}
	h := result.Height
	fmt.Printf("Daemon Height: %d \n", h)
	str := strconv.FormatUint(result.Height, 10)
	currentHeight.SetText("Height: " + str)

	return err

}

func getSC(p *rpc.GetSC_Params) error { /// search sc using getsc method
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD := jsonrpc.NewClient(daemonAddress)
	var result *rpc.GetSC_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	bal := result.Balances /// retrieve sc balances
	balM, _ := json.Marshal(bal)
	strB := string(balM)

	sKeys := result.VariableStringKeys /// retrieve sc string keys
	sKeysM, _ := json.Marshal(sKeys)
	strK := string(sKeysM)

	uintKeys := result.VariableUint64Keys /// retrieve sc uint keys
	uintKeysM, _ := json.Marshal(uintKeys)
	uintK := string(uintKeysM)

	searchPopUp(strB, strK, uintK, result.Code) /// displays results

	return err
}
