package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
)

func main() {
	c, _ := client.NewClient("https://eth.merkle.io")
	defer c.Close()

	start := time.Now()

	result, err := c.ReadContract(context.Background(), client.ReadContractOptions{
		Address:      common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		ABI:          `[{"name":"balanceOf","type":"function","inputs":[{"name":"owner","type":"address"}],"outputs":[{"type":"uint256"}]}]`,
		FunctionName: "balanceOf",
		Args:         []any{common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")},
	})

	elapsed := time.Since(start)

	if err != nil {
		panic(err)
	}
	balance := result[0].(*big.Int)
	fmt.Printf("Balance: %s\n", balance)
	fmt.Printf("Time: %s\n", elapsed)
}
