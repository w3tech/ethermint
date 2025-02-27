package v6

import (
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v5types "github.com/evmos/ethermint/x/evm/migrations/v6/types"
	"github.com/evmos/ethermint/x/evm/types"
)

// MigrateStore migrates the x/evm module state from the consensus version 5 to
// version 6. Specifically, it migrates the geth chain configuration
// that changed from geth v1.10 to v1.13.
func MigrateStore(
	ctx sdk.Context,
	storeService corestore.KVStoreService,
	cdc codec.BinaryCodec,
) error {
	var v5Params v5types.V5Params
	store := storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyPrefixParams)
	if err != nil {
		return err
	}
	cdc.MustUnmarshal(bz, &v5Params)

	updatedParams := types.Params{
		EvmDenom:     v5Params.EvmDenom,
		EnableCreate: v5Params.EnableCreate,
		EnableCall:   v5Params.EnableCall,
		ExtraEIPs:    v5Params.ExtraEIPs,
		ChainConfig: types.ChainConfig{
			HomesteadBlock:      v5Params.ChainConfig.HomesteadBlock,
			DAOForkBlock:        v5Params.ChainConfig.DAOForkBlock,
			DAOForkSupport:      v5Params.ChainConfig.DAOForkSupport,
			EIP150Block:         v5Params.ChainConfig.EIP150Block,
			EIP155Block:         v5Params.ChainConfig.EIP155Block,
			EIP158Block:         v5Params.ChainConfig.EIP158Block,
			ByzantiumBlock:      v5Params.ChainConfig.ByzantiumBlock,
			ConstantinopleBlock: v5Params.ChainConfig.ConstantinopleBlock,
			PetersburgBlock:     v5Params.ChainConfig.PetersburgBlock,
			IstanbulBlock:       v5Params.ChainConfig.IstanbulBlock,
			MuirGlacierBlock:    v5Params.ChainConfig.MuirGlacierBlock,
			BerlinBlock:         v5Params.ChainConfig.BerlinBlock,
			LondonBlock:         v5Params.ChainConfig.LondonBlock,
			ArrowGlacierBlock:   v5Params.ChainConfig.ArrowGlacierBlock,
			GrayGlacierBlock:    v5Params.ChainConfig.GrayGlacierBlock,
			MergeNetsplitBlock:  v5Params.ChainConfig.MergeNetsplitBlock,
			ShanghaiTime:        v5Params.ChainConfig.ShanghaiBlock,
			CancunTime:          v5Params.ChainConfig.CancunBlock,
		},
		AllowUnprotectedTxs: v5Params.AllowUnprotectedTxs,
	}

	if err := updatedParams.Validate(); err != nil {
		return err
	}
	updatedBz := cdc.MustMarshal(&updatedParams)
	err = store.Set(types.KeyPrefixParams, updatedBz)
	if err != nil {
		return err
	}

	return nil
}
