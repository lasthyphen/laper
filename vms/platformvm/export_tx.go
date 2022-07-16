// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"errors"
	"fmt"

	"github.com/lasthyphen/beacongo/chains/atomic"
	"github.com/lasthyphen/beacongo/database"
	"github.com/lasthyphen/beacongo/ids"
	"github.com/lasthyphen/beacongo/snow"
	"github.com/lasthyphen/beacongo/utils/crypto"
	"github.com/lasthyphen/beacongo/utils/math"
	"github.com/lasthyphen/beacongo/vms/components/djtx"
	"github.com/lasthyphen/beacongo/vms/components/verify"
	"github.com/lasthyphen/beacongo/vms/platformvm/stakeable"
	"github.com/lasthyphen/beacongo/vms/secp256k1fx"
)

var (
	errNoExportOutputs  = errors.New("no export outputs")
	errOutputsNotSorted = errors.New("outputs not sorted")
	errOverflowExport   = errors.New("overflow when computing export amount + txFee")

	_ UnsignedAtomicTx = &UnsignedExportTx{}
)

// UnsignedExportTx is an unsigned ExportTx
type UnsignedExportTx struct {
	BaseTx `serialize:"true"`

	// Which chain to send the funds to
	DestinationChain ids.ID `serialize:"true" json:"destinationChain"`

	// Outputs that are exported to the chain
	ExportedOutputs []*djtx.TransferableOutput `serialize:"true" json:"exportedOutputs"`
}

// InitCtx sets the FxID fields in the inputs and outputs of this
// [UnsignedExportTx]. Also sets the [ctx] to the given [vm.ctx] so that
// the addresses can be json marshalled into human readable format
func (tx *UnsignedExportTx) InitCtx(ctx *snow.Context) {
	tx.BaseTx.InitCtx(ctx)
	for _, out := range tx.ExportedOutputs {
		out.FxID = secp256k1fx.ID
		out.InitCtx(ctx)
	}
}

// InputUTXOs returns an empty set
func (tx *UnsignedExportTx) InputUTXOs() ids.Set { return nil }

// SyntacticVerify this transaction is well-formed
func (tx *UnsignedExportTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return errNilTx
	case tx.syntacticallyVerified: // already passed syntactic verification
		return nil
	case len(tx.ExportedOutputs) == 0:
		return errNoExportOutputs
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return err
	}

	for _, out := range tx.ExportedOutputs {
		if err := out.Verify(); err != nil {
			return fmt.Errorf("output failed verification: %w", err)
		}
		if _, ok := out.Output().(*stakeable.LockOut); ok {
			return errWrongLocktime
		}
	}
	if !djtx.IsSortedTransferableOutputs(tx.ExportedOutputs, Codec) {
		return errOutputsNotSorted
	}

	tx.syntacticallyVerified = true
	return nil
}

// Attempts to verify this transaction with the provided state.
func (tx *UnsignedExportTx) SemanticVerify(vm *VM, parentState MutableState, stx *Tx) error {
	_, err := tx.AtomicExecute(vm, parentState, stx)
	return err
}

// Execute this transaction.
func (tx *UnsignedExportTx) Execute(
	vm *VM,
	vs VersionedState,
	stx *Tx,
) (func() error, error) {
	if err := tx.SyntacticVerify(vm.ctx); err != nil {
		return nil, err
	}

	outs := make([]*djtx.TransferableOutput, len(tx.Outs)+len(tx.ExportedOutputs))
	copy(outs, tx.Outs)
	copy(outs[len(tx.Outs):], tx.ExportedOutputs)

	if vm.bootstrapped.GetValue() {
		if err := verify.SameSubnet(vm.ctx, tx.DestinationChain); err != nil {
			return nil, err
		}
	}

	// Verify the flowcheck
	if err := vm.semanticVerifySpend(vs, tx, tx.Ins, outs, stx.Creds, vm.TxFee, vm.ctx.DJTXAssetID); err != nil {
		return nil, fmt.Errorf("failed semanticVerifySpend: %w", err)
	}

	// Consume the UTXOS
	consumeInputs(vs, tx.Ins)
	// Produce the UTXOS
	txID := tx.ID()
	produceOutputs(vs, txID, vm.ctx.DJTXAssetID, tx.Outs)
	return nil, nil
}

