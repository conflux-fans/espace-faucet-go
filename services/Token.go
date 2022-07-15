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
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

func SendCFX(toAddress string) (string, error) {
	client := ethSdk.NewClient(getClient())
	oneCfx := big.NewInt(1000000000000000000)

	amount := new(big.Int).Mul(big.NewInt(viper.GetInt64("sendcfx")), oneCfx)
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
		_, isPending, err =  client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return "", err
		}
	}

	return signedTx.Hash().String(), nil
}

func SendERC20(request models.ERC20) (string, error)  {
	client := ethSdk.NewClient(getClient())
	erc20Data := viper.GetStringMap("erc20")

	token := erc20Data[request.Name].(map[string]interface{})
	if token == nil {
		return "", errors.New("Unsupported token")
	}

	parsed, err := abi.JSON(strings.NewReader(token["abijson"].(string)))
	if err != nil {
		return "", err
	}

	fromAddr, _, err := getFromAddress()
	if err != nil {
		return "", err
	}
	oneToken := big.NewInt(10000000000)
	amount := new(big.Int).Mul(big.NewInt(int64(token["value"].(float64))), oneToken)

	input, err := parsed.Pack("transfer", fromAddr, amount)
	signedTx, err := createTx(client, token["address"].(string), input, big.NewInt(0))
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
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

func createTx(client *ethSdk.Client, toAddress string, data []byte, amount *big.Int) (*types.Transaction, error){
	fromAddr, fromPrivKey, err := getFromAddress()
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(toAddress)

	var gasLimit uint64 = 300000

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
	var privKey = viper.GetString("providerkey")
	fromPrivkey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		panic(err)
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
	client := ethSdk.NewClient(getClient())
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
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	oneCfx := new(big.Int).Mul(big.NewInt(1e9), big.NewInt(1e9))
	addr, tx, _, err := testData.DeployToken(auth, client, new(big.Int).Mul(big.NewInt(1000000), oneCfx), "k", uint8(10), "KKK")
	return &models.Resp{
		ContractAddr: addr.String(),
		Hash: tx.Hash().String(),
	}, nil
}

func getClient() *rpc.Client {
	var URL = viper.GetString("espaceurl")
	rpcClient, _ := rpc.Dial(URL)
	return rpcClient
}
