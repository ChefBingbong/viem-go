package wallet

import (
	"context"

	"github.com/ChefBingbong/viem-go/types"
)

// SignAuthorizationParameters contains the parameters for the SignAuthorization action.
// This mirrors viem's SignAuthorizationParameters type, which extends PrepareAuthorizationParameters.
type SignAuthorizationParameters = PrepareAuthorizationParameters

// SignAuthorizationReturnType is the return type for the SignAuthorization action.
// This mirrors viem's SignAuthorizationReturnType (a signed EIP-7702 authorization).
type SignAuthorizationReturnType = *types.SignedAuthorization

// SignAuthorization signs an EIP-7702 Authorization object.
//
// This action first prepares the authorization (filling in chainId and nonce if needed),
// then signs it using the account's signAuthorization method.
//
// Note: This action requires a local account that implements AuthorizationSignableAccount.
// JSON-RPC accounts are not supported for this action.
//
// With the calculated signature, you can:
// - use verifyAuthorization to verify the signed Authorization object
// - use recoverAuthorizationAddress to recover the signing address
//
// This is equivalent to viem's `signAuthorization` action.
//
// Example:
//
//	signed, err := wallet.SignAuthorization(ctx, client, wallet.SignAuthorizationParameters{
//	    Account:         myLocalAccount,
//	    ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
//	})
//
// Example with account hoisting:
//
//	// client created with: CreateWalletClient({ account: privateKeyToAccount('0x...'), ... })
//	signed, err := wallet.SignAuthorization(ctx, client, wallet.SignAuthorizationParameters{
//	    ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
//	})
func SignAuthorization(ctx context.Context, client Client, params SignAuthorizationParameters) (SignAuthorizationReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return nil, &AccountNotFoundError{DocsPath: "/docs/eip7702/signAuthorization"}
	}

	// Verify the account supports signing authorizations (must be local)
	// This mirrors viem's `if (!account.signAuthorization) throw new AccountTypeNotSupportedError`
	signable, ok := account.(AuthorizationSignableAccount)
	if !ok {
		return nil, &AccountTypeNotSupportedError{
			DocsPath: "/docs/eip7702/signAuthorization",
			MetaMessages: []string{
				"The `signAuthorization` Action does not support JSON-RPC Accounts.",
			},
		}
	}

	// Prepare the authorization (fill chainId, nonce if needed)
	// This mirrors viem's: const authorization = await prepareAuthorization(client, parameters)
	auth, err := PrepareAuthorization(ctx, client, params)
	if err != nil {
		return nil, err
	}

	// Sign the prepared authorization
	// This mirrors viem's: return account.signAuthorization(authorization)
	signed, err := signable.SignAuthorization(types.AuthorizationRequest{
		Address: auth.Address,
		ChainId: auth.ChainId,
		Nonce:   auth.Nonce,
	})
	if err != nil {
		return nil, err
	}

	return signed, nil
}
