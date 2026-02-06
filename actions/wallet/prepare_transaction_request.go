package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// DefaultParameters is the default set of parameters to prepare for a transaction request.
// This mirrors viem's `defaultParameters` export.
var DefaultParameters = []string{
	"blobVersionedHashes",
	"chainId",
	"fees",
	"gas",
	"nonce",
	"type",
}

// PrepareTransactionRequestParameters contains the parameters for the PrepareTransactionRequest action.
// This mirrors viem's PrepareTransactionRequestParameters type.
type PrepareTransactionRequestParameters struct {
	// Account is the account to prepare the transaction for. If nil, uses the client's account.
	Account Account

	// Chain optionally overrides the client's chain.
	Chain *viemchain.Chain

	// Parameters is the list of parameters to prepare.
	// Defaults to DefaultParameters: ["blobVersionedHashes", "chainId", "fees", "gas", "nonce", "type"].
	Parameters []string

	// ChainID optionally specifies the chain ID.
	ChainID *int64

	// Transaction fields
	AccessList           []formatters.AccessListItem       `json:"accessList,omitempty"`
	AuthorizationList    []transaction.SignedAuthorization `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string                          `json:"blobVersionedHashes,omitempty"`
	Blobs                []string                          `json:"blobs,omitempty"`
	Data                 string                            `json:"data,omitempty"`
	From                 string                            `json:"from,omitempty"`
	Gas                  *big.Int                          `json:"gas,omitempty"`
	GasPrice             *big.Int                          `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     *big.Int                          `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int                          `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int                          `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int                              `json:"nonce,omitempty"`
	To                   string                            `json:"to,omitempty"`
	Type                 formatters.TransactionType        `json:"type,omitempty"`
	Value                *big.Int                          `json:"value,omitempty"`
}

// PrepareTransactionRequestReturnType is the return type for PrepareTransactionRequest.
// It is a fully-prepared transaction request with all required fields populated.
type PrepareTransactionRequestReturnType = *PrepareTransactionRequestParameters

// PrepareTransactionRequest prepares a transaction request for signing.
//
// This action fills in missing fields such as nonce, chainId, gas, fees, and type
// based on the current network state. This mirrors viem's `prepareTransactionRequest` action.
//
// Example:
//
//	prepared, err := wallet.PrepareTransactionRequest(ctx, client, wallet.PrepareTransactionRequestParameters{
//	    To:    "0x0000000000000000000000000000000000000000",
//	    Value: big.NewInt(1),
//	})
//
// Example with account hoisting:
//
//	prepared, err := wallet.PrepareTransactionRequest(ctx, client, wallet.PrepareTransactionRequestParameters{
//	    Account: myAccount,
//	    To:      "0x0000000000000000000000000000000000000000",
//	    Value:   big.NewInt(1),
//	})
func PrepareTransactionRequest(ctx context.Context, client Client, params PrepareTransactionRequestParameters) (PrepareTransactionRequestReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}

	// Resolve chain
	ch := params.Chain
	if ch == nil {
		ch = client.Chain()
	}

	// Resolve parameters to prepare
	parameters := params.Parameters
	if len(parameters) == 0 {
		parameters = DefaultParameters
	}

	// Set from address if account is available
	if account != nil && params.From == "" {
		params.From = account.Address().Hex()
	}

	// Helper to resolve chain ID
	resolveChainID := func() (int64, error) {
		if params.ChainID != nil {
			return *params.ChainID, nil
		}
		if ch != nil {
			return ch.ID, nil
		}
		chainID, err := public.GetChainID(ctx, client)
		if err != nil {
			return 0, err
		}
		return int64(chainID), nil
	}

	// ---------- Fill nonce ----------
	if containsParam(parameters, "nonce") && params.Nonce == nil && account != nil {
		nonce, err := public.GetTransactionCount(ctx, client, public.GetTransactionCountParameters{
			Address:  common.HexToAddress(account.Address().Hex()),
			BlockTag: "pending",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
		n := int(nonce)
		params.Nonce = &n
	}

	// ---------- Fill chainId ----------
	if containsParam(parameters, "chainId") {
		chainID, err := resolveChainID()
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %w", err)
		}
		params.ChainID = &chainID
	}

	// ---------- Fill type ----------
	if (containsParam(parameters, "fees") || containsParam(parameters, "type")) && params.Type == "" {
		// Try to infer the transaction type from existing fields
		txType := inferTransactionType(&params)
		if txType != "" {
			params.Type = txType
		} else {
			// Determine from the network: if the latest block has baseFeePerGas, it's EIP-1559
			block, blockErr := public.GetBlock(ctx, client, public.GetBlockParameters{
				BlockTag: "latest",
			})
			if blockErr != nil {
				return nil, fmt.Errorf("failed to get latest block: %w", blockErr)
			}
			if block != nil && block.BaseFeePerGas != nil {
				params.Type = formatters.TransactionTypeEIP1559
			} else {
				params.Type = formatters.TransactionTypeLegacy
			}
		}
	}

	// ---------- Fill fees ----------
	if containsParam(parameters, "fees") {
		if params.Type != formatters.TransactionTypeLegacy && params.Type != formatters.TransactionTypeEIP2930 {
			// EIP-1559 fees
			if params.MaxFeePerGas == nil || params.MaxPriorityFeePerGas == nil {
				fees, feeErr := public.EstimateFeesPerGas(ctx, client, public.EstimateFeesPerGasParameters{
					Type: public.FeeValuesTypeEIP1559,
				})
				if feeErr != nil {
					return nil, fmt.Errorf("failed to estimate fees: %w", feeErr)
				}

				// Check if existing maxFeePerGas is too low
				if params.MaxPriorityFeePerGas == nil && params.MaxFeePerGas != nil &&
					params.MaxFeePerGas.Cmp(fees.MaxPriorityFeePerGas) < 0 {
					return nil, fmt.Errorf(
						"maxFeePerGas (%s) cannot be lower than maxPriorityFeePerGas (%s)",
						params.MaxFeePerGas.String(), fees.MaxPriorityFeePerGas.String(),
					)
				}

				if params.MaxPriorityFeePerGas == nil {
					params.MaxPriorityFeePerGas = fees.MaxPriorityFeePerGas
				}
				if params.MaxFeePerGas == nil {
					params.MaxFeePerGas = fees.MaxFeePerGas
				}
			}
		} else {
			// Legacy fees
			if params.MaxFeePerGas != nil || params.MaxPriorityFeePerGas != nil {
				return nil, fmt.Errorf("EIP-1559 fees (maxFeePerGas/maxPriorityFeePerGas) are not supported for legacy/eip2930 transactions")
			}

			if params.GasPrice == nil {
				fees, feeErr := public.EstimateFeesPerGas(ctx, client, public.EstimateFeesPerGasParameters{
					Type: public.FeeValuesTypeLegacy,
				})
				if feeErr != nil {
					return nil, fmt.Errorf("failed to estimate gas price: %w", feeErr)
				}
				params.GasPrice = fees.GasPrice
			}
		}
	}

	// ---------- Fill gas ----------
	if containsParam(parameters, "gas") && params.Gas == nil {
		estimateParams := public.EstimateGasParameters{
			To:                   toCommonAddressPtr(params.To),
			Data:                 hexToBytes(params.Data),
			Value:                params.Value,
			GasPrice:             params.GasPrice,
			MaxFeePerGas:         params.MaxFeePerGas,
			MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
			MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		}
		if account != nil {
			addr := common.HexToAddress(account.Address().Hex())
			estimateParams.Account = &addr
		}
		if params.Nonce != nil {
			n := uint64(*params.Nonce)
			estimateParams.Nonce = &n
		}

		gas, gasErr := public.EstimateGas(ctx, client, estimateParams)
		if gasErr != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", gasErr)
		}
		params.Gas = new(big.Int).SetUint64(gas)
	}

	// Validate the final request
	if err := transaction.AssertRequest(transaction.AssertRequestParams{
		Account:              params.From,
		To:                   params.To,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
	}); err != nil {
		return nil, err
	}

	return &params, nil
}

// inferTransactionType attempts to infer the transaction type from the request fields.
// This mirrors viem's getTransactionType logic.
func inferTransactionType(params *PrepareTransactionRequestParameters) formatters.TransactionType {
	if len(params.AuthorizationList) > 0 {
		return formatters.TransactionTypeEIP7702
	}
	if len(params.BlobVersionedHashes) > 0 || len(params.Blobs) > 0 || params.MaxFeePerBlobGas != nil {
		return formatters.TransactionTypeEIP4844
	}
	if params.MaxFeePerGas != nil || params.MaxPriorityFeePerGas != nil {
		return formatters.TransactionTypeEIP1559
	}
	if len(params.AccessList) > 0 {
		return formatters.TransactionTypeEIP2930
	}
	if params.GasPrice != nil {
		return formatters.TransactionTypeLegacy
	}
	return "" // Cannot infer
}

// containsParam checks if a parameter list contains a given parameter.
func containsParam(params []string, param string) bool {
	for _, p := range params {
		if p == param {
			return true
		}
	}
	return false
}

// toCommonAddressPtr converts a hex address string to *common.Address, or nil if empty.
func toCommonAddressPtr(addr string) *common.Address {
	if addr == "" {
		return nil
	}
	a := common.HexToAddress(addr)
	return &a
}

// hexToBytes converts a hex string to bytes, returning nil for empty strings.
func hexToBytes(hex string) []byte {
	if hex == "" || hex == "0x" {
		return nil
	}
	return common.FromHex(hex)
}
