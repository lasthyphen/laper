// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"testing"

	"github.com/lasthyphen/beacongo/codec"
	"github.com/lasthyphen/beacongo/codec/linearcodec"
	"github.com/lasthyphen/beacongo/ids"
	"github.com/lasthyphen/beacongo/snow"
	"github.com/lasthyphen/beacongo/vms/components/djtx"
	"github.com/lasthyphen/beacongo/vms/components/verify"
)

type testOperable struct {
	djtx.TestTransferable `serialize:"true"`

	Outputs []verify.State `serialize:"true"`
}

func (o *testOperable) InitCtx(ctx *snow.Context) {}

func (o *testOperable) Outs() []verify.State { return o.Outputs }

func TestOperationVerifyNil(t *testing.T) {
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	op := (*Operation)(nil)
	if err := op.Verify(m); err == nil {
		t.Fatalf("Should have erred due to nil operation")
	}
}

func TestOperationVerifyEmpty(t *testing.T) {
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	op := &Operation{
		Asset: djtx.Asset{ID: ids.Empty},
	}
	if err := op.Verify(m); err == nil {
		t.Fatalf("Should have erred due to empty operation")
	}
}

func TestOperationVerifyUTXOIDsNotSorted(t *testing.T) {
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	op := &Operation{
		Asset: djtx.Asset{ID: ids.Empty},
		UTXOIDs: []*djtx.UTXOID{
			{
				TxID:        ids.Empty,
				OutputIndex: 1,
			},
			{
				TxID:        ids.Empty,
				OutputIndex: 0,
			},
		},
		Op: &testOperable{},
	}
	if err := op.Verify(m); err == nil {
		t.Fatalf("Should have erred due to unsorted utxoIDs")
	}
}

func TestOperationVerify(t *testing.T) {
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	assetID := ids.GenerateTestID()
	op := &Operation{
		Asset: djtx.Asset{ID: assetID},
		UTXOIDs: []*djtx.UTXOID{
			{
				TxID:        assetID,
				OutputIndex: 1,
			},
		},
		Op: &testOperable{},
	}
	if err := op.Verify(m); err != nil {
		t.Fatal(err)
	}
}

func TestOperationSorting(t *testing.T) {
	c := linearcodec.NewDefault()
	if err := c.RegisterType(&testOperable{}); err != nil {
		t.Fatal(err)
	}

	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	ops := []*Operation{
		{
			Asset: djtx.Asset{ID: ids.Empty},
			UTXOIDs: []*djtx.UTXOID{
				{
					TxID:        ids.Empty,
					OutputIndex: 1,
				},
			},
			Op: &testOperable{},
		},
		{
			Asset: djtx.Asset{ID: ids.Empty},
			UTXOIDs: []*djtx.UTXOID{
				{
					TxID:        ids.Empty,
					OutputIndex: 0,
				},
			},
			Op: &testOperable{},
		},
	}
	if IsSortedAndUniqueOperations(ops, m) {
		t.Fatalf("Shouldn't be sorted")
	}
	SortOperations(ops, m)
	if !IsSortedAndUniqueOperations(ops, m) {
		t.Fatalf("Should be sorted")
	}
	ops = append(ops, &Operation{
		Asset: djtx.Asset{ID: ids.Empty},
		UTXOIDs: []*djtx.UTXOID{
			{
				TxID:        ids.Empty,
				OutputIndex: 1,
			},
		},
		Op: &testOperable{},
	})
	if IsSortedAndUniqueOperations(ops, m) {
		t.Fatalf("Shouldn't be unique")
	}
}

func TestOperationTxNotState(t *testing.T) {
	intf := interface{}(&OperationTx{})
	if _, ok := intf.(verify.State); ok {
		t.Fatalf("shouldn't be marked as state")
	}
}
