// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"testing"
)

func TestCalculate(t *testing.T) {
	expression := "(1--3)+2*(3+-4)"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	result := calc.Tree().Calculate(1.0)
	if result-2 != 0 {
		t.Fatal("got incorrect result", result)
	}
}

func TestString(t *testing.T) {
	expression := "(((1 - -(3)) / 3) + (2 * (3 + -(4))))"
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
		t.Fatal("strings don't match", parsed)
	}
}

func TestDerivative(t *testing.T) {
	expression := "x^2"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	a := calc.Tree()
	da := a.Derivative()
	t.Log(da.String())
}

func TestSource(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	s := NewSource()
	r := s.Samples(5, rng)
	for _, v := range r {
		v.Root.Calculate(1)
	}
	r.Statistics(s)
}

func TestNewMode(t *testing.T) {
	expression := "4*x^3 + 2*x"
	Integrate(expression)
}
