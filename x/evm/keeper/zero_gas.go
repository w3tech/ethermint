package keeper

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/evmos/ethermint/x/evm/types"
)

func (k Keeper) GetAllZeroGas(ctx sdk.Context) []types.ZeroGas {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.KeyPrefixZeroGas, nil)
	defer iter.Close()

	// Use map to group signatures by contract address
	addrMap := make(map[string]*types.ZeroGas)

	for ; iter.Valid(); iter.Next() {
		addr, sig := types.ParseZeroGasByAddrSigKey(iter.Key())
		addrStr := common.BytesToAddress(addr).Hex()
		sigStr := hex.EncodeToString(sig)

		if zeroGas, exists := addrMap[addrStr]; exists {
			// If contract address already exists, append the signature
			zeroGas.Signatures = append(zeroGas.Signatures, sigStr)
		} else {
			// If new contract address, create new ZeroGas entry
			addrMap[addrStr] = &types.ZeroGas{
				ContractAddress: addrStr,
				Signatures:      []string{sigStr},
			}
		}
	}

	// Convert map to slice
	result := make([]types.ZeroGas, 0, len(addrMap))
	for _, zeroGas := range addrMap {
		result = append(result, *zeroGas)
	}

	return result
}

// GetAllZeroGasSigsByAddr returns all the zero gas signatures by contract address.
func (k Keeper) GetAllZeroGasSigsByAddr(ctx sdk.Context, addr []byte) []string {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(types.GetZeroGasByAddrKey(addr), nil)
	defer iter.Close()

	var sigs []string
	for ; iter.Valid(); iter.Next() {
		_, sig := types.ParseZeroGasByAddrSigKey(iter.Key())
		sigStr := hex.EncodeToString(sig)
		sigs = append(sigs, sigStr)
	}
	return sigs
}

func (k Keeper) HasZeroGas(ctx sdk.Context, addr []byte, sig []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetZeroGasByAddrSigKey(addr, sig))
}

// SetZeroGas sets the EVM zero gas.
func (k Keeper) SetZeroGas(ctx sdk.Context, addr []byte, sig []byte) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetZeroGasByAddrSigKey(addr, sig), []byte{})
	return nil
}

func (k Keeper) DeleteZeroGas(ctx sdk.Context, addr []byte, sig []byte) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetZeroGasByAddrSigKey(addr, sig))
	return nil
}
