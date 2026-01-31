package definitions

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
)

// Bsc is the BNB Smart Chain definition.
var Bsc = chain.DefineChain(chain.Chain{
	ID:   56,
	Name: "BNB Smart Chain",
	NativeCurrency: chain.ChainNativeCurrency{
		Name:     "BNB",
		Symbol:   "BNB",
		Decimals: 18,
	},
	BlockTime: int64Ptr(750),
	RpcUrls: map[string]chain.ChainRpcUrls{
		"default": {
			HTTP: []string{"https://56.rpc.thirdweb.com"},
		},
	},
	BlockExplorers: map[string]chain.ChainBlockExplorer{
		"default": {
			Name:   "BscScan",
			URL:    "https://bscscan.com",
			ApiURL: "https://api.bscscan.com/api",
		},
	},
	Contracts: &chain.ChainContracts{
		Multicall3: &chain.ChainContract{
			Address:      common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11"),
			BlockCreated: uint64Ptr(15_921_452),
		},
	},
})
