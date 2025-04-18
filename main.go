// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate peg -switch -inline calculator.peg

import (
	"math/big"
)

func main() {
	expression := "( 1 - -3 ) / 3 + 2 * ( 3 + -4 ) + 3 % 2^2"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		panic(err)
	}
	if err := calc.Parse(); err != nil {
		panic(err)
	}
	if Calculate(calc.Eval()).Cmp(big.NewInt(2)) != 0 {
		panic("got incorrect result")
	}
}
