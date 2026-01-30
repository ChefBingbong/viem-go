package test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils/transaction"
)

var _ = Describe("Transaction", func() {
	Describe("GetTransactionType", func() {
		It("should detect EIP-1559 transaction", func() {
			tx := &transaction.Transaction{
				MaxFeePerGas:         big.NewInt(1000000000),
				MaxPriorityFeePerGas: big.NewInt(100000000),
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP1559))
		})

		It("should detect legacy transaction", func() {
			tx := &transaction.Transaction{
				GasPrice: big.NewInt(1000000000),
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeLegacy))
		})

		It("should detect EIP-2930 transaction", func() {
			tx := &transaction.Transaction{
				GasPrice: big.NewInt(1000000000),
				AccessList: transaction.AccessList{
					{Address: "0x1234567890123456789012345678901234567890", StorageKeys: []string{}},
				},
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP2930))
		})

		It("should detect EIP-4844 transaction", func() {
			tx := &transaction.Transaction{
				MaxFeePerGas:        big.NewInt(1000000000),
				MaxFeePerBlobGas:    big.NewInt(100000),
				BlobVersionedHashes: []string{"0x01..."},
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP4844))
		})

		It("should detect EIP-7702 transaction", func() {
			tx := &transaction.Transaction{
				MaxFeePerGas: big.NewInt(1000000000),
				AuthorizationList: []transaction.SignedAuthorization{
					{
						Authorization: transaction.Authorization{
							Address: "0x1234567890123456789012345678901234567890",
							ChainId: 1,
							Nonce:   0,
						},
					},
				},
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP7702))
		})

		It("should use explicit type if set", func() {
			tx := &transaction.Transaction{
				Type:     transaction.TransactionTypeEIP1559,
				GasPrice: big.NewInt(1000000000), // Would normally be legacy
			}
			txType, err := transaction.GetTransactionType(tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP1559))
		})
	})

	Describe("GetSerializedTransactionType", func() {
		It("should detect EIP-1559 from serialized", func() {
			txType, err := transaction.GetSerializedTransactionType("0x02f8...")
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP1559))
		})

		It("should detect EIP-2930 from serialized", func() {
			txType, err := transaction.GetSerializedTransactionType("0x01f8...")
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP2930))
		})

		It("should detect EIP-4844 from serialized", func() {
			txType, err := transaction.GetSerializedTransactionType("0x03f8...")
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP4844))
		})

		It("should detect EIP-7702 from serialized", func() {
			txType, err := transaction.GetSerializedTransactionType("0x04f8...")
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeEIP7702))
		})

		It("should detect legacy from serialized", func() {
			// Legacy starts with 0xc0-0xff (RLP list prefix)
			txType, err := transaction.GetSerializedTransactionType("0xf8...")
			Expect(err).NotTo(HaveOccurred())
			Expect(txType).To(Equal(transaction.TransactionTypeLegacy))
		})
	})

	Describe("SerializeAccessList", func() {
		It("should serialize an empty access list", func() {
			result, err := transaction.SerializeAccessList(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should serialize access list with storage keys", func() {
			accessList := transaction.AccessList{
				{
					Address: "0x1234567890123456789012345678901234567890",
					StorageKeys: []string{
						"0x0000000000000000000000000000000000000000000000000000000000000001",
					},
				},
			}
			result, err := transaction.SerializeAccessList(accessList)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(result)).To(Equal(1))
		})

		It("should fail on invalid address", func() {
			accessList := transaction.AccessList{
				{
					Address:     "0xinvalid",
					StorageKeys: []string{},
				},
			}
			_, err := transaction.SerializeAccessList(accessList)
			Expect(err).To(HaveOccurred())
		})

		It("should fail on invalid storage key size", func() {
			accessList := transaction.AccessList{
				{
					Address: "0x1234567890123456789012345678901234567890",
					StorageKeys: []string{
						"0x01", // Too short
					},
				},
			}
			_, err := transaction.SerializeAccessList(accessList)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("AssertTransactionEIP1559", func() {
		It("should pass for valid EIP-1559 transaction", func() {
			tx := &transaction.Transaction{
				ChainId:              1,
				MaxFeePerGas:         big.NewInt(1000000000),
				MaxPriorityFeePerGas: big.NewInt(100000000),
				To:                   "0x1234567890123456789012345678901234567890",
			}
			err := transaction.AssertTransactionEIP1559(tx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail for invalid chain ID", func() {
			tx := &transaction.Transaction{
				ChainId:      0,
				MaxFeePerGas: big.NewInt(1000000000),
			}
			err := transaction.AssertTransactionEIP1559(tx)
			Expect(err).To(HaveOccurred())
		})

		It("should fail when tip exceeds fee cap", func() {
			tx := &transaction.Transaction{
				ChainId:              1,
				MaxFeePerGas:         big.NewInt(100000000),
				MaxPriorityFeePerGas: big.NewInt(1000000000), // Higher than maxFeePerGas
			}
			err := transaction.AssertTransactionEIP1559(tx)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SerializeTransaction and ParseTransaction", func() {
		It("should serialize and parse EIP-1559 transaction", func() {
			tx := &transaction.Transaction{
				Type:                 transaction.TransactionTypeEIP1559,
				ChainId:              1,
				Nonce:                0,
				MaxPriorityFeePerGas: big.NewInt(1000000000),
				MaxFeePerGas:         big.NewInt(2000000000),
				Gas:                  big.NewInt(21000),
				To:                   "0x1234567890123456789012345678901234567890",
				Value:                big.NewInt(1000000000000000000),
				Data:                 "0x",
			}

			serialized, err := transaction.SerializeTransaction(tx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(serialized).To(HavePrefix("0x02"))

			parsed, err := transaction.ParseTransaction(serialized)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.Type).To(Equal(transaction.TransactionTypeEIP1559))
			Expect(parsed.ChainId).To(Equal(1))
			Expect(parsed.Nonce).To(Equal(0))
			Expect(parsed.To).To(Equal("0x1234567890123456789012345678901234567890"))
		})

		It("should serialize and parse legacy transaction", func() {
			tx := &transaction.Transaction{
				Type:     transaction.TransactionTypeLegacy,
				ChainId:  1,
				Nonce:    0,
				GasPrice: big.NewInt(20000000000),
				Gas:      big.NewInt(21000),
				To:       "0x1234567890123456789012345678901234567890",
				Value:    big.NewInt(1000000000000000000),
				Data:     "0x",
			}

			serialized, err := transaction.SerializeTransaction(tx, nil)
			Expect(err).NotTo(HaveOccurred())

			parsed, err := transaction.ParseTransaction(serialized)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.Type).To(Equal(transaction.TransactionTypeLegacy))
			Expect(parsed.Nonce).To(Equal(0))
			Expect(parsed.To).To(Equal("0x1234567890123456789012345678901234567890"))
		})
	})
})
