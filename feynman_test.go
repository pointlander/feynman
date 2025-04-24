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
	expression := "(1--3)+2*(3+-4)"
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		t.Fatal(err)
	}
	if err := calc.Parse(); err != nil {
		t.Fatal(err)
	}
	result := calc.Tree().Calculate(big.NewFloat(1))
	if result.Cmp(big.NewFloat(2)) != 0 {
		t.Fatal("got incorrect result", result.String())
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

func TestGenerate(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	g := NewGaussian()
	s := Samples{}
	for i := 0; i < 33; i++ {
		s.Samples = append(s.Samples, Set{})
		expression := s.Generate(5, &g, rng)
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
	expression := "3*x^2"
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
		for k := 0; k < 1024; k++ {
			s.Samples = append(s.Samples, Set{})
			query := s.Generate(5, &g, rng)
			t.Log(k, query.String())
			b := query.Derivative()

			fitness := big.NewFloat(0)
			fit := func() *big.Float {
				defer func() {
					recover()
				}()
				z := float64(rng.Intn(256) + 1)
				aa := a.Calculate(big.NewFloat(z))
				bb := b.Calculate(big.NewFloat(z))
				diff := big.NewFloat(0).Sub(aa, bb)
				diff = diff.Mul(diff, diff)
				return diff
			}
			for j := 0; j < 256; j++ {
				fit := fit()
				if fit == nil {
					fit = big.NewFloat(1337)
				}
				fitness = fitness.Add(fitness, fit)
			}
			var set func(*Samples)
			set = func(samples *Samples) {
				if len(samples.Samples) == 0 {
					return
				}
				if samples.Samples[len(samples.Samples)-1].Fitness == nil {
					samples.Samples[len(samples.Samples)-1].Fitness = fitness
				}
				if samples.Left != nil {
					set(samples.Left)
				}
				if samples.Right != nil {
					set(samples.Right)
				}
			}
			set(&s)
			t.Log("fitness", fitness)
			if fitness.Cmp(big.NewFloat(0)) == 0 {
				t.Log("result", query)
				t.Log("dresult", b.String())
				break outer
			}
		}
		var re func(*Samples, *G)
		re = func(s *Samples, g *G) {
			length := len(s.Samples)
			sort.Slice(s.Samples, func(i, j int) bool {
				return s.Samples[i].Fitness.Cmp(s.Samples[j].Fitness) < 0
			})
			for k := 0; k < Width; k++ {
				sum, count := 0.0, 0.0
				for _, v := range s.Samples[:length/2] {
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
				for _, v := range s.Samples[:length/2] {
					for _, vv := range v.Set[k].Value {
						diff := avg - vv
						stddev += diff * diff
					}
				}
				stddev = math.Sqrt(stddev / count)
				g.G[k].Mean = avg
				g.G[k].Stddev = stddev
			}
			if s.Left != nil {
				if g.Left == nil {
					gg := NewGaussian()
					g.Left = &gg
				}
				re(s.Left, g.Left)
			}
			if s.Right != nil {
				if g.Right == nil {
					gg := NewGaussian()
					g.Right = &gg
				}
				re(s.Right, g.Right)
			}
		}
		re(&s, &g)
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
