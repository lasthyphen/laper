// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"testing"
	"time"

	"github.com/lasthyphen/beacongo/ids"
	"github.com/lasthyphen/beacongo/vms/components/djtx"
	"github.com/lasthyphen/beacongo/vms/components/verify"
	"github.com/lasthyphen/beacongo/vms/platformvm/stakeable"
	"github.com/lasthyphen/beacongo/vms/secp256k1fx"
)

type dummyUnsignedTx struct {
	BaseTx
}

func (du *dummyUnsignedTx) SemanticVerify(vm *VM, parentState MutableState, stx *Tx) error {
	return nil
}

func TestSemanticVerifySpendUTXOs(t *testing.T) {
	vm, _, _ := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
		vm.ctx.Lock.Unlock()
	}()

	// The VM time during a test, unless [chainTimestamp] is set
	now := time.Unix(1607133207, 0)

	unsignedTx := dummyUnsignedTx{
		BaseTx: BaseTx{},
	}
	unsignedTx.Initialize([]byte{0}, []byte{1})

	// Note that setting [chainTimestamp] also set's the VM's clock.
	// Adjust input/output locktimes accordingly.
	tests := []struct {
		description string
		utxos       []*djtx.UTXO
		ins         []*djtx.TransferableInput
		outs        []*djtx.TransferableOutput
		creds       []verify.Verifiable
		fee         uint64
		assetID     ids.ID
		shouldErr   bool
	}{
		{
			description: "no inputs, no outputs, no fee",
			utxos:       []*djtx.UTXO{},
			ins:         []*djtx.TransferableInput{},
			outs:        []*djtx.TransferableOutput{},
			creds:       []verify.Verifiable{},
			fee:         0,
			assetID:     vm.ctx.DJTXAssetID,
			shouldErr:   false,
		},
		{
			description: "no inputs, no outputs, positive fee",
			utxos:       []*djtx.UTXO{},
			ins:         []*djtx.TransferableInput{},
			outs:        []*djtx.TransferableOutput{},
			creds:       []verify.Verifiable{},
			fee:         1,
			assetID:     vm.ctx.DJTXAssetID,
			shouldErr:   true,
		},
		{
			description: "no inputs, no outputs, positive fee",
			utxos:       []*djtx.UTXO{},
			ins:         []*djtx.TransferableInput{},
			outs:        []*djtx.TransferableOutput{},
			creds:       []verify.Verifiable{},
			fee:         1,
			assetID:     vm.ctx.DJTXAssetID,
			shouldErr:   true,
		},
		{
			description: "one input, no outputs, positive fee",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
		{
			description: "wrong number of credentials",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs:      []*djtx.TransferableOutput{},
			creds:     []verify.Verifiable{},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: true,
		},
		{
			description: "wrong number of UTXOs",
			utxos:       []*djtx.UTXO{},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: true,
		},
		{
			description: "invalid credential",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				(*secp256k1fx.Credential)(nil),
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: true,
		},
		{
			description: "one input, no outputs, positive fee",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &secp256k1fx.TransferInput{
					Amt: 1,
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
		{
			description: "locked one input, no outputs, no fee",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Unix()) + 1,
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &stakeable.LockIn{
					Locktime: uint64(now.Unix()) + 1,
					TransferableIn: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       0,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
		{
			description: "locked one input, no outputs, positive fee",
			utxos: []*djtx.UTXO{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				Out: &stakeable.LockOut{
					Locktime: uint64(now.Unix()) + 1,
					TransferableOut: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			}},
			ins: []*djtx.TransferableInput{{
				Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
				In: &stakeable.LockIn{
					Locktime: uint64(now.Unix()) + 1,
					TransferableIn: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			}},
			outs: []*djtx.TransferableOutput{},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: true,
		},
		{
			description: "one locked one unlock input, one locked output, positive fee",
			utxos: []*djtx.UTXO{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			ins: []*djtx.TransferableInput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					In: &stakeable.LockIn{
						Locktime: uint64(now.Unix()) + 1,
						TransferableIn: &secp256k1fx.TransferInput{
							Amt: 1,
						},
					},
				},
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*djtx.TransferableOutput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
		{
			description: "one locked one unlock input, one locked output, positive fee, partially locked",
			utxos: []*djtx.UTXO{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 2,
					},
				},
			},
			ins: []*djtx.TransferableInput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					In: &stakeable.LockIn{
						Locktime: uint64(now.Unix()) + 1,
						TransferableIn: &secp256k1fx.TransferInput{
							Amt: 1,
						},
					},
				},
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 2,
					},
				},
			},
			outs: []*djtx.TransferableOutput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) + 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 2,
						},
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
				&secp256k1fx.Credential{},
			},
			fee:       1,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
		{
			description: "one unlock input, one locked output, zero fee, unlocked",
			utxos: []*djtx.UTXO{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &stakeable.LockOut{
						Locktime: uint64(now.Unix()) - 1,
						TransferableOut: &secp256k1fx.TransferOutput{
							Amt: 1,
						},
					},
				},
			},
			ins: []*djtx.TransferableInput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
					},
				},
			},
			outs: []*djtx.TransferableOutput{
				{
					Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
					Out: &secp256k1fx.TransferOutput{
						Amt: 1,
					},
				},
			},
			creds: []verify.Verifiable{
				&secp256k1fx.Credential{},
			},
			fee:       0,
			assetID:   vm.ctx.DJTXAssetID,
			shouldErr: false,
		},
	}

	for _, test := range tests {
		vm.clock.Set(now)

		t.Run(test.description, func(t *testing.T) {
			err := vm.semanticVerifySpendUTXOs(
				&unsignedTx,
				test.utxos,
				test.ins,
				test.outs,
				test.creds,
				test.fee,
				test.assetID,
			)

			if err == nil && test.shouldErr {
				t.Fatalf("expected error but got none")
			} else if err != nil && !test.shouldErr {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