// AtomicOperations returns the shared memory requests
func (tx *UnsignedExportTx) AtomicOperations() (ids.ID, *atomic.Requests, error) {
	txID := tx.ID()

	elems := make([]*atomic.Element, len(tx.ExportedOutputs))
	for i, out := range tx.ExportedOutputs {
		utxo := &djtx.UTXO{
			UTXOID: djtx.UTXOID{
				TxID:        txID,
				OutputIndex: uint32(len(tx.Outs) + i),
			},
			Asset: djtx.Asset{ID: out.AssetID()},
			Out:   out.Out,
		}

		utxoBytes, err := Codec.Marshal(CodecVersion, utxo)
		if err != nil {
			return ids.ID{}, nil, fmt.Errorf("failed to marshal UTXO: %w", err)
		}
		utxoID := utxo.InputID()
		elem := &atomic.Element{
			Key:   utxoID[:],
			Value: utxoBytes,
		}
		if out, ok := utxo.Out.(djtx.Addressable); ok {
			elem.Traits = out.Addresses()
		}

		elems[i] = elem
	}
	return tx.DestinationChain, &atomic.Requests{PutRequests: elems}, nil
}

// Execute this transaction and return the versioned state.
func (tx *UnsignedExportTx) AtomicExecute(
	vm *VM,
	parentState MutableState,
	stx *Tx,
) (VersionedState, error) {
	// Set up the state if this tx is committed
	newState := newVersionedState(
		parentState,
		parentState.CurrentStakerChainState(),
		parentState.PendingStakerChainState(),
	)
	_, err := tx.Execute(vm, newState, stx)
	return newState, err
}

// Accept this transaction.
func (tx *UnsignedExportTx) AtomicAccept(ctx *snow.Context, batch database.Batch) error {
	chainID, requests, err := tx.AtomicOperations()
	if err != nil {
		return err
	}
	return ctx.SharedMemory.Apply(map[ids.ID]*atomic.Requests{chainID: requests}, batch)
}

// Create a new transaction
func (vm *VM) newExportTx(
	amount uint64, // Amount of tokens to export
	chainID ids.ID, // Chain to send the UTXOs to
	to ids.ShortID, // Address of chain recipient
	keys []*crypto.PrivateKeySECP256K1R, // Pay the fee and provide the tokens
	changeAddr ids.ShortID, // Address to send change to, if there is any
) (*Tx, error) {
	toBurn, err := math.Add64(amount, vm.TxFee)
	if err != nil {
		return nil, errOverflowExport
	}
	ins, outs, _, signers, err := vm.stake(keys, 0, toBurn, changeAddr)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate tx inputs/outputs: %w", err)
	}

	// Create the transaction
	utx := &UnsignedExportTx{
		BaseTx: BaseTx{BaseTx: djtx.BaseTx{
			NetworkID:    vm.ctx.NetworkID,
			BlockchainID: vm.ctx.ChainID,
			Ins:          ins,
			Outs:         outs, // Non-exported outputs
		}},
		DestinationChain: chainID,
		ExportedOutputs: []*djtx.TransferableOutput{{ // Exported to X-Chain
			Asset: djtx.Asset{ID: vm.ctx.DJTXAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: amount,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Threshold: 1,
					Addrs:     []ids.ShortID{to},
				},
			},
		}},
	}
	tx := &Tx{UnsignedTx: utx}
	if err := tx.Sign(Codec, signers); err != nil {
		return nil, err
	}
	return tx, utx.SyntacticVerify(vm.ctx)
}
