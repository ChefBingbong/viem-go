package wallet

import (
	"context"
	"fmt"
	"time"

	viemchain "github.com/ChefBingbong/viem-go/chain"
)

// SendCallsSyncParameters contains the parameters for the SendCallsSync action.
// This mirrors viem's SendCallsSyncParameters type which extends SendCallsParameters
// with polling/wait fields.
type SendCallsSyncParameters struct {
	SendCallsParameters

	// PollingInterval is the polling interval to poll for the calls status.
	// Defaults to client.PollingInterval().
	PollingInterval time.Duration

	// Status is a function that determines whether to stop polling.
	// Default: statusCode == 200 || statusCode >= 300.
	Status func(result *GetCallsStatusReturnType) bool

	// ThrowOnFailure when true, returns an error if the call bundle fails.
	// Default: false.
	ThrowOnFailure bool

	// Timeout is the timeout (in ms) to wait for calls to be included in a block.
	// Default: max(chain.blockTime * 3, 5000)ms.
	Timeout *time.Duration
}

// SendCallsSyncReturnType is the return type for the SendCallsSync action.
type SendCallsSyncReturnType = *GetCallsStatusReturnType

// SendCallsSync requests the connected wallet to send a batch of calls, and waits
// for the calls to be included in a block.
//
// This is equivalent to viem's `sendCallsSync` action.
//
// JSON-RPC Methods: wallet_sendCalls + wallet_getCallsStatus (EIP-5792)
//
// Example:
//
//	status, err := wallet.SendCallsSync(ctx, client, wallet.SendCallsSyncParameters{
//	    SendCallsParameters: wallet.SendCallsParameters{
//	        Calls: []wallet.Call{
//	            {Data: "0xdeadbeef", To: "0x70997970c51812dc3a010c7d01b50e0d17dc79c8"},
//	            {To: "0x70997970c51812dc3a010c7d01b50e0d17dc79c8", Value: big.NewInt(69420)},
//	        },
//	    },
//	})
func SendCallsSync(ctx context.Context, client Client, params SendCallsSyncParameters) (SendCallsSyncReturnType, error) {
	// Resolve timeout (mirrors viem's: parameters.timeout ?? Math.max((chain?.blockTime ?? 0) * 3, 5_000))
	timeout := resolveSendCallsTimeout(params.Timeout, client.Chain())

	// Send the calls
	result, err := SendCalls(ctx, client, params.SendCallsParameters)
	if err != nil {
		return nil, fmt.Errorf("sendCalls failed: %w", err)
	}

	// Wait for the calls to be included in a block
	status, waitErr := WaitForCallsStatus(ctx, client, WaitForCallsStatusParameters{
		ID:              result.ID,
		PollingInterval: params.PollingInterval,
		Status:          params.Status,
		ThrowOnFailure:  params.ThrowOnFailure,
		Timeout:         timeout,
	})
	if waitErr != nil {
		return nil, waitErr
	}

	return status, nil
}

// resolveSendCallsTimeout resolves the timeout for sendCallsSync.
// Mirrors viem's: parameters.timeout ?? Math.max((chain?.blockTime ?? 0) * 3, 5_000)
func resolveSendCallsTimeout(paramTimeout *time.Duration, ch *viemchain.Chain) time.Duration {
	if paramTimeout != nil {
		return *paramTimeout
	}

	var blockTimeMs int64
	if ch != nil && ch.BlockTime != nil {
		blockTimeMs = *ch.BlockTime
	}

	timeoutMs := blockTimeMs * 3
	if timeoutMs < 5000 {
		timeoutMs = 5000
	}

	return time.Duration(timeoutMs) * time.Millisecond
}
