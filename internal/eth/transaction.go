package eth

import (
	"encoding/hex"
	"errors"
	"strings"

	"flashblock/internal/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Errors
var (
	ErrInvalidRawTx = errors.New("invalid raw transaction format")
)

// DecodeRawTransaction decodes a raw Ethereum transaction from hex format
func DecodeRawTransaction(rawTxHex string) (*types.Transaction, error) {
	// Remove "0x" prefix if present
	if strings.HasPrefix(rawTxHex, "0x") {
		rawTxHex = rawTxHex[2:]
	}

	// Decode hex string to bytes
	rawTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		return nil, err
	}

	// Decode RLP encoded transaction
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(rawTxBytes, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// ConvertToModelTransaction converts an Ethereum transaction to a model.Transaction
func ConvertToModelTransaction(ethTx *types.Transaction, rawTxHex string) (*model.Transaction, error) {
	var from string
	signer := types.LatestSignerForChainID(ethTx.ChainId())
	sender, err := types.Sender(signer, ethTx)
	if err == nil {
		from = sender.Hex()
	}

	var to string
	if ethTx.To() != nil {
		to = ethTx.To().Hex()
	}

	// Extract transaction data
	data := ethTx.Data()
	value := ethTx.Value()
	gasPrice := ethTx.GasPrice()
	gasLimit := ethTx.Gas()
	nonce := ethTx.Nonce()

	return model.NewEthereumTransaction(
		from,
		to,
		value,
		gasPrice,
		gasLimit,
		nonce,
		data,
		rawTxHex,
	), nil
}

// ParseRawTransaction parses a raw transaction hex string and returns a model.Transaction
func ParseRawTransaction(rawTxHex string) (*model.Transaction, error) {
	// Decode the raw transaction
	ethTx, err := DecodeRawTransaction(rawTxHex)
	if err != nil {
		return nil, err
	}

	// Convert to our transaction model
	return ConvertToModelTransaction(ethTx, rawTxHex)
}

// RecoverSender attempts to recover the sender address from a raw transaction
func RecoverSender(rawTxHex string) (common.Address, error) {
	tx, err := DecodeRawTransaction(rawTxHex)
	if err != nil {
		return common.Address{}, err
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	return types.Sender(signer, tx)
}
