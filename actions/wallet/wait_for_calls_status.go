package wallet

import (
	"context"
	"fmt"
	"math"
	"time"
)

// WaitForCallsStatusParameters contains the parameters for the WaitForCallsStatus action.
// This mirrors viem's WaitForCallsStatusParameters type.
type WaitForCallsStatusParameters struct {
	// ID is the identifier of the call batch to wait for.
	ID string

	// PollingInterval is the polling frequency.
	// Defaults to the client's pollingInterval config.
	PollingInterval time.Duration

	// RetryCount is the number of times to retry if the call bundle request fails.
	// Default: 4 (exponential backoff).
	RetryCount *int

	// RetryDelay is a function that returns the delay between retries.
	// Default: exponential backoff: (1 << count) * 200ms.
	RetryDelay func(count int) time.Duration

	// Status is a function that determines whether to stop polling based on the status.
	// Default: statusCode == 200 || statusCode >= 300.
	Status func(result *GetCallsStatusReturnType) bool

	// ThrowOnFailure when true, returns an error if the call bundle fails.
	// Default: false.
	ThrowOnFailure bool

	// Timeout is the maximum time to wait before stopping polling.
	// Default: 60 seconds.
	Timeout time.Duration
}

// WaitForCallsStatusReturnType is the return type for the WaitForCallsStatus action.
type WaitForCallsStatusReturnType = *GetCallsStatusReturnType

// WaitForCallsStatusTimeoutError is returned when waiting for calls status times out.
type WaitForCallsStatusTimeoutError struct {
	ID string
}

func (e *WaitForCallsStatusTimeoutError) Error() string {
	return fmt.Sprintf("timed out while waiting for call bundle with id %q to be confirmed", e.ID)
}

// BundleFailedError is returned when a call bundle fails and ThrowOnFailure is true.
type BundleFailedError struct {
	Result *GetCallsStatusReturnType
}

func (e *BundleFailedError) Error() string {
	return fmt.Sprintf("call bundle failed with status %q (code: %d)", e.Result.Status, e.Result.StatusCode)
}

// WaitForCallsStatus waits for the status & receipts of a call bundle that was sent via SendCalls.
//
// This is equivalent to viem's `waitForCallsStatus` action.
//
// JSON-RPC Method: wallet_getCallsStatus (EIP-5792) — polled repeatedly.
//
// Example:
//
//	status, err := wallet.WaitForCallsStatus(ctx, client, wallet.WaitForCallsStatusParameters{
//	    ID: "0xdeadbeef",
//	})
func WaitForCallsStatus(ctx context.Context, client Client, params WaitForCallsStatusParameters) (WaitForCallsStatusReturnType, error) {
	// Resolve defaults
	pollingInterval := params.PollingInterval
	if pollingInterval == 0 {
		pollingInterval = client.PollingInterval()
	}

	retryCount := 4
	if params.RetryCount != nil {
		retryCount = *params.RetryCount
	}

	retryDelay := params.RetryDelay
	if retryDelay == nil {
		retryDelay = func(count int) time.Duration {
			return time.Duration(int(math.Pow(2, float64(count)))*200) * time.Millisecond
		}
	}

	statusCheck := params.Status
	if statusCheck == nil {
		statusCheck = func(result *GetCallsStatusReturnType) bool {
			return result.StatusCode == 200 || result.StatusCode >= 300
		}
	}

	timeout := params.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Poll loop (mirrors viem's poll + observe pattern)
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	// Emit on begin: check immediately
	result, err := getCallsStatusWithRetry(timeoutCtx, client, params.ID, retryCount, retryDelay)
	if err != nil {
		return nil, err
	}
	if params.ThrowOnFailure && result.Status == "failure" {
		return nil, &BundleFailedError{Result: result}
	}
	if statusCheck(result) {
		return result, nil
	}

	for {
		select {
		case <-timeoutCtx.Done():
			return nil, &WaitForCallsStatusTimeoutError{ID: params.ID}
		case <-ticker.C:
			result, err := getCallsStatusWithRetry(timeoutCtx, client, params.ID, retryCount, retryDelay)
			if err != nil {
				return nil, err
			}
			if params.ThrowOnFailure && result.Status == "failure" {
				return nil, &BundleFailedError{Result: result}
			}
			if statusCheck(result) {
				return result, nil
			}
		}
	}
}

// getCallsStatusWithRetry calls GetCallsStatus with retry logic.
// This mirrors viem's withRetry wrapper.
func getCallsStatusWithRetry(
	ctx context.Context,
	client Client,
	id string,
	retryCount int,
	retryDelay func(count int) time.Duration,
) (*GetCallsStatusReturnType, error) {
	var lastErr error
	for attempt := 0; attempt <= retryCount; attempt++ {
		result, err := GetCallsStatus(ctx, client, GetCallsStatusParameters{ID: id})
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempt < retryCount {
			delay := retryDelay(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}
	return nil, fmt.Errorf("getCallsStatus failed after %d retries: %w", retryCount+1, lastErr)
}
