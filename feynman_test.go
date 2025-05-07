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
	x := map[string]float64{"x": 1.0}
	result := calc.Tree().Calculate(x)
	if result-2 != 0 {
		t.Fatal("got incorrect result", result)
	}
}

func TestSin(t *testing.T) {
	expression := "sin(pi)"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	x := map[string]float64{"x": 1.0}
	result := calc.Tree().Calculate(x)
	if result > 1e-10 {
		t.Log(calc.Tree().String())
		t.Fatal("got incorrect result", result)
	}
}

func TestCos(t *testing.T) {
	expression := "cos(pi)"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	x := map[string]float64{"x": 1.0}
	result := calc.Tree().Calculate(x)
	if result != -1 {
		t.Log(calc.Tree().String())
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
	x := map[string]float64{"x": 1.0}
	for _, v := range r {
		v.Root.Calculate(x)
	}
	r.Statistics(s)
}

func TestNewMode(t *testing.T) {
	expression := []string{
		"x",
		"2*x",
		"4*x^3",
		"x^3",
		"4*x^3 + 2*x",
		"2*x*cos(x^2)",
	}
	for _, e := range expression {
		t.Log(e)
		calc := &Calculator[uint32]{Buffer: e}
		err := calc.Init()
		if err != nil {
			panic(err)
		}
		if err := calc.Parse(); err != nil {
			panic(err)
		}
		input := calc.Tree()
		result := Integrate(5, e)
		result = result.Derivative()
		for i := 0; i < 256; i++ {
			z := map[string]float64{"x": float64(i + 1)}
			aa := input.Calculate(z)
			bb := result.Calculate(z)
			diff := aa - bb
			diff *= diff
			if diff != 0 {
				t.Fatalf("%s != %s", input, result)
			}
		}
		t.Log("done")
	}
}
