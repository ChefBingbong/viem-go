package wallet

import (
	"context"
	"fmt"
)

// ShowCallsStatusParameters contains the parameters for the ShowCallsStatus action.
// This mirrors viem's ShowCallsStatusParameters type.
type ShowCallsStatusParameters struct {
	// ID is the identifier of the call batch.
	ID string
}

// ShowCallsStatus requests the wallet to show information about a call batch
// that was sent via SendCalls.
//
// This is equivalent to viem's `showCallsStatus` action.
//
// JSON-RPC Method: wallet_showCallsStatus (EIP-5792)
//
// Example:
//
//	err := wallet.ShowCallsStatus(ctx, client, wallet.ShowCallsStatusParameters{
//	    ID: "0xdeadbeef",
//	})
func ShowCallsStatus(ctx context.Context, client Client, params ShowCallsStatusParameters) error {
	_, err := client.Request(ctx, "wallet_showCallsStatus", params.ID)
	if err != nil {
		return fmt.Errorf("wallet_showCallsStatus failed: %w", err)
	}
	return nil
}
