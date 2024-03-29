package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
)

const (
	DAEMON_MAINNET_DEFAULT   = "127.0.0.1:10102"
	DAEMON_TESTNET_DEFAULT   = "127.0.0.1:40402"
	DAEMON_SIMULATOR_DEFAULT = "127.0.0.1:20000"
)

// Gets current height and displays
func GetHeight() {
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
	defer cancel()

	var result *dero.Daemon_GetHeight_Result
	if err := rpcClientD.CallFor(ctx, &result, "DERO.GetHeight"); err != nil {
		return
	}

	h := result.Height
	str := strconv.FormatUint(h, 10)
	currentHeight.SetText("Height: " + str)
}

// Search SC using getsc method
func getSC(p *dero.GetSC_Params) {
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
	defer cancel()

	var result *dero.GetSC_Result
	if err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", p); err != nil {
		logger.Errorln("[getSC]", err)
		return
	}

	bal := result.Balances /// retrieve all sc balances
	balM, _ := json.Marshal(bal)
	balances := string(balM)

	string_keys := SortStringMap(result.VariableStringKeys) /// retrieve all sc string keys, use result.VariableStringKeys["KEY"] for single value

	uint_keys := SortUintMap(result.VariableUint64Keys) /// retrieve all sc uint keys use result.VariableUint64Keys[0] for single value

	go searchPopUp(balances, string_keys, uint_keys, result.Code) /// displays results

}

// Get SC code and print result in terminal
func getSCcode(scid string) {
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
	defer cancel()

	var result *dero.GetSC_Result
	params := dero.GetSC_Params{
		SCID:      scid,
		Code:      true,
		Variables: false,
	}

	if err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params); err != nil {
		logger.Errorln("[getSCcode]", err)
		return
	}

	fmt.Println(result.Code)
}

// Switch for interface type to string
func findKey(i interface{}) (text string) {
	switch v := i.(type) {
	case uint64:
		text = strconv.Itoa(int(v))
	case string:
		text = v
	case float64:
		text = strconv.Itoa(int(v))
	default:

	}

	return
}

func SortStringMap(m map[string]interface{}) (str string) {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if m[k] == m["C"] {
			/// skipping C
		} else {
			str = str + k + " " + findKey(m[k]) + " \n"
		}
	}

	return
}

func SortUintMap(m map[uint64]interface{}) (str string) {
	keys := make([]uint64, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		str = str + strconv.Itoa(int(k)) + " " + findKey(m[k]) + " \n"

	}

	return
}
