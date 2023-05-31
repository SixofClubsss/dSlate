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

var passHash [32]byte

// Get wallet address with user:pass auth
func GetAddress() {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	var result *dero.GetAddress_Result
	if err := rpcClientW.CallFor(ctx, &result, "GetAddress"); err != nil {
		rpc.Wallet.Connect = false
		rpc.Wallet.Rpc = ""
		walletCheckBox.SetChecked(false)
		log.Println("[GetAddress]", err)
		return
	}

	if len(result.Address) == 66 {
		rpc.Wallet.Connect = true
		rpc.Wallet.Address = result.Address
		walletCheckBox.SetChecked(true)
		log.Println("[dSlate] Wallet Connected")
		log.Println("[dSlate] Dero Address: " + result.Address)
		passHash = sha256.Sum256([]byte(rpc.Wallet.UserPass))
	}
}

// Get wallet Dero balance
func GetBalance() {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	var result *dero.GetBalance_Result
	if err := rpcClientW.CallFor(ctx, &result, "GetBalance"); err != nil {
		log.Println("[GetBalance]", err)
		return
	}

	atomic := float64(result.Unlocked_Balance) /// unlocked balance in atomic units
	div := atomic / 100000
	str := strconv.FormatFloat(div, 'f', 5, 64)
	walletBalance.SetText("Balance: " + str)
}

// Update existing contracts with 'UpdateCode' entrypoint
func updateContract(scid, code string, fee uint64) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "UpdateCode"}
	arg2 := dero.Argument{Name: "code", DataType: "S", Value: code}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	var addr string
	switch rpc.Daemon.Rpc {
	case "127.0.0.1:10102":
		addr = "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn"
	case "127.0.0.1:40402":
		addr = "deto1qyre7td6x9r88y4cavdgpv6k7lvx6j39lfsx420hpvh3ydpcrtxrxqg8v8e3z"
	case "127.0.0.1:20000":
		addr = "deto1qyre7td6x9r88y4cavdgpv6k7lvx6j39lfsx420hpvh3ydpcrtxrxqg8v8e3z"
	default:
		addr = "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn"
	}

	t1 := dero.Transfer{
		Destination: addr,
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		log.Println("[updateContract]", err)
		return
	}

	log.Println("[updateContract] Update TX:", txid)

	return
}

// Upload a NFA SC to testnet by string
func uploadNFA(code string) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	txid := dero.Transfer_Result{}
	t1 := dero.Transfer{
		Destination: "deto1qyre7td6x9r88y4cavdgpv6k7lvx6j39lfsx420hpvh3ydpcrtxrxqg8v8e3z",
		Amount:      0,
	}

	params := &dero.Transfer_Params{
		Transfers: []dero.Transfer{t1},
		SC_Code:   code,
		SC_Value:  0,
		SC_RPC:    dero.Arguments{},
		Ringsize:  2,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		log.Println("[uploadNFA]", err)
		return
	}

	log.Println("[uploadNFA] TXID:", txid)

	return txid.TXID
}
