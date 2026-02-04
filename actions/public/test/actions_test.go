package public_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
)

// mockClient implements the public.Client interface for testing.
type mockClient struct {
	transport       transport.Transport
	chain           *chain.Chain
	cacheTime       time.Duration
	blockTag        types.BlockTag
	batch           *types.BatchOptions
	ccipRead        *types.CCIPReadOptions
	uid             string
	requestRecorder func(method string, params []any)
}

func (c *mockClient) Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error) {
	if c.requestRecorder != nil {
		c.requestRecorder(method, params)
	}
	return c.transport.Request(ctx, transport.RPCRequest{Method: method, Params: params})
}

func (c *mockClient) Chain() *chain.Chain {
	return c.chain
}

func (c *mockClient) CacheTime() time.Duration {
	if c.cacheTime == 0 {
		return 4 * time.Second
	}
	return c.cacheTime
}

func (c *mockClient) ExperimentalBlockTag() types.BlockTag {
	return c.blockTag
}

func (c *mockClient) Batch() *types.BatchOptions {
	return c.batch
}

func (c *mockClient) CCIPRead() *types.CCIPReadOptions {
	return c.ccipRead
}

func (c *mockClient) UID() string {
	if c.uid == "" {
		return "test-mock-client"
	}
	return c.uid
}

