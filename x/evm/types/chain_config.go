// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ethereum/go-ethereum/params"
)

// EthereumConfig returns an Ethereum ChainConfig for EVM state transitions.
// All the negative or nil values are converted to nil
func (cc ChainConfig) EthereumConfig(chainID *big.Int) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                       chainID,
		HomesteadBlock:                getBlockValue(cc.HomesteadBlock),
		DAOForkBlock:                  getBlockValue(cc.DAOForkBlock),
		DAOForkSupport:                cc.DAOForkSupport,
		EIP150Block:                   getBlockValue(cc.EIP150Block),
		EIP155Block:                   getBlockValue(cc.EIP155Block),
		EIP158Block:                   getBlockValue(cc.EIP158Block),
		ByzantiumBlock:                getBlockValue(cc.ByzantiumBlock),
		ConstantinopleBlock:           getBlockValue(cc.ConstantinopleBlock),
		PetersburgBlock:               getBlockValue(cc.PetersburgBlock),
		IstanbulBlock:                 getBlockValue(cc.IstanbulBlock),
		MuirGlacierBlock:              getBlockValue(cc.MuirGlacierBlock),
		BerlinBlock:                   getBlockValue(cc.BerlinBlock),
		LondonBlock:                   getBlockValue(cc.LondonBlock),
		ArrowGlacierBlock:             getBlockValue(cc.ArrowGlacierBlock),
		GrayGlacierBlock:              getBlockValue(cc.GrayGlacierBlock),
		MergeNetsplitBlock:            getBlockValue(cc.MergeNetsplitBlock),
		ShanghaiTime:                  getTimestampValue(cc.ShanghaiTime),
		CancunTime:                    getTimestampValue(cc.CancunTime),
		PragueTime:                    getTimestampValue(cc.PragueTime),
		VerkleTime:                    getTimestampValue(cc.VerkleTime),
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: true,
		Ethash:                        nil,
		Clique:                        nil,
	}
}

// DefaultChainConfig returns default evm parameters.
func DefaultChainConfig() ChainConfig {
	homesteadBlock := sdkmath.ZeroInt()
	daoForkBlock := sdkmath.ZeroInt()
	eip150Block := sdkmath.ZeroInt()
	eip155Block := sdkmath.ZeroInt()
	eip158Block := sdkmath.ZeroInt()
	byzantiumBlock := sdkmath.ZeroInt()
	constantinopleBlock := sdkmath.ZeroInt()
	petersburgBlock := sdkmath.ZeroInt()
	istanbulBlock := sdkmath.ZeroInt()
	muirGlacierBlock := sdkmath.ZeroInt()
	berlinBlock := sdkmath.ZeroInt()
	londonBlock := sdkmath.ZeroInt()
	arrowGlacierBlock := sdkmath.ZeroInt()
	grayGlacierBlock := sdkmath.ZeroInt()
	mergeNetsplitBlock := sdkmath.ZeroInt()
	shanghaiTime := sdkmath.ZeroInt()
	cancunTime := sdkmath.ZeroInt()

	return ChainConfig{
		HomesteadBlock:      &homesteadBlock,
		DAOForkBlock:        &daoForkBlock,
		DAOForkSupport:      true,
		EIP150Block:         &eip150Block,
		EIP155Block:         &eip155Block,
		EIP158Block:         &eip158Block,
		ByzantiumBlock:      &byzantiumBlock,
		ConstantinopleBlock: &constantinopleBlock,
		PetersburgBlock:     &petersburgBlock,
		IstanbulBlock:       &istanbulBlock,
		MuirGlacierBlock:    &muirGlacierBlock,
		BerlinBlock:         &berlinBlock,
		LondonBlock:         &londonBlock,
		ArrowGlacierBlock:   &arrowGlacierBlock,
		GrayGlacierBlock:    &grayGlacierBlock,
		MergeNetsplitBlock:  &mergeNetsplitBlock,
		ShanghaiTime:        &shanghaiTime,
		CancunTime:          &cancunTime,
		PragueTime:          nil,
		VerkleTime:          nil,
	}
}

