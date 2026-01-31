package definitions

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
)

// Arbitrum is the Arbitrum One chain definition.
var Arbitrum = chain.DefineChain(chain.Chain{
	ID:   42_161,
	Name: "Arbitrum One",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "Ether",
		Symbol:   "ETH",
		Decimals: 18,
	},
	BlockTime: int64Ptr(250),
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://arb1.arbitrum.io/rpc"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "Arbiscan",
			URL:    "https://arbiscan.io",
			ApiURL: "https://api.arbiscan.io/api",
		},
	},
	Contracts: &chain.ChainContracts{
		Multicall3: &chain.ChainContract{
			Address:      common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(7_654_707),
		},
	},
})
