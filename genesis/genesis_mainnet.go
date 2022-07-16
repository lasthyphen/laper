// Copyright (C) 2019-2021, Dijets, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	"github.com/lasthyphen/beacongo/utils/units"
	"github.com/lasthyphen/beacongo/vms/platformvm/reward"
)

var (

	mainnetGenesisConfigJSON = `{
		"networkID": 1,
		"allocations": [
			{
				"ethAddr": "0xd86b355443158939c2f1b2A00961F8453b33E74E",
				"djtxAddr": "X-dijets1v8wat5z4cxh7wh873d7n6d9m6mpnynr8sgl059",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1668115272
					}
				]
			},
			{
				"ethAddr": "0x1B00f59fff05F6591c13e32740377eAE72661061",
				"djtxAddr": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1674115272
					}
				]
			},
			{
				"ethAddr": "0x80231567cD6E270c8360B80b97034Ec26dad83b8",
				"djtxAddr": "X-dijets16yd4ams4xdp9c6ht9zfnp90225ukwmqnj964sw",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1676115272
					}
				]
			},
			{
				"ethAddr": "0x2CF25c0323F52bb8d3f46EBc2DFcd34Ba057b781",
				"djtxAddr": "X-dijets1mc6ar37ggvvh80ezkwnnyrkfummc32k0z9v9la",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1688115272
					}
				]
			},
			{
				"ethAddr": "0x1a8398064cbee749c62cc2943da234f75350390e",
				"djtxAddr": "X-dijets1u3c5j2f3kg38eg6ryflljj3n7ark48gtr7s6dh",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1688115272
					}
				]
			},
			{
				"ethAddr": "0x62D4B8D6bdc42D4BB32B082c43fAdA8fF12E7130",
				"djtxAddr": "X-dijets1vw0kjghl4nz3g9v5hsgn3l0xvgqetm678gcv09",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1698115272
					}
				]
			},
			{
				"ethAddr": "0xA0f60016666800936c534aB78cF3D9656654fAdE",
				"djtxAddr": "X-dijets1seyylc9v8w9x5l6dpmf3se66hr37c34xnmd0r5",
				"initialAmount": 9285714000000000,
				"unlockSchedule": [
					{
						"amount": 500000000000000,
						"locktime": 1698115272
					}
				]
			}
		],
		"startTime": 1655861006,
		"initialStakeDuration": 31536000,
		"initialStakeDurationOffset": 54000,
		"initialStakedFunds": [
			"X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e"
		],
		"initialStakers": [
			{
				"nodeID": "NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg",
				"rewardAddress": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"delegationFee": 1000000
			},
			{
				"nodeID": "NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ",
				"rewardAddress": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"delegationFee": 500000
			},
			{
				"nodeID": "NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN",
				"rewardAddress": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"delegationFee": 250000
			},
			{
				"nodeID": "NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu",
				"rewardAddress": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"delegationFee": 125000
			},
			{
				"nodeID": "NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5",
				"rewardAddress": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"delegationFee": 62500
			}
		],
		"cChainGenesis": "{\"config\":{\"chainId\":98200,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC\":{\"balance\":\"0x295BE96E64066972000000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}",
		"message": "Azi Da Piyu."
	}`

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                 units.MilliDjtx,
			CreateAssetTxFee:      10 * units.MilliDjtx,
			CreateSubnetTxFee:     1 * units.Djtx,
			CreateBlockchainTxFee: 1 * units.Djtx,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloDjtx,
			MaxValidatorStake: 3 * units.MegaDjtx,
			MinDelegatorStake: 25 * units.Djtx,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  2 * 7 * 24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          100 * units.MegaDjtx,
			},
		},
	}
)
