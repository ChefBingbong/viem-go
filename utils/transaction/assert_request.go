package transaction

import (
	"fmt"
	"math/big"
	"regexp"
)

var addressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// AssertRequestParams contains parameters for transaction request validation.
type AssertRequestParams struct {
	Account              string
	To                   string
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

// AssertRequest validates a transaction request.
// Checks address validity and fee constraints.
//
// Example:
//
//	err := AssertRequest(AssertRequestParams{
//		To: "0x1234...",
//		MaxFeePerGas: big.NewInt(1000000000),
//	})
func AssertRequest(params AssertRequestParams) error {
	// Validate account address
	if params.Account != "" && !isValidAddress(params.Account) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, params.Account)
	}

	// Validate to address
	if params.To != "" && !isValidAddress(params.To) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, params.To)
	}

	// Check maxFeePerGas doesn't exceed max uint256
	if params.MaxFeePerGas != nil && params.MaxFeePerGas.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("%w: maxFeePerGas exceeds maximum value", ErrFeeCapTooHigh)
	}

	// Check tip doesn't exceed fee cap
	if params.MaxPriorityFeePerGas != nil && params.MaxFeePerGas != nil {
		if params.MaxPriorityFeePerGas.Cmp(params.MaxFeePerGas) > 0 {
			return fmt.Errorf("%w: maxPriorityFeePerGas (%s) > maxFeePerGas (%s)",
				ErrTipAboveFeeCap, params.MaxPriorityFeePerGas.String(), params.MaxFeePerGas.String())
		}
	}

	return nil
}

// isValidAddress checks if a string is a valid Ethereum address format.
func isValidAddress(address string) bool {
	return addressRegex.MatchString(address)
}
