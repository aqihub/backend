package incentivizer

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Global variables
var (
	tokenABI     abi.ABI
	privateKey   *ecdsa.PrivateKey
	tokenAddress common.Address
)

func TransferToken(contractAddress string, amount float64) (string, error) {

	txHash, err := transferTokens(contractAddress, amount)
	if err != nil {
		log.Printf("Error transferring tokens: %v", err)
		return err.Error(), err
	}

	return txHash, nil
}

// Helper function to convert amount to wei (assuming 18 decimals)
func toWei(amount float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18))
	result := new(big.Int)
	wei.Int(result)
	return result
}

// Transfer tokens function
func transferTokens(toAddress string, amount float64) (string, error) {

	apiKey := os.Getenv("ALCHEMY_API_KEY")
	rpcURL := fmt.Sprintf(os.Getenv("ALCHEMY_ETH_NETWORK")+"/%s", apiKey)

	// Connect to Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	// Parse private key
	privateKeyHex := os.Getenv("SIGNER_PRIVATE_KEY")
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}
	privateKey, err = crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Parse token address
	tokenAddress = common.HexToAddress(os.Getenv("AQI_ERC20_ADDRESS"))

	// Parse ABI
	tokenABI, err = abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf("Failed to parse token ABI: %v", err)
	}

	// Create contract instance
	tokenContract := bind.NewBoundContract(tokenAddress, tokenABI, client, client, client)

	// Get auth for transaction
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(43113))
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %v", err)
	}

	// Set gas limit
	auth.GasLimit = uint64(100000)

	// Convert amount to wei
	amountWei := toWei(amount)

	// Send transaction
	tx, err := tokenContract.Transact(auth, "transfer", common.HexToAddress(toAddress), amountWei)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return tx.Hash().Hex(), nil
}
