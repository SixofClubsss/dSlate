package main

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"

	"github.com/SixofClubsss/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
)

const (
	DAEMON_MAINNET_DEFAULT   = "127.0.0.1:10102"
	DAEMON_TESTNET_DEFAULT   = "127.0.0.1:40402"
	DAEMON_SIMULATOR_DEFAULT = "127.0.0.1:20000"
)

var (
	daemonAddress string
	daemonConnect bool
)

func Ping() error { /// ping blockchain for connection
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(daemonAddress)
	defer cancel()

	var result string
	err := rpcClientD.CallFor(ctx, &result, "DERO.Ping")
	if err != nil {
		daemonConnect = false
		return nil
	}

	if result == "Pong " {
		daemonConnect = true
	} else {
		daemonConnect = false
	}

	return err
}

func GetHeight() error { /// get current height and displays
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(daemonAddress)
	defer cancel()

	var result *dero.Daemon_GetHeight_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetHeight")

	if err != nil {
		return nil
	}
	h := result.Height
	log.Printf("Daemon Height: %d \n", h)
	str := strconv.FormatUint(h, 10)
	currentHeight.SetText("Height: " + str)

	return err
}

func getSC(p *dero.GetSC_Params) error { /// search sc using getsc method
	rpcClientD, ctx, cancel := rpc.SetDaemonClient(daemonAddress)
	defer cancel()

	var result *dero.GetSC_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", p)

	if err != nil {
		log.Println(err)
		return nil
	}

	bal := result.Balances /// retrieve all sc balances
	balM, _ := json.Marshal(bal)
	balances := string(balM)

	string_keys := SortStringMap(result.VariableStringKeys) /// retrieve all sc string keys, use result.VariableStringKeys["KEY"] for single value

	uint_keys := SortUintMap(result.VariableUint64Keys) /// retrieve all sc uint keys use result.VariableUint64Keys[0] for single value

	go searchPopUp(balances, string_keys, uint_keys, result.Code) /// displays results

	return err
}

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
