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

func uploadContract(code string, fee uint64) error {
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
