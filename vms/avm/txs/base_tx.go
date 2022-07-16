// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/lasthyphen/beacongo/codec"
	"github.com/lasthyphen/beacongo/ids"
	"github.com/lasthyphen/beacongo/snow"
	"github.com/lasthyphen/beacongo/vms/components/djtx"
)

var _ UnsignedTx = &BaseTx{}

// BaseTx is the basis of all transactions.
type BaseTx struct {
	djtx.BaseTx `serialize:"true"`
}

func (t *BaseTx) InitCtx(ctx *snow.Context) {
	for _, out := range t.Outs {
		out.InitCtx(ctx)
	}
}

func (t *BaseTx) SyntacticVerify(
	ctx *snow.Context,
	c codec.Manager,
	txFeeAssetID ids.ID,
	txFee uint64,
	_ uint64,
	_ int,
) error {
	if t == nil {
		return errNilTx
	}

	if err := t.MetadataVerify(ctx); err != nil {
		return err
	}

	return djtx.VerifyTx(
		txFee,
		txFeeAssetID,
		[][]*djtx.TransferableInput{t.Ins},
		[][]*djtx.TransferableOutput{t.Outs},
		c,
	)
}

func (t *BaseTx) Visit(v Visitor) error {
	return v.BaseTx(t)
}
