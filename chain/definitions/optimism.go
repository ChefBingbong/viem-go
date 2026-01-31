package definitions

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
)

// Optimism is the OP Mainnet (Optimism) chain definition.
var Optimism = chain.DefineChain(chain.Chain{
	ID:   10,
	Name: "OP Mainnet",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "Ether",
		Symbol:   "ETH",
		Decimals: 18,
	},
	BlockTime: int64Ptr(2_000),
	SourceID:  int64Ptr(1), // mainnet L1
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://mainnet.optimism.io"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "Optimism Explorer",
			URL:    "https://optimistic.etherscan.io",
			ApiURL: "https://api-optimistic.etherscan.io/api",
		},
	},
	Contracts: &chain.ChainContracts{
		Multicall3: &chain.ChainContract{
			Address:      common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(4_286_263),
		},
	},
})
