// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/big"
	"math/rand"
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
	if calc.Tree().Calculate(big.NewInt(1)).Cmp(big.NewInt(2)) != 0 {
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

func TestGenerate(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	g := NewGaussian()
	s := Samples{}
	for i := 0; i < 33; i++ {
		expression := s.Generate(g, rng)
		t.Log(i, expression)
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
			t.Log(parsed)
			t.Log(expression)
			t.Fatal("strings don't match")
		}
	}
}

func TestRandomSearch(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	g := NewGaussian()
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
	s := Samples{}
	for i := 0; i < 1024; i++ {
		query := s.Generate(g, rng)
		y := &Calculator[uint32]{Buffer: query}
		err := y.Init()
		if err != nil {
			t.Fatal(err)
		}
		if err := y.Parse(); err != nil {
			t.Fatal(err)
		}
		b := y.Tree()

		fitness := big.NewInt(0)
		for j := 0; j < 2048; j++ {
			z := int64(rng.Intn(1024*1024) + 1)
			diff := big.NewInt(0).Sub(a.Calculate(big.NewInt(z)), b.Calculate(big.NewInt(z)))
			diff = diff.Abs(diff)
			fitness = fitness.Add(fitness, diff)
		}
		if fitness.Cmp(big.NewInt(0)) == 0 {
			t.Log("result", query)
			break
		}
	}
}
