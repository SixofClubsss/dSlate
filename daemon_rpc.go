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
	DAEMON_TESTNET_DEFAULT   = "http://127.0.0.1:40402/json_rpc"
	DAEMON_SIMULATOR_DEFAULT = "http://127.0.0.1:20000/json_rpc"
)

var rpcClientD = jsonrpc.NewClient(DAEMON_MAINNET_DEFAULT) /// daemon default to mainnet

var daemonConnectBool bool

var card1 int
var card2 int

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

func getSC_ex1() error { /// search for two cards
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD := jsonrpc.NewClient(daemonAddress)
	var result *rpc.GetSC_Result
	p := &rpc.GetSC_Params{
		SCID:      "c8b015cffe9dccec02541792e6142ac3dca9b66392315525ad34a9f4df8d65d9",
		Code:      false,
		Variables: true,
	}
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if result.VariableStringKeys["Player 1 Card:"] != nil {
		card1 = int(result.VariableStringKeys["Player 1 Card:"].(float64))
		card2 = int(result.VariableStringKeys["Player 2 Card:"].(float64))
		fmt.Println("Card 1 is:", card1)
		fmt.Println("Card 2 is:", card2)
	} else {
		card1 = 0
		card2 = 0
	}

	return err
}
