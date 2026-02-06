package wallet

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/utils/authorization"
)

// PrepareAuthorizationParameters contains the parameters for the PrepareAuthorization action.
// This mirrors viem's PrepareAuthorizationParameters type.
type PrepareAuthorizationParameters struct {
	// Account is the account to prepare authorization for. If nil, uses the client's account.
	Account Account

	// ContractAddress is the contract address being authorized.
	// Either ContractAddress or Address must be provided.
	ContractAddress string

	// Address is an alias for ContractAddress (for compatibility).
	Address string

	// ChainID optionally specifies the chain ID for the authorization.
	// If nil, it will be fetched from the client's chain or via eth_chainId.
	ChainID *int

	// Nonce optionally specifies the nonce for the authorization.
	// If nil, it will be fetched via eth_getTransactionCount with "pending" block tag.
	Nonce *int

	// Executor specifies who will execute the EIP-7702 transaction.
	// - nil: assumes another account will execute
	// - "self": the signing account will execute (nonce += 1)
	// - Account: a specific account will execute (nonce += 1 if same as signing account)
	Executor any
}

// PrepareAuthorizationReturnType is the return type for the PrepareAuthorization action.
// This mirrors viem's Authorization type.
type PrepareAuthorizationReturnType = authorization.AuthorizationRequest

// PrepareAuthorization prepares an EIP-7702 Authorization object for signing.
// This Action will fill the required fields of the Authorization object if they are
// not provided (e.g. `nonce` and `chainId`).
//
// With the prepared Authorization object, you can use SignAuthorization to sign
// over the Authorization object.
//
// This is equivalent to viem's `prepareAuthorization` action.
//
// Example:
//
//	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
//	    ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
//	})
//
// Example with account hoisting:
//
//	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
//	    Account:         myAccount,
//	    ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
//	})
func PrepareAuthorization(ctx context.Context, client Client, params PrepareAuthorizationParameters) (PrepareAuthorizationReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return PrepareAuthorizationReturnType{}, &AccountNotFoundError{DocsPath: "/docs/eip7702/prepareAuthorization"}
	}

	// Resolve executor (mirrors viem's executor parsing)
	var executorAddr *string
	isSelfExecutor := false
	if params.Executor != nil {
		switch v := params.Executor.(type) {
		case string:
			if v == "self" {
				isSelfExecutor = true
			} else {
				executorAddr = &v
			}
		case Account:
			addr := v.Address().Hex()
			executorAddr = &addr
		}
	}

	// Build the authorization
	addr := params.ContractAddress
	if addr == "" {
		addr = params.Address
	}

	auth := authorization.AuthorizationRequest{
		Address: addr,
	}

	// Resolve chain ID (mirrors viem's chainId resolution)
	if params.ChainID != nil {
		auth.ChainId = *params.ChainID
	} else {
		// Try client chain first, then RPC
		if ch := client.Chain(); ch != nil {
			auth.ChainId = int(ch.ID)
		} else {
			chainID, err := public.GetChainID(ctx, client)
			if err != nil {
				return PrepareAuthorizationReturnType{}, fmt.Errorf("failed to get chain ID: %w", err)
			}
			auth.ChainId = int(chainID)
		}
	}

	// Resolve nonce (mirrors viem's nonce resolution)
	if params.Nonce != nil {
		auth.Nonce = *params.Nonce
	} else {
		// Fetch the pending nonce for the account
		nonce, err := public.GetTransactionCount(ctx, client, public.GetTransactionCountParameters{
			Address:  common.HexToAddress(account.Address().Hex()),
			BlockTag: "pending",
		})
		if err != nil {
			return PrepareAuthorizationReturnType{}, fmt.Errorf("failed to get transaction count: %w", err)
		}
		auth.Nonce = int(nonce)

		// If executor is self or executor address matches signing account, increment nonce
		// This mirrors viem's: if (executor === 'self' || (executor?.address && isAddressEqual(...))) authorization.nonce += 1
		if isSelfExecutor {
			auth.Nonce++
		} else if executorAddr != nil {
			if strings.EqualFold(*executorAddr, account.Address().Hex()) {
				auth.Nonce++
			}
		}
	}

	return auth, nil
}
