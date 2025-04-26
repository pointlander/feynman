// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
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
	expression := "4*x^3"
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

			fitness := 0.0
			fit := func() float64 {
				z := float64(rng.Intn(256) + 1)
				aa := a.Calculate(z)
				bb := b.Calculate(z)
				diff := aa - bb
				return diff * diff
			}
			for j := 0; j < 256; j++ {
				fit := fit()
				if math.IsInf(fit, 0) || math.IsNaN(fit) {
					fit = 1337.0
				}
				fitness += fit
			}
			var set func(*Samples)
			set = func(samples *Samples) {
				if len(samples.Samples) == 0 {
					return
				}
				samples.Samples[len(samples.Samples)-1].Fitness = fitness
				if samples.Left != nil {
					set(samples.Left)
				}
				if samples.Right != nil {
					set(samples.Right)
				}
			}
			set(&s)
			t.Log("fitness", fitness)
			if fitness == 0 {
				t.Log("result", query)
				t.Log("dresult", b.String())
				break outer
			}
		}
		var re func(*Samples, *G)
		re = func(s *Samples, g *G) {
			length := len(s.Samples)
			sort.Slice(s.Samples, func(i, j int) bool {
				return s.Samples[i].Fitness < s.Samples[j].Fitness
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
