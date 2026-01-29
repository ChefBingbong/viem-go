package erc20_test

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// MockTransport is a mock transport for testing contract calls.
type MockTransport struct {
	responses map[string]json.RawMessage
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]json.RawMessage),
	}
}

func (m *MockTransport) SetResponse(method string, result any) {
	data, _ := json.Marshal(result)
	m.responses[method] = data
}

func (m *MockTransport) Call(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	if resp, ok := m.responses[method]; ok {
		return resp, nil
	}
	// Return a default response
	return json.RawMessage(`"0x"`), nil
}

func (m *MockTransport) Close() error {
	return nil
}

var _ = Describe("ERC20 Contract Bindings", func() {
	var (
		mockTransport *MockTransport
		rpcClient     *client.Client
		token         *erc20.ERC20
		tokenAddress  common.Address
		ctx           context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		tokenAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") // USDC address

		// Create mock transport and client
		mockTransport = NewMockTransport()

		// Mock chain ID response
		mockTransport.SetResponse("eth_chainId", "0x1")

		var err error
		rpcClient, err = client.NewClientWithTransport(mockTransport)
		Expect(err).ToNot(HaveOccurred())

		// Create ERC20 binding
		token, err = erc20.New(tokenAddress, rpcClient)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when reading token metadata", func() {
		It("should return the token name", func() {
			// Mock the eth_call response for name()
			// "USD Coin" encoded as ABI string
			// Offset (32) + Length (8) + "USD Coin" padded
			nameEncoded := "0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" + // offset
				"0000000000000000000000000000000000000000000000000000000000000008" + // length = 8
				"55534420436f696e000000000000000000000000000000000000000000000000" // "USD Coin"
			mockTransport.SetResponse("eth_call", nameEncoded)

			name, err := token.Name(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("USD Coin"))
		})

		It("should return the token symbol", func() {
			// "USDC" encoded as ABI string
			symbolEncoded := "0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" + // offset
				"0000000000000000000000000000000000000000000000000000000000000004" + // length = 4
				"5553444300000000000000000000000000000000000000000000000000000000" // "USDC"
			mockTransport.SetResponse("eth_call", symbolEncoded)

			symbol, err := token.Symbol(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(symbol).To(Equal("USDC"))
		})

		It("should return the token decimals", func() {
			// 6 decimals encoded as uint8
			decimalsEncoded := "0x0000000000000000000000000000000000000000000000000000000000000006"
			mockTransport.SetResponse("eth_call", decimalsEncoded)

			decimals, err := token.Decimals(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(decimals).To(Equal(uint8(6)))
		})
	})

	Context("when reading balances", func() {
		It("should return the balance of an address", func() {
			// 1000000 (1 USDC with 6 decimals) encoded as uint256
			balanceEncoded := "0x00000000000000000000000000000000000000000000000000000000000f4240"
			mockTransport.SetResponse("eth_call", balanceEncoded)

			owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
			balance, err := token.BalanceOf(ctx, owner)
			Expect(err).ToNot(HaveOccurred())
			Expect(balance.Cmp(big.NewInt(1000000))).To(Equal(0))
		})

		It("should return the total supply", func() {
			// 1 billion USDC (1000000000 * 10^6)
			totalSupply := new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(1000000))
			supplyEncoded := hexutil.EncodeBig(totalSupply)
			// Pad to 32 bytes
			supplyEncoded = "0x" + padLeft(supplyEncoded[2:], 64)
			mockTransport.SetResponse("eth_call", supplyEncoded)

			supply, err := token.TotalSupply(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(supply.Cmp(totalSupply)).To(Equal(0))
		})
	})

	Context("when reading allowances", func() {
		It("should return the allowance for a spender", func() {
			// 500000 (0.5 USDC) allowance
			allowanceEncoded := "0x000000000000000000000000000000000000000000000000000000000007a120"
			mockTransport.SetResponse("eth_call", allowanceEncoded)

			owner := common.HexToAddress("0x1111111111111111111111111111111111111111")
			spender := common.HexToAddress("0x2222222222222222222222222222222222222222")

			allowance, err := token.Allowance(ctx, owner, spender)
			Expect(err).ToNot(HaveOccurred())
			Expect(allowance.Cmp(big.NewInt(500000))).To(Equal(0))
		})
	})

	Context("when accessing contract properties", func() {
		It("should return the contract address", func() {
			Expect(token.Address()).To(Equal(tokenAddress))
		})

		It("should return the underlying contract", func() {
			Expect(token.Contract()).ToNot(BeNil())
			Expect(token.Contract().Address()).To(Equal(tokenAddress))
		})
	})
})

// padLeft pads a hex string with zeros to the specified length.
func padLeft(s string, length int) string {
	for len(s) < length {
		s = "0" + s
	}
	return s
}
