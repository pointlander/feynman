// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

//go:generate peg -switch -inline calculator.peg

func Integrate(depth int, expression string) *Node {
	calc := &Calculator[uint32]{Buffer: expression}
	err := calc.Init()
	if err != nil {
		panic(err)
	}
	if err := calc.Parse(); err != nil {
		panic(err)
	}
	a := calc.Tree()
	seed := 1
	values := []float64{.01, -.01, .1, -.1, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5}
	cache := make([]float64, len(values))
	for i, z := range values {
		zz := map[string]float64{"x": z}
		cache[i] = a.Calculate(zz)
	}
	type Element struct {
		Index int
		Value float64
	}
	for {
		rng := rand.New(rand.NewSource(int64(seed)))
		s := NewSource()
		last := ""
		for {
			r := s.Samples(depth, rng)
			d := make([][]Element, len(values))
			for j, v := range r {
				b := v.Root.Derivative(map[string]bool{"x": true})
				for k := range values {
					z := map[string]float64{"x": values[k]}
					aa := cache[k]
					bb := b.Calculate(z)
					diff := aa - bb
					if math.IsInf(diff, 0) || math.IsNaN(diff) {
						r[j].Fitness = math.Inf(1)
						d[k] = append(d[k], Element{
							Index: v.Index,
							Value: math.Inf(1),
						})
					} else {
						d[k] = append(d[k], Element{
							Index: v.Index,
							Value: math.Abs(diff),
						})
					}
					if !(math.IsInf(r[j].Fitness, 0) || math.IsNaN(r[j].Fitness)) {
						r[j].Fitness += diff * diff
					}
				}
			}
			sort.Slice(r, func(i, j int) bool {
				return r[i].Fitness < r[j].Fitness
			})
			if r[0].Fitness == 0 {
				return r[0].Root
			}

			if last == r[0].Root.String() {
				break
			}
			last = r[0].Root.String()

			for k := range d {
				sort.Slice(d[k], func(i, j int) bool {
					return d[k][i].Value < d[k][j].Value
				})
			}
			index := 0
		outer:
			for i := range d[0] {
				for k := range d {
					if math.IsInf(d[k][i].Value, 0) {
						break outer
					}
				}
				index++
			}
			for k := range d {
				d[k] = d[k][:index]
			}

			if index > 0 {
				common := make(map[int]int)
				index /= 2
				for k := range d {
					for j := range index {
						common[d[k][j].Index]++
					}
				}
				sort.Slice(r, func(i, j int) bool {
					return common[r[i].Index] > common[r[j].Index]
				})
				r = r[:index]
			}
			r.Statistics(s)
		}
		seed++
	}
}

func main() {
	calc := &Calculator[uint32]{Buffer: "(x3*x^(x1/x2))/x4"}
	err := calc.Init()
	if err != nil {
		panic(err)
	}
	if err := calc.Parse(); err != nil {
		panic(err)
	}
	a := calc.Tree()
	b := a.Derivative(map[string]bool{"x": true})
	fmt.Println(b.Simplify())
}
