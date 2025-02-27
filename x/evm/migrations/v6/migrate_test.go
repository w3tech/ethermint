package v6_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"

	"github.com/evmos/ethermint/encoding"
	"github.com/evmos/ethermint/x/evm/migrations/v6"
	v5types "github.com/evmos/ethermint/x/evm/migrations/v6/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func TestMigrate(t *testing.T) {
	encCfg := encoding.MakeTestEncodingConfig()
	cdc := encCfg.Codec

	storeKey := storetypes.NewKVStoreKey(types.ModuleName)
	tKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	storeService := runtime.NewKVStoreService(storeKey)

	// set v5 defaultParams
	defaultParams := types.DefaultParams()
	v5Params := v5types.V5Params{
		EvmDenom:     defaultParams.EvmDenom,
		EnableCreate: defaultParams.EnableCreate,
		EnableCall:   defaultParams.EnableCall,
		ExtraEIPs:    defaultParams.ExtraEIPs,
		ChainConfig: v5types.V5ChainConfig{
			HomesteadBlock:      defaultParams.ChainConfig.HomesteadBlock,
			DAOForkBlock:        defaultParams.ChainConfig.DAOForkBlock,
			DAOForkSupport:      defaultParams.ChainConfig.DAOForkSupport,
			EIP150Block:         defaultParams.ChainConfig.EIP150Block,
			EIP155Block:         defaultParams.ChainConfig.EIP155Block,
			EIP158Block:         defaultParams.ChainConfig.EIP158Block,
			ByzantiumBlock:      defaultParams.ChainConfig.ByzantiumBlock,
			ConstantinopleBlock: defaultParams.ChainConfig.ConstantinopleBlock,
			PetersburgBlock:     defaultParams.ChainConfig.PetersburgBlock,
			IstanbulBlock:       defaultParams.ChainConfig.IstanbulBlock,
			MuirGlacierBlock:    defaultParams.ChainConfig.MuirGlacierBlock,
			BerlinBlock:         defaultParams.ChainConfig.BerlinBlock,
			LondonBlock:         defaultParams.ChainConfig.LondonBlock,
			ArrowGlacierBlock:   defaultParams.ChainConfig.ArrowGlacierBlock,
			GrayGlacierBlock:    defaultParams.ChainConfig.GrayGlacierBlock,
			MergeNetsplitBlock:  defaultParams.ChainConfig.MergeNetsplitBlock,
			ShanghaiBlock:       defaultParams.ChainConfig.ShanghaiTime,
			CancunBlock:         defaultParams.ChainConfig.CancunTime,
		},
		AllowUnprotectedTxs: defaultParams.AllowUnprotectedTxs,
	}

	// Set the params in the store
	bz := cdc.MustMarshal(&v5Params)
	kvStore := storeService.OpenKVStore(ctx)
	kvStore.Set(types.KeyPrefixParams, bz)

	// Migrate the store
	err := v6.MigrateStore(ctx, storeService, cdc)
	require.NoError(t, err)

	var updatedParams types.Params
	paramsBz, err := kvStore.Get(types.KeyPrefixParams)
	require.NoError(t, err)
	cdc.MustUnmarshal(paramsBz, &updatedParams)

	// test that the params have been migrated correctly
	require.Equal(t, defaultParams, updatedParams)
}
