// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	"github.com/lasthyphen/beacongo/utils/units"
	"github.com/lasthyphen/beacongo/vms/platformvm/reward"
)

var (

	localGenesisConfigJSON = `{
		"networkID": 12345,
		"allocations": [
			{
				"ethAddr": "0xd86b355443158939c2f1b2A00961F8453b33E74E",
				"djtxAddr": "X-dijets1v8wat5z4cxh7wh873d7n6d9m6mpnynr8sgl059",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 13999990000000000,
						"locktime": 1633824000
					}
				]
			},
			{
				"ethAddr": "0x1B00f59fff05F6591c13e32740377eAE72661061",
				"djtxAddr": "X-dijets1uxkc262hvgsvstguvawwaanmsh8zvtcsd64l5e",
				"initialAmount": 30000000000000000,
				"unlockSchedule": [
					{
						"amount": 20000000000000000
					},
					{
						"amount": 17999990000000000,
						"locktime": 1633824000
					}
				]
			},
			{
				"ethAddr": "0x80231567cD6E270c8360B80b97034Ec26dad83b8",
				"djtxAddr": "X-dijets16yd4ams4xdp9c6ht9zfnp90225ukwmqnj964sw",
				"initialAmount": 35000000000000000,
				"unlockSchedule": [
					{
						"amount": 4000000000000000,
						"locktime": 1633824000
					}
				]
			}
		],
		"startTime": 1630987200,
		"initialStakeDuration": 31536000,
		"initialStakeDurationOffset": 5400,
		"initialStakedFunds": [
			"X-dijets1v8wat5z4cxh7wh873d7n6d9m6mpnynr8sgl059"
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
		"cChainGenesis": "{\"config\":{\"chainId\":43112,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC\":{\"balance\":\"0x295BE96E64066972000000\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}",
		"message": "{{ fun_quote }}"
	}`

	// LocalParams are the params used for local networks
	LocalParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                 units.MilliDjtx,
			CreateAssetTxFee:      units.MilliDjtx,
			CreateSubnetTxFee:     100 * units.MilliDjtx,
			CreateBlockchainTxFee: 100 * units.MilliDjtx,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloDjtx,
			MaxValidatorStake: 3 * units.MegaDjtx,
			MinDelegatorStake: 25 * units.Djtx,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  24 * time.Hour,
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