func getBlockValue(block *sdkmath.Int) *big.Int {
	if block == nil || block.IsNegative() {
		return nil
	}
	return block.BigInt()
}

func getTimestampValue(ts *sdkmath.Int) *uint64 {
	if ts == nil || ts.IsNegative() {
		return nil
	}
	res := ts.Uint64()
	return &res
}

// Validate performs a basic validation of the ChainConfig params. The function will return an error
// if any of the block values is uninitialized (i.e nil) or if the EIP150Hash is an invalid hash.
func (cc ChainConfig) Validate() error {
	if err := validateBlockOrTimestamp(cc.HomesteadBlock); err != nil {
		return errorsmod.Wrap(err, "homesteadBlock")
	}
	if err := validateBlockOrTimestamp(cc.DAOForkBlock); err != nil {
		return errorsmod.Wrap(err, "daoForkBlock")
	}
	if err := validateBlockOrTimestamp(cc.EIP150Block); err != nil {
		return errorsmod.Wrap(err, "eip150Block")
	}
	if err := validateBlockOrTimestamp(cc.EIP155Block); err != nil {
		return errorsmod.Wrap(err, "eip155Block")
	}
	if err := validateBlockOrTimestamp(cc.EIP158Block); err != nil {
		return errorsmod.Wrap(err, "eip158Block")
	}
	if err := validateBlockOrTimestamp(cc.ByzantiumBlock); err != nil {
		return errorsmod.Wrap(err, "byzantiumBlock")
	}
	if err := validateBlockOrTimestamp(cc.ConstantinopleBlock); err != nil {
		return errorsmod.Wrap(err, "constantinopleBlock")
	}
	if err := validateBlockOrTimestamp(cc.PetersburgBlock); err != nil {
		return errorsmod.Wrap(err, "petersburgBlock")
	}
	if err := validateBlockOrTimestamp(cc.IstanbulBlock); err != nil {
		return errorsmod.Wrap(err, "istanbulBlock")
	}
	if err := validateBlockOrTimestamp(cc.MuirGlacierBlock); err != nil {
		return errorsmod.Wrap(err, "muirGlacierBlock")
	}
	if err := validateBlockOrTimestamp(cc.BerlinBlock); err != nil {
		return errorsmod.Wrap(err, "berlinBlock")
	}
	if err := validateBlockOrTimestamp(cc.LondonBlock); err != nil {
		return errorsmod.Wrap(err, "londonBlock")
	}
	if err := validateBlockOrTimestamp(cc.ArrowGlacierBlock); err != nil {
		return errorsmod.Wrap(err, "arrowGlacierBlock")
	}
	if err := validateBlockOrTimestamp(cc.GrayGlacierBlock); err != nil {
		return errorsmod.Wrap(err, "GrayGlacierBlock")
	}
	if err := validateBlockOrTimestamp(cc.MergeNetsplitBlock); err != nil {
		return errorsmod.Wrap(err, "MergeNetsplitBlock")
	}
	if err := validateBlockOrTimestamp(cc.ShanghaiTime); err != nil {
		return errorsmod.Wrap(err, "ShanghaiTime")
	}
	if err := validateBlockOrTimestamp(cc.CancunTime); err != nil {
		return errorsmod.Wrap(err, "CancunTime")
	}
	if err := validateBlockOrTimestamp(cc.PragueTime); err != nil {
		return errorsmod.Wrap(err, "PragueTime")
	}
	if err := validateBlockOrTimestamp(cc.VerkleTime); err != nil {
		return errorsmod.Wrap(err, "VerkleTime")
	}
	// NOTE: chain ID is not needed to check config order
	if err := cc.EthereumConfig(nil).CheckConfigForkOrder(); err != nil {
		return errorsmod.Wrap(err, "invalid config fork order")
	}
	return nil
}

func validateBlockOrTimestamp(value *sdkmath.Int) error {
	// nil value means that the fork has not yet been applied
	if value == nil {
		return nil
	}

	if value.IsNegative() {
		return errorsmod.Wrapf(
			ErrInvalidChainConfig, "block or timestamp value cannot be negative: %s", value,
		)
	}

	return nil
}
