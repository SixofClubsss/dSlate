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
		log.Println("[GetAddress]", err)
		return nil
	}

	address := len(result.Address)
	if address == 66 {
		walletConnect = true
		walletCheckBox.SetChecked(true)
		log.Println("[dSlate] Wallet Connected")
		log.Println("[dSlate] Dero Address: " + result.Address)
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
		log.Println("[GetBalance]", err)
		return nil
	}

	atomic := float64(result.Unlocked_Balance) /// unlocked balance in atomic units
	div := atomic / 100000
	str := strconv.FormatFloat(div, 'f', 5, 64)
	walletBalance.SetText("Balance: " + str)

	return err
}

func uploadContract(code string, fee uint64) error { /// install new contract
	rpcClientW, ctx, cancel := rpc.SetWalletClient(walletAddress, rpcLoginInput.Text)
	defer cancel()

	txid := dero.Transfer_Result{}

	params := &dero.Transfer_Params{
		Transfers: []dero.Transfer{},
		SC_Code:   code,
		SC_Value:  0,
		SC_RPC:    dero.Arguments{},
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println("[uploadContract]", err)
		return nil
	}

	log.Println("[uploadContract] TXID:", txid)

	return err
}

func updateContract(scid, code string, fee uint64) error { /// update existing contracts with 'UpdateCode' entrypoint
	rpcClientW, ctx, cancel := rpc.SetWalletClient(walletAddress, rpcLoginInput.Text)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "UpdateCode"}
	arg2 := dero.Argument{Name: "code", DataType: "S", Value: code}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	var addr string
	switch daemonAddress {
	case "127.0.0.1:10102":
		addr = "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn"
	case "127.0.0.1:40102":
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

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println("[updateContract]", err)
		return nil
	}

	log.Println("[updateContract] Update TX:", txid)

	return err
}
