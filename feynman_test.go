// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"math/big"
	"math/rand"
	"sort"
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
		s.Samples = append(s.Samples, Set{})
		expression := s.Generate(g, rng)
		t.Log(i, expression.String())
		parsed := expression.String()
		if parsed != expression.String() {
			t.Log(parsed)
			t.Log(expression)
			t.Fatal("strings don't match")
		}
	}
}

func TestRandomSearch(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	g := NewGaussian()
	expression := "2*x"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	a := calc.Tree()
outer:
	for i := 0; i < 1024; i++ {
		s := Samples{}
		for k := 0; k < 128; k++ {
			s.Samples = append(s.Samples, Set{})
			query := s.Generate(g, rng)
			t.Log(k, query.String())
			b := query.Derivative()

			fitness := big.NewInt(0)
			fit := func() *big.Int {
				defer func() {
					recover()
				}()
				z := int64(rng.Intn(1024) + 1)
				aa := a.Calculate(big.NewInt(z))
				bb := b.Calculate(big.NewInt(z))
				diff := big.NewInt(0).Sub(aa, bb)
				diff = diff.Mul(diff, diff)
				t.Log(aa, bb, diff)
				return diff
			}
			for j := 0; j < 256; j++ {
				fit := fit()
				if fit == nil {
					fit = big.NewInt(1337)
				}
				fitness = fitness.Add(fitness, fit)
			}
			s.Samples[len(s.Samples)-1].Fitness = fitness
			t.Log("fitness", fitness)
			if fitness.Cmp(big.NewInt(0)) == 0 {
				t.Log("result", query)
				t.Log("dresult", b.String())
				break outer
			}
		}
		sort.Slice(s.Samples, func(i, j int) bool {
			return s.Samples[i].Fitness.Cmp(s.Samples[j].Fitness) < 0
		})
		for k := 0; k < Width; k++ {
			sum, count := 0.0, 0.0
			for _, v := range s.Samples[:64] {
				for _, vv := range v.Set[k].Value {
					count++
					sum += vv
				}
			}
			if count < 7 {
				continue
			}
			avg := sum / count
			stddev := 0.0
			for _, v := range s.Samples[:64] {
				for _, vv := range v.Set[k].Value {
					diff := avg - vv
					stddev += diff * diff
				}
			}
			stddev = math.Sqrt(stddev / count)
			g[k].Mean = avg
			g[k].Stddev = stddev
		}
		t.Log(g)
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
