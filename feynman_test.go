// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/big"
	"testing"
)

func TestCalculate(t *testing.T) {
	expression := "(1--3)/3+2*(3+-4)+3%2^2"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	if calc.Tree().Calculate().Cmp(big.NewInt(2)) != 0 {
		t.Fatal("got incorrect result")
	}
}

func TestString(t *testing.T) {
	expression := "(1--3)/3+2*(3+-4)+3%2^2"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	parsed := calc.Tree().String()
	if parsed != expression {
		t.Fatal("strings don't match")
	}
}
