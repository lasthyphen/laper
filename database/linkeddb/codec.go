// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package linkeddb

import (
	"math"

	"github.com/lasthyphen/beacongo/codec"
	"github.com/lasthyphen/beacongo/codec/linearcodec"
)

const (
	codecVersion = 0
)

// c does serialization and deserialization
var (
	c codec.Manager
)

func init() {
	lc := linearcodec.NewCustomMaxLength(math.MaxUint32)
	c = codec.NewManager(math.MaxInt32)

	if err := c.RegisterCodec(codecVersion, lc); err != nil {
		panic(err)
	}
}
