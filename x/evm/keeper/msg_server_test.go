package keeper_test

import (
	"math/big"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/evmos/ethermint/testutil"
	utiltx "github.com/evmos/ethermint/testutil/tx"
	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/evmos/ethermint/x/evm/types"
)

type MsgServerTestSuite struct {
	testutil.BaseTestSuiteWithAccount
	testContractAddr common.Address
}

func TestMsgServerTestSuite(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}

func (suite *MsgServerTestSuite) TestEthereumTx() {
	var (
		err             error
		msg             *types.MsgEthereumTx
		signer          ethtypes.Signer
		vmdb            *statedb.StateDB
		expectedGasUsed uint64
	)

	testCases := []struct {
		name     string
		malleate func()
		expErr   bool
	}{
		{
			"Deploy contract tx - insufficient gas",
			func() {
				msg, err = utiltx.CreateContractMsgTx(
					vmdb.GetNonce(suite.Address),
					signer,
					big.NewInt(1),
					suite.Address,
					suite.Signer,
				)
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"Transfer funds tx",
			func() {
				msg, _, err = newEthMsgTx(
					vmdb.GetNonce(suite.Address),
					suite.Address,
					suite.Signer,
					signer,
					ethtypes.AccessListTxType,
					nil,
					nil,
				)
				suite.Require().NoError(err)
				expectedGasUsed = params.TxGas
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest(suite.T())
			signer = ethtypes.LatestSignerForChainID(suite.App.EvmKeeper.ChainID())
			vmdb = suite.StateDB()

			tc.malleate()
			res, err := suite.App.EvmKeeper.EthereumTx(suite.Ctx, msg)
			if tc.expErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
			suite.Require().Equal(expectedGasUsed, res.GasUsed)
			suite.Require().False(res.Failed())
		})
	}
}

func (suite *MsgServerTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name:      "fail - invalid authority",
			request:   &types.MsgUpdateParams{Authority: "foobar"},
			expectErr: true,
		},
		{
			name: "pass - valid Update msg",
			request: &types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    types.DefaultParams(),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run("MsgUpdateParams", func() {
			suite.SetupTest(suite.T())
			_, err := suite.App.EvmKeeper.UpdateParams(suite.Ctx, tc.request)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *MsgServerTestSuite) TestUpdateZeroGas() {
	var testContractAddr = common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	testCases := []struct {
		name      string
		maleate   func()
		request   *types.MsgUpdateZeroGas
		checkFn   func()
		expectErr bool
	}{
		{
			name:      "fail - invalid authority",
			maleate:   nil,
			request:   &types.MsgUpdateZeroGas{Authority: "foobar"},
			checkFn:   nil,
			expectErr: true,
		},
		{
			name:    "pass - valid Update msg with AddItems",
			maleate: nil,
			request: &types.MsgUpdateZeroGas{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Metadata: types.UpdateZeroGasMetadata{
					AddItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x12341234"},
						},
					},
				},
			},
			checkFn: func() {
				zeroGas := suite.App.EvmKeeper.GetAllZeroGas(suite.Ctx)
				suite.Require().Len(zeroGas, 1)

				hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, testContractAddr.Bytes(), []byte{0x12, 0x34, 0x12, 0x34})
				suite.Require().True(hasZeroGas)
			},
			expectErr: false,
		},
		{
			name: "pass - valid Update msg with RemoveItems",
			maleate: func() {
				testMethodId := []byte{0x12, 0x34, 0x12, 0x34}
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, testContractAddr.Bytes(), testMethodId)
			},
			request: &types.MsgUpdateZeroGas{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Metadata: types.UpdateZeroGasMetadata{
					RemoveItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x12341234"},
						},
					},
				},
			},
			checkFn: func() {
				zeroGas := suite.App.EvmKeeper.GetAllZeroGas(suite.Ctx)
				suite.Require().Len(zeroGas, 0)

				hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, testContractAddr.Bytes(), []byte{0x12, 0x34, 0x12, 0x34})
				suite.Require().False(hasZeroGas)
			},
		},
		{
			name: "pass - valid Update msg with AddItems and RemoveItems",
			maleate: func() {
				testMethodId := []byte{0x12, 0x34, 0x12, 0x34}
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, testContractAddr.Bytes(), testMethodId)
			},
			request: &types.MsgUpdateZeroGas{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Metadata: types.UpdateZeroGasMetadata{
					AddItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x23452345"},
						},
					},
					RemoveItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x12341234"},
						},
					},
				},
			},
			checkFn: func() {
				zeroGas := suite.App.EvmKeeper.GetAllZeroGas(suite.Ctx)
				suite.Require().Len(zeroGas, 1)

				hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, testContractAddr.Bytes(), []byte{0x23, 0x45, 0x23, 0x45})
				suite.Require().True(hasZeroGas)
			},
			expectErr: false,
		},
		{
			name: "pass - valid Update msg with AddItems and RemoveItems with same signature",
			maleate: func() {
				testMethodId := []byte{0x12, 0x34, 0x12, 0x34}
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, testContractAddr.Bytes(), testMethodId)
			},
			request: &types.MsgUpdateZeroGas{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Metadata: types.UpdateZeroGasMetadata{
					AddItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x23452345"},
						},
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x34563456"},
						},
					},
					RemoveItems: []*types.ZeroGas{
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x12341234"},
						},
						{
							ContractAddress: testContractAddr.Hex(),
							Signatures:      []string{"0x34563456"},
						},
					},
				},
			},
			checkFn: func() {
				zeroGas := suite.App.EvmKeeper.GetAllZeroGas(suite.Ctx)
				suite.Require().Len(zeroGas, 1)

				hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, testContractAddr.Bytes(), []byte{0x23, 0x45, 0x23, 0x45})
				suite.Require().True(hasZeroGas)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run("MsgUpdateZeroGas", func() {
			suite.SetupTest(suite.T())
			if tc.maleate != nil {
				tc.maleate()
			}
			_, err := suite.App.EvmKeeper.UpdateZeroGas(suite.Ctx, tc.request)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
