package test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utils "github.com/ChefBingbong/viem-go/accounts/utils"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// Test private key (Anvil account 0)
const testPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

var _ = Describe("Accounts Utils", func() {
	Describe("Sign", func() {
		It("should sign a hash and return a signature object", func() {
			hash := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
			result, err := utils.Sign(utils.SignParameters{
				Hash:       hash,
				PrivateKey: testPrivateKey,
				To:         utils.SignReturnFormatObject,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Signature).NotTo(BeNil())
			Expect(result.Signature.R).To(HavePrefix("0x"))
			Expect(result.Signature.S).To(HavePrefix("0x"))
			Expect(result.Signature.YParity).To(BeNumerically(">=", 0))
			Expect(result.Signature.YParity).To(BeNumerically("<=", 1))
		})

		It("should sign a hash and return hex", func() {
			hash := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
			result, err := utils.Sign(utils.SignParameters{
				Hash:       hash,
				PrivateKey: testPrivateKey,
				To:         utils.SignReturnFormatHex,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Hex).To(HavePrefix("0x"))
			Expect(len(result.Hex)).To(Equal(132)) // 0x + 64 + 64 + 2 = 132
		})

		It("should fail with invalid private key", func() {
			hash := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
			_, err := utils.Sign(utils.SignParameters{
				Hash:       hash,
				PrivateKey: "0xinvalid",
			})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SignMessage", func() {
		It("should sign a message", func() {
			sig, err := utils.SignMessage(utils.SignMessageParameters{
				Message:    signature.NewSignableMessage("hello world"),
				PrivateKey: testPrivateKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sig).To(HavePrefix("0x"))
			Expect(len(sig)).To(Equal(132))
		})

		It("should return consistent signatures", func() {
			sig1, err := utils.SignMessage(utils.SignMessageParameters{
				Message:    signature.NewSignableMessage("hello"),
				PrivateKey: testPrivateKey,
			})
			Expect(err).NotTo(HaveOccurred())

			sig2, err := utils.SignMessage(utils.SignMessageParameters{
				Message:    signature.NewSignableMessage("hello"),
				PrivateKey: testPrivateKey,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(sig1).To(Equal(sig2))
		})
	})

	Describe("SignTypedData", func() {
		It("should sign EIP-712 typed data", func() {
			sig, err := utils.SignTypedData(utils.SignTypedDataParameters{
				Domain: signature.TypedDataDomain{
					Name:              "Ether Mail",
					Version:           "1",
					ChainId:           big.NewInt(1),
					VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
				},
				Types: map[string][]signature.TypedDataField{
					"Person": {
						{Name: "name", Type: "string"},
						{Name: "wallet", Type: "address"},
					},
					"Mail": {
						{Name: "from", Type: "Person"},
						{Name: "to", Type: "Person"},
						{Name: "contents", Type: "string"},
					},
				},
				PrimaryType: "Mail",
				Message: map[string]any{
					"from": map[string]any{
						"name":   "Cow",
						"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
					},
					"to": map[string]any{
						"name":   "Bob",
						"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
					},
					"contents": "Hello, Bob!",
				},
				PrivateKey: testPrivateKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sig).To(HavePrefix("0x"))
		})
	})

	Describe("SignTransaction", func() {
		It("should sign an EIP-1559 transaction", func() {
			tx := &transaction.Transaction{
				Type:                 transaction.TransactionTypeEIP1559,
				ChainId:              1,
				Nonce:                0,
				MaxPriorityFeePerGas: big.NewInt(1000000000),
				MaxFeePerGas:         big.NewInt(2000000000),
				Gas:                  big.NewInt(21000),
				To:                   "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
				Value:                big.NewInt(1000000000000000000),
			}

			signedTx, err := utils.SignTransaction(utils.SignTransactionParameters{
				PrivateKey:  testPrivateKey,
				Transaction: tx,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(signedTx).To(HavePrefix("0x02"))
		})

		It("should sign a legacy transaction", func() {
			tx := &transaction.Transaction{
				Type:     transaction.TransactionTypeLegacy,
				ChainId:  1,
				Nonce:    0,
				GasPrice: big.NewInt(20000000000),
				Gas:      big.NewInt(21000),
				To:       "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
				Value:    big.NewInt(1000000000000000000),
			}

			signedTx, err := utils.SignTransaction(utils.SignTransactionParameters{
				PrivateKey:  testPrivateKey,
				Transaction: tx,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(signedTx).To(HavePrefix("0x"))
		})
	})

	Describe("SignAuthorization", func() {
		It("should sign an EIP-7702 authorization", func() {
			result, err := utils.SignAuthorization(utils.SignAuthorizationParameters{
				Address:    "0x1234567890123456789012345678901234567890",
				ChainId:    1,
				Nonce:      0,
				PrivateKey: testPrivateKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.SignedAuthorization).NotTo(BeNil())
			Expect(result.SignedAuthorization.Address).To(Equal("0x1234567890123456789012345678901234567890"))
			Expect(result.SignedAuthorization.ChainId).To(Equal(1))
			Expect(result.SignedAuthorization.R).To(HavePrefix("0x"))
			Expect(result.SignedAuthorization.S).To(HavePrefix("0x"))
		})
	})

	Describe("PrivateKeyToAddress", func() {
		It("should convert private key to address", func() {
			addr, err := utils.PrivateKeyToAddress(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())
			// Compare case-insensitively since checksum format may differ
			Expect(addr).To(BeEquivalentTo(testAddress))
		})

		It("should fail with invalid private key", func() {
			_, err := utils.PrivateKeyToAddress("0xinvalid")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("PrivateKeyToPublicKey", func() {
		It("should derive public key from private key", func() {
			pubKey, err := utils.PrivateKeyToPublicKey(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(pubKey).To(HavePrefix("0x04"))
			Expect(len(pubKey)).To(Equal(132)) // 0x + 04 + 128 hex chars
		})
	})

	Describe("PublicKeyToAddress", func() {
		It("should convert public key to address", func() {
			pubKey, err := utils.PrivateKeyToPublicKey(testPrivateKey)
			Expect(err).NotTo(HaveOccurred())

			addr, err := utils.PublicKeyToAddress(pubKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(BeEquivalentTo(testAddress))
		})
	})

	Describe("ParseAccount", func() {
		It("should parse address string to json-rpc account", func() {
			account := utils.ParseAccount(testAddress)
			Expect(account.Address).To(Equal(testAddress))
			Expect(account.Type).To(Equal(utils.AccountTypeJSONRPC))
		})

		It("should return account as-is", func() {
			original := utils.Account{
				Address: testAddress,
				Type:    utils.AccountTypeLocal,
			}
			account := utils.ParseAccountFromAccount(original)
			Expect(account.Address).To(Equal(testAddress))
			Expect(account.Type).To(Equal(utils.AccountTypeLocal))
		})

		It("should handle generic input", func() {
			// String input
			account := utils.ParseAccountGeneric(testAddress)
			Expect(account.Type).To(Equal(utils.AccountTypeJSONRPC))

			// Account input
			original := utils.Account{Address: testAddress, Type: utils.AccountTypeLocal}
			account2 := utils.ParseAccountGeneric(original)
			Expect(account2.Type).To(Equal(utils.AccountTypeLocal))
		})
	})
})
