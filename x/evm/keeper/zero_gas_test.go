package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/app"
	"github.com/evmos/ethermint/testutil"
	"github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/suite"
)

type ZeroGasTestSuite struct {
	testutil.EVMTestSuiteWithAccountAndQueryClient
	testContractAddr common.Address
	enableFeemarket  bool
	enableLondonHF   bool
}

func (suite *ZeroGasTestSuite) SetupTest() {
	suite.EVMTestSuiteWithAccountAndQueryClient.SetupTestWithCb(suite.T(), func(app *app.EthermintApp, genesis app.GenesisState) app.GenesisState {
		feemarketGenesis := feemarkettypes.DefaultGenesisState()
		if suite.enableFeemarket {
			feemarketGenesis.Params.EnableHeight = 1
			feemarketGenesis.Params.NoBaseFee = false
		} else {
			feemarketGenesis.Params.NoBaseFee = true
		}
		genesis[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(feemarketGenesis)
		if !suite.enableLondonHF {
			evmGenesis := types.DefaultGenesisState()
			maxInt := sdkmath.NewInt(math.MaxInt64)
			evmGenesis.Params.ChainConfig.LondonBlock = &maxInt
			evmGenesis.Params.ChainConfig.ArrowGlacierBlock = &maxInt
			evmGenesis.Params.ChainConfig.GrayGlacierBlock = &maxInt
			evmGenesis.Params.ChainConfig.MergeNetsplitBlock = &maxInt
			evmGenesis.Params.ChainConfig.ShanghaiTime = &maxInt
			genesis[types.ModuleName] = app.AppCodec().MustMarshalJSON(evmGenesis)
		}
		return genesis
	})

	suite.testContractAddr = suite.deployTestContract(suite.Address)
	suite.Commit(suite.T())
}

func TestZeroGasTestSuite(t *testing.T) {
	s := new(ZeroGasTestSuite)
	s.enableFeemarket = false
	s.enableLondonHF = true
	suite.Run(t, s)
}

// deployTestContract deploy a test erc20 contract and returns the contract address
func (suite *ZeroGasTestSuite) deployTestContract(owner common.Address) common.Address {
	supply := sdkmath.NewIntWithDecimal(1000, 18).BigInt()
	return suite.EVMTestSuiteWithAccountAndQueryClient.DeployTestContract(
		suite.T(),
		owner,
		supply,
		suite.enableFeemarket,
	)
}

func (suite *ZeroGasTestSuite) TestGetAllZeroGas() {
	testCases := []struct {
		msg      string
		malleate func()
		expected []types.ZeroGas
	}{
		{
			"get all zero gas methods",
			func() {
				// register ERC20 transfer method as zero gas method
				signature := types.ERC20Contract.ABI.Methods["transfer"].ID
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
				suite.Commit(suite.T())
			},
			[]types.ZeroGas{
				{
					ContractAddress: suite.testContractAddr.String(),
					Signatures: []string{
						hex.EncodeToString(types.ERC20Contract.ABI.Methods["transfer"].ID),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			zeroGas := suite.App.EvmKeeper.GetAllZeroGas(suite.Ctx)
			suite.Require().Len(zeroGas, 1)
		})
	}
}

func (suite *ZeroGasTestSuite) TestGetAllZeroGasSigsByAddr() {
	testCases := []struct {
		msg      string
		malleate func()
		expected []string
	}{
		{
			"get all zero gas methods by contract address",
			func() {
				// register ERC20 transfer method as zero gas method
				signature := types.ERC20Contract.ABI.Methods["transfer"].ID
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
				suite.Commit(suite.T())
			},
			[]string{
				hex.EncodeToString(types.ERC20Contract.ABI.Methods["transfer"].ID),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			zeroGas := suite.App.EvmKeeper.GetAllZeroGasSigsByAddr(suite.Ctx, suite.testContractAddr.Bytes())
			suite.Require().Equal(tc.expected, zeroGas)
		})
	}
}

func (suite *ZeroGasTestSuite) TestHasZeroGas() {
	testCases := []struct {
		msg      string
		malleate func()
		expected bool
	}{
		{
			"check if a zero gas method exists",
			func() {
				// register ERC20 transfer method as zero gas method
				signature := types.ERC20Contract.ABI.Methods["transfer"].ID
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
				suite.Commit(suite.T())
			},
			true,
		},
		{
			"check if a zero gas method does not exist",
			func() {},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			tc.malleate()
			hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), types.ERC20Contract.ABI.Methods["transfer"].ID)
			suite.Require().Equal(tc.expected, hasZeroGas)
		})
	}
}

func (suite *ZeroGasTestSuite) TestSetZeroGas() {
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"set zero gas method",
			func() {
				// register ERC20 transfer method as zero gas method
				signature := types.ERC20Contract.ABI.Methods["transfer"].ID
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
				suite.Commit(suite.T())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), types.ERC20Contract.ABI.Methods["transfer"].ID)
			suite.Require().True(hasZeroGas)
		})
	}
}

func (suite *ZeroGasTestSuite) TestDeleteZeroGas() {
	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"delete zero gas method",
			func() {
				// register ERC20 transfer method as zero gas method
				signature := types.ERC20Contract.ABI.Methods["transfer"].ID
				suite.App.EvmKeeper.SetZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
				suite.Commit(suite.T())
				suite.App.EvmKeeper.DeleteZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), signature)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			hasZeroGas := suite.App.EvmKeeper.HasZeroGas(suite.Ctx, suite.testContractAddr.Bytes(), types.ERC20Contract.ABI.Methods["transfer"].ID)
			suite.Require().False(hasZeroGas)
		})
	}
}
