package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/conflux-fans/espace-faucet-go/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethSdk "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

func SendCFX(toAddress string) (string, error) {
	client := ethSdk.NewClient(getClient())

	// 使用 float64 处理配置，支持 1.5 CFX 等情况
	amount := toWei(viper.GetFloat64("sendcfx"), 18)
	signedTx, err := createTx(client, toAddress, nil, amount)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	var isPending bool
	for isPending {
		_, isPending, err = client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return "", err
		}
	}

	return signedTx.Hash().String(), nil
}

func SendERC20(request models.ERC20) (string, error) {
	client := ethSdk.NewClient(getClient())
	erc20Data := viper.GetStringMap("erc20")
	erc20AbiStr := viper.GetString("abijson")

	token := erc20Data[request.Name].(map[string]interface{})
	if token == nil {
		return "", errors.New("Unsupported token")
	}

	parsed, err := abi.JSON(strings.NewReader(erc20AbiStr))
	if err != nil {
		return "", err
	}

	decimal := token["decimals"].(int)
	result := toWei(token["value"].(float64), decimal)

	input, err := parsed.Pack("transfer", common.HexToAddress(request.Address), result)
	signedTx, err := createTx(client, token["address"].(string), input, big.NewInt(0))
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		// panic(err)
		return "", err
	}

	var isPending bool
	for isPending {
		_, isPending, err = client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return "", err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return signedTx.Hash().String(), nil
}

func createTx(client *ethSdk.Client, toAddress string, data []byte, amount *big.Int) (*types.Transaction, error) {
	fromAddr, fromPrivKey, err := getFromAddress()
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(toAddress)

	var gasLimit uint64 = 300000

	err = checkBalance(client, fromAddr, amount)
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
	var privKey = viper.GetString("providerkey")
	fromPrivkey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		// panic(err)
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

func getClient() *rpc.Client {
	var URL = viper.GetString("espaceurl")
	rpcClient, _ := rpc.Dial(URL)
	return rpcClient
}

// toWei 将 float64 金额根据 decimals 转换为 *big.Int (Wei)
func toWei(amount float64, decimals int) *big.Int {
	// 计算 10^decimals
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// 使用 big.Float 处理浮点数乘法以保证精度
	fAmount := new(big.Float).SetFloat64(amount)
	fMultiplier := new(big.Float).SetInt(multiplier)
	fResult := new(big.Float).Mul(fAmount, fMultiplier)

	// 转换为 big.Int
	result := new(big.Int)
	fResult.Int(result)
	return result
}