// createTestServer creates a test HTTP server that responds to JSON-RPC requests.
func createTestServer(t *testing.T, handler func(method string, params []any) any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			JSONRPC string `json:"jsonrpc"`
			ID      any    `json:"id"`
			Method  string `json:"method"`
			Params  []any  `json:"params"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		result := handler(req.Method, req.Params)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  result,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

// createMockClient creates a mock client for testing.
func createMockClient(t *testing.T, serverURL string) *mockClient {
	tr, err := transport.HTTP(serverURL)(transport.TransportParams{})
	require.NoError(t, err)

	return &mockClient{
		transport: tr,
	}
}

// ============================================================================
// Call Tests
// ============================================================================

func TestCall_Basic(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			// Return 42 encoded as uint256
			return "0x000000000000000000000000000000000000000000000000000000000000002a"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	result, err := public.Call(ctx, client, public.CallParameters{
		To: &to,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Data)
	assert.Equal(t, 32, len(result.Data))
}

func TestCall_WithData(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			return "0xabcdef"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	data := common.Hex2Bytes("a9059cbb") // transfer selector

	result, err := public.Call(ctx, client, public.CallParameters{
		To:   &to,
		Data: data,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Data)
}

func TestCall_EmptyResult(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			return "0x"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	result, err := public.Call(ctx, client, public.CallParameters{
		To: &to,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.Data)
}

func TestCall_WithBlockNumber(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	blockNum := uint64(12345)

	_, err := public.Call(ctx, client, public.CallParameters{
		To:          &to,
		BlockNumber: &blockNum,
	})

	require.NoError(t, err)
	require.Len(t, capturedParams, 2)
	assert.Equal(t, "0x3039", capturedParams[1]) // 12345 in hex
}

func TestCall_WithBlockTag(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	_, err := public.Call(ctx, client, public.CallParameters{
		To:       &to,
		BlockTag: public.BlockTagPending,
	})

	require.NoError(t, err)
	require.Len(t, capturedParams, 2)
	assert.Equal(t, "pending", capturedParams[1])
}

func TestCall_WithStateOverride(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0000000000000000000000000000000000000000000000000000000000000001"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	overrideAddr := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// Create state override with modified balance
	balance := big.NewInt(1000000000000000000) // 1 ETH
	stateOverride := types.StateOverride{
		overrideAddr: types.StateOverrideAccount{
			Balance: balance,
		},
	}

	result, err := public.Call(ctx, client, public.CallParameters{
		To:            &to,
		StateOverride: stateOverride,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	// Verify state override was passed (should be 3 params: request, block, stateOverride)
	require.GreaterOrEqual(t, len(capturedParams), 3)
}

func TestCall_WithBlockOverride(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0000000000000000000000000000000000000000000000000000000000000001"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create block override with modified gas limit
	gasLimit := uint64(30000000)
	blockOverrides := &types.BlockOverrides{
		GasLimit: &gasLimit,
	}

	result, err := public.Call(ctx, client, public.CallParameters{
		To:             &to,
		BlockOverrides: blockOverrides,
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	// Verify block override was passed (should be 4 params: request, block, stateOverride, blockOverride)
	require.GreaterOrEqual(t, len(capturedParams), 3)
}

func TestCall_InvalidParams_CodeAndFactory(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	factory := common.HexToAddress("0xfactory0000000000000000000000000000000")

	// Providing both code and factory should fail
	_, err := public.Call(ctx, client, public.CallParameters{
		Code:        []byte{0x60, 0x80, 0x60, 0x40, 0x52},
		Factory:     &factory,
		FactoryData: []byte{0x12, 0x34},
		Data:        []byte{0xab, 0xcd},
	})

	require.Error(t, err)
	_, ok := err.(*public.InvalidCallParamsError)
	assert.True(t, ok, "expected InvalidCallParamsError")
	assert.Contains(t, err.Error(), "cannot provide both 'code' and 'factory'")
}

func TestCall_InvalidParams_CodeAndTo(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Providing both code and to should fail
	_, err := public.Call(ctx, client, public.CallParameters{
		Code: []byte{0x60, 0x80, 0x60, 0x40, 0x52},
		To:   &to,
		Data: []byte{0xab, 0xcd},
	})

	require.Error(t, err)
	_, ok := err.(*public.InvalidCallParamsError)
	assert.True(t, ok, "expected InvalidCallParamsError")
	assert.Contains(t, err.Error(), "cannot provide both 'code' and 'to'")
}

func TestCall_WithValue(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	value := big.NewInt(1000000000000000000) // 1 ETH

	_, err := public.Call(ctx, client, public.CallParameters{
		To:    &to,
		Value: value,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(capturedParams), 1)
	reqMap, ok := capturedParams[0].(map[string]any)
	if ok {
		assert.Equal(t, "0xde0b6b3a7640000", reqMap["value"])
	}
}

func TestCall_WithGasAndGasPrice(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	gas := uint64(100000)
	gasPrice := big.NewInt(20000000000) // 20 gwei

	_, err := public.Call(ctx, client, public.CallParameters{
		To:       &to,
		Gas:      &gas,
		GasPrice: gasPrice,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(capturedParams), 1)
	reqMap, ok := capturedParams[0].(map[string]any)
	if ok {
		assert.Equal(t, "0x186a0", reqMap["gas"])
		assert.Equal(t, "0x4a817c800", reqMap["gasPrice"])
	}
}

func TestCall_WithEIP1559Fees(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	maxFeePerGas := big.NewInt(50000000000)     // 50 gwei
	maxPriorityFeePerGas := big.NewInt(2000000000) // 2 gwei

	_, err := public.Call(ctx, client, public.CallParameters{
		To:                   &to,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(capturedParams), 1)
	reqMap, ok := capturedParams[0].(map[string]any)
	if ok {
		assert.Equal(t, "0xba43b7400", reqMap["maxFeePerGas"])
		assert.Equal(t, "0x77359400", reqMap["maxPriorityFeePerGas"])
	}
}

func TestCall_WithAccount(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	from := common.HexToAddress("0xf00d000000000000000000000000000000000001")
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	_, err := public.Call(ctx, client, public.CallParameters{
		Account: &from,
		To:      &to,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(capturedParams), 1)
	reqMap, ok := capturedParams[0].(map[string]any)
	if ok {
		assert.Equal(t, "0xf00D000000000000000000000000000000000001", reqMap["from"])
	}
}

func TestCall_WithAccessList(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			capturedParams = params
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	accessList := types.AccessList{
		{
			Address: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa00"),
			StorageKeys: []common.Hash{
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
			},
		},
	}

	_, err := public.Call(ctx, client, public.CallParameters{
		To:         &to,
		AccessList: accessList,
	})

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(capturedParams), 1)
}

func TestCall_ErrorWrapping(t *testing.T) {
	// Test that errors are properly wrapped in CallExecutionError
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"error": map[string]any{
				"code":    3,
				"message": "execution reverted",
				"data":    "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b5465737420726576657274000000000000000000000000000000000000000000",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	_, err := public.Call(ctx, client, public.CallParameters{
		To: &to,
	})

	require.Error(t, err)
	// Should be wrapped in CallExecutionError
	_, ok := err.(*public.CallExecutionError)
	assert.True(t, ok, "expected CallExecutionError, got %T", err)
}

// ============================================================================
// GetBalance Tests
// ============================================================================

func TestGetBalance_Basic(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBalance" {
			// Return 1 ETH in wei
			return "0xde0b6b3a7640000"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	balance, err := public.GetBalance(ctx, client, public.GetBalanceParameters{
		Address: addr,
	})

	require.NoError(t, err)
	assert.NotNil(t, balance)

	expected := new(big.Int)
	expected.SetString("1000000000000000000", 10)
	assert.Equal(t, 0, balance.Cmp(expected))
}

func TestGetBalance_Zero(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBalance" {
			return "0x0"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	balance, err := public.GetBalance(ctx, client, public.GetBalanceParameters{
		Address: addr,
	})

	require.NoError(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, 0, balance.Cmp(big.NewInt(0)))
}

func TestGetBalance_WithBlockNumber(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBalance" {
			capturedParams = params
			return "0x1"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	blockNum := uint64(100)

	_, err := public.GetBalance(ctx, client, public.GetBalanceParameters{
		Address:     addr,
		BlockNumber: &blockNum,
	})

	require.NoError(t, err)
	require.Len(t, capturedParams, 2)
	assert.Equal(t, "0x64", capturedParams[1]) // 100 in hex
}

func TestGetBalance_WithBlockTag(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBalance" {
			capturedParams = params
			return "0x1"
		}
		return "0x0"
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	_, err := public.GetBalance(ctx, client, public.GetBalanceParameters{
		Address:  addr,
		BlockTag: public.BlockTagSafe,
	})

	require.NoError(t, err)
	require.Len(t, capturedParams, 2)
	assert.Equal(t, "safe", capturedParams[1])
}

// ============================================================================
// GetBlock Tests
// ============================================================================

func TestGetBlock_Latest(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBlockByNumber" {
			return map[string]any{
				"number":           "0x10",
				"hash":             "0x1234567890123456789012345678901234567890123456789012345678901234",
				"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce":            "0x0000000000000000",
				"sha3Uncles":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"stateRoot":        "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot":     "0x0000000000000000000000000000000000000000000000000000000000000000",
				"miner":            "0x0000000000000000000000000000000000000000",
				"difficulty":       "0x0",
				"totalDifficulty":  "0x0",
				"size":             "0x100",
				"gasLimit":         "0x1c9c380",
				"gasUsed":          "0x0",
				"timestamp":        "0x60000000",
				"transactions":     []string{},
				"uncles":           []string{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{})

	require.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, uint64(16), block.Number)
}

func TestGetBlock_ByNumber(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBlockByNumber" {
			capturedParams = params
			return map[string]any{
				"number":           "0x64",
				"hash":             "0x1234567890123456789012345678901234567890123456789012345678901234",
				"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce":            "0x0000000000000000",
				"sha3Uncles":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"stateRoot":        "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot":     "0x0000000000000000000000000000000000000000000000000000000000000000",
				"miner":            "0x0000000000000000000000000000000000000000",
				"difficulty":       "0x0",
				"totalDifficulty":  "0x0",
				"size":             "0x100",
				"gasLimit":         "0x1c9c380",
				"gasUsed":          "0x0",
				"timestamp":        "0x60000000",
				"transactions":     []string{},
				"uncles":           []string{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	blockNum := uint64(100)
	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{
		BlockNumber: &blockNum,
	})

	require.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, uint64(100), block.Number)
	require.Len(t, capturedParams, 2)
	assert.Equal(t, "0x64", capturedParams[0])
}

func TestGetBlock_ByHash(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		if method == "eth_getBlockByHash" {
			return map[string]any{
				"number":           "0x10",
				"hash":             "0x1234567890123456789012345678901234567890123456789012345678901234",
				"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"nonce":            "0x0000000000000000",
				"sha3Uncles":       "0x0000000000000000000000000000000000000000000000000000000000000000",
				"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
				"stateRoot":        "0x0000000000000000000000000000000000000000000000000000000000000000",
				"receiptsRoot":     "0x0000000000000000000000000000000000000000000000000000000000000000",
				"miner":            "0x0000000000000000000000000000000000000000",
				"difficulty":       "0x0",
				"totalDifficulty":  "0x0",
				"size":             "0x100",
				"gasLimit":         "0x1c9c380",
				"gasUsed":          "0x0",
				"timestamp":        "0x60000000",
				"transactions":     []string{},
				"uncles":           []string{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	hash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")
	block, err := public.GetBlock(ctx, client, public.GetBlockParameters{
		BlockHash: &hash,
	})

	require.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, "eth_getBlockByHash", capturedMethod)
}

func TestGetBlock_NotFound(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	blockNum := uint64(999999999)
	_, err := public.GetBlock(ctx, client, public.GetBlockParameters{
		BlockNumber: &blockNum,
	})

	require.Error(t, err)
	_, ok := err.(*public.BlockNotFoundError)
	assert.True(t, ok, "expected BlockNotFoundError")
}

// ============================================================================
// GetTransaction Tests
// ============================================================================

func TestGetTransaction_ByHash(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getTransactionByHash" {
			return map[string]any{
				"blockHash":        "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":      "0x10",
				"from":             "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"gas":              "0x5208",
				"gasPrice":         "0x3b9aca00",
				"hash":             "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				"input":            "0x",
				"nonce":            "0x1",
				"to":               "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				"transactionIndex": "0x0",
				"value":            "0xde0b6b3a7640000",
				"type":             "0x0",
				"v":                "0x1c",
				"r":                "0x1234",
				"s":                "0x5678",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	hash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
		Hash: &hash,
	})

	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, hash, tx.Hash)
	assert.Equal(t, uint64(16), *tx.BlockNumber)
	assert.Equal(t, common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), tx.From)
}

func TestGetTransaction_ByBlockHashAndIndex(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		if method == "eth_getTransactionByBlockHashAndIndex" {
			return map[string]any{
				"blockHash":        "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":      "0x10",
				"from":             "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"gas":              "0x5208",
				"gasPrice":         "0x3b9aca00",
				"hash":             "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				"input":            "0x",
				"nonce":            "0x1",
				"to":               "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				"transactionIndex": "0x0",
				"value":            "0xde0b6b3a7640000",
				"type":             "0x0",
				"v":                "0x1c",
				"r":                "0x1234",
				"s":                "0x5678",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	blockHash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")
	index := 0
	tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
		BlockHash: &blockHash,
		Index:     &index,
	})

	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "eth_getTransactionByBlockHashAndIndex", capturedMethod)
}

func TestGetTransaction_ByBlockNumberAndIndex(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		if method == "eth_getTransactionByBlockNumberAndIndex" {
			return map[string]any{
				"blockHash":        "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":      "0x64",
				"from":             "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				"gas":              "0x5208",
				"gasPrice":         "0x3b9aca00",
				"hash":             "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				"input":            "0x",
				"nonce":            "0x1",
				"to":               "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				"transactionIndex": "0x0",
				"value":            "0xde0b6b3a7640000",
				"type":             "0x0",
				"v":                "0x1c",
				"r":                "0x1234",
				"s":                "0x5678",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	blockNum := uint64(100)
	index := 0
	tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
		BlockNumber: &blockNum,
		Index:       &index,
	})

	require.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, "eth_getTransactionByBlockNumberAndIndex", capturedMethod)
}

func TestGetTransaction_NotFound(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	hash := common.HexToHash("0xdeadbeef1234567890abcdef1234567890abcdef1234567890abcdef12345678")
	_, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
		Hash: &hash,
	})

	require.Error(t, err)
	_, ok := err.(*public.TransactionNotFoundError)
	assert.True(t, ok, "expected TransactionNotFoundError")
}

func TestGetTransaction_InvalidParams(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	// No parameters provided
	_, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameters")
}
