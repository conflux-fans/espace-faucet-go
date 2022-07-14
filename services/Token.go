package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/conflux-fans/espace-faucet-go/models"
	"github.com/conflux-fans/espace-faucet-go/testData"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethSdk "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strings"
)

const URL = "https://evmtestnet.confluxrpc.com"
const privKey = "36b37e2fb26276bc6344305e1fc6b829f6a8fa0e151b3c1b9c82e251d8f2c2f3"
//var privKey = os.Getenv("PRIVATE_KEY")
//var URL = os.Getenv("URL")
//var password = os.Getenv("PASSWORD")
var rpcClient, _  = rpc.Dial(URL)
const abiJSON = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_initialAmount\",\"type\":\"uint256\"},{\"name\":\"_tokenName\",\"type\":\"string\"},{\"name\":\"_decimalUnits\",\"type\":\"uint8\"},{\"name\":\"_tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

func SendCFX(toAddress string) (string, error) {
	client := ethSdk.NewClient(rpcClient)
	signedTx, err := createTx(client, toAddress, nil)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	var isPending bool
	for isPending {
		_, isPending, err =  client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return "", err
		}
	}

	return signedTx.Hash().String(), nil
}

func SendERC20(request models.ERC20) (string, error)  {
	client := ethSdk.NewClient(rpcClient)
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return "", err
	}

	fromAddr, _, err := getFromAddress()
	if err != nil {
		return "", err
	}

	input, err := parsed.Pack("transfer", fromAddr, big.NewInt(1000000000))
	signedTx, err := createTx(client, request.ContractAddress, input)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	var isPending bool
	for isPending {
		_, isPending, err =  client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return "", err
		}
	}
	return signedTx.Hash().String(), nil
}

func createTx(client *ethSdk.Client, toAddress string, data []byte) (*types.Transaction, error){
	fromAddr, fromPrivKey, err := getFromAddress()
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(toAddress)

	var gasLimit uint64 = 300000
	//amount := big.NewInt(1000000000000000000)
	amount := big.NewInt(0)

	err = checkBalance(client,fromAddr, amount)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), *fromAddr)
	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func getFromAddress() (*common.Address, *ecdsa.PrivateKey, error) {
	fromPrivkey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, nil, err
	}
	fromPubkey := fromPrivkey.PublicKey
	fromAddr := crypto.PubkeyToAddress(fromPubkey)
	return &fromAddr, fromPrivkey, nil
}

func checkBalance(client *ethSdk.Client, fromAddr *common.Address, amount *big.Int) error {
	balance, err := client.BalanceAt(context.Background(), *fromAddr, nil)
	if err != nil {
		return err
	}
	if balance.Cmp(amount) == -1 {
		return errors.New("Insufficient balance")
	}
	return nil
}

func DeployERC20() (*models.Resp, error){
	client := ethSdk.NewClient(rpcClient)
	fromAddr, fromPrivKey, err := getFromAddress()
	if err != nil {
		return nil, err
	}
	nonce, err := client.PendingNonceAt(context.Background(), *fromAddr)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth := bind.NewKeyedTransactor(fromPrivKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	oneCfx := new(big.Int).Mul(big.NewInt(1e9), big.NewInt(1e9))
	addr, tx, _, err := testData.DeployToken(auth, client, new(big.Int).Mul(big.NewInt(1000000), oneCfx), "biu", uint8(10), "BIU")
	return &models.Resp{
		ContractAddr: addr.String(),
		Hash: tx.Hash().String(),
	}, nil

}

//func QueryERC20() (error){
//	client := ethSdk.NewClient(rpcClient)
//	fromAddr, _, err := getFromAddress()
//	if err != nil {
//		return  err
//	}
//
//	token, err := testData.NewToken(common.HexToAddress("0xb2193d7f5b072978585e75cf5ed2ec06fc2340f7"), client)
//	if err != nil {
//		return  err
//	}
//
//	num, err := token.BalanceOf(nil, *fromAddr)
//	if err != nil {
//		return  err
//	}
//	return nil
//}