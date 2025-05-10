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
	expression := "(" + b.Simplify().String() + "- x5)^2"
	fmt.Println(expression)
	fmt.Println()
	{
		calc := &Calculator[uint32]{Buffer: expression}
		err := calc.Init()
		if err != nil {
			panic(err)
		}
		if err := calc.Parse(); err != nil {
			panic(err)
		}
		c := calc.Tree()
		partials := [4]*Node{}
		partials[0] = c.Derivative(map[string]bool{"x1": true}).Simplify()
		partials[1] = c.Derivative(map[string]bool{"x2": true}).Simplify()
		partials[2] = c.Derivative(map[string]bool{"x3": true}).Simplify()
		partials[3] = c.Derivative(map[string]bool{"x4": true}).Simplify()
		for _, v := range partials {
			fmt.Println(v)
		}

		rng := rand.New(rand.NewSource(1))
		values := map[string]float64{
			"x":  3.0,
			"x1": rng.Float64(),
			"x2": rng.Float64(),
			"x3": rng.Float64(),
			"x4": rng.Float64(),
			"x5": 9.0,
		}
		for i := 0; i < 33; i++ {
			dx := make([]float64, 0, 8)
			for _, v := range partials {
				dx = append(dx, v.Calculate(values))
			}
			sum := 0.0
			for _, v := range dx {
				sum += v * v
			}
			factor, length := 1.0, math.Sqrt(sum)
			fmt.Println(length)
			if length > 1 {
				factor /= length
			}
			values["x1"] = values["x1"] + .03*factor*dx[0]
			values["x2"] = values["x2"] + .03*factor*dx[1]
			values["x3"] = values["x3"] + .03*factor*dx[2]
			values["x4"] = values["x4"] + .03*factor*dx[3]
		}
		fmt.Println(values["x1"]/values["x2"], values["x3"]/values["x4"])
	}
}
