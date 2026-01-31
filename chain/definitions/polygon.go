package definitions

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
)

// Polygon is the Polygon PoS chain definition.
var Polygon = chain.DefineChain(chain.Chain{
	ID:   137,
	Name: "Polygon",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "POL",
		Symbol:   "POL",
		Decimals: 18,
	},
	BlockTime: int64Ptr(2_000),
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://polygon-rpc.com"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "PolygonScan",
			URL:    "https://polygonscan.com",
			ApiURL: "https://api.etherscan.io/v2/api",
		},
	},
	Contracts: &chain.ChainContracts{
		Multicall3: &chain.ChainContract{
			Address:      common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(25_770_160),
		},
	},
})
