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

func Integrate(depth int, expression string) {
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
outer:
	for {
		fmt.Println("----------------------------------------------")
		rng := rand.New(rand.NewSource(int64(seed)))
		s := NewSource()
		last := ""
		points := 1
		for {
			r := s.Samples(depth, rng)
			for j, v := range r {
				b := v.Root.Derivative()
				for k := 0; k < points; k++ {
					z := float64(k + 1)
					aa := a.Calculate(z)
					bb := b.Calculate(z)
					diff := aa - bb
					if math.IsInf(diff, 0) || math.IsNaN(diff) {
						r[j].Fitness = math.Inf(1)
					}
					if !(math.IsInf(r[j].Fitness, 0) || math.IsNaN(r[j].Fitness)) {
						r[j].Fitness += diff * diff
					}
				}
			}
			sort.Slice(r, func(i, j int) bool {
				return r[i].Fitness < r[j].Fitness
			})
			if last == r[0].Root.String() {
				break
			}
			last = r[0].Root.String()
			fmt.Println(r[0].Fitness, r[0].Root.Simplify().String())
			if r[0].Fitness == 0 {
				if points > 2 {
					break outer
				}
				points++
			}
			index := 0
			for _, v := range r {
				if math.IsInf(v.Fitness, 0) {
					break
				}
				index++
			}
			r = r[:index]
			if index > 0 {
				count, sum := 0.0, 0.0
				for _, v := range r {
					count++
					sum += v.Fitness
				}
				average := sum / count
				variance := 0.0
				for _, v := range r {
					diff := average - v.Fitness
					variance += diff * diff
				}
				variance /= float64(len(r))
				max, cut := 0.0, 0
				for i := 1; i < len(r)-1; i++ {
					avga, avgb := 0.0, 0.0
					vara, varb := 0.0, 0.0
					for j := 0; j < i; j++ {
						avga += r[j].Fitness
					}
					avga /= float64(i)
					for j := 0; j < i; j++ {
						diff := r[j].Fitness - avga
						vara += diff * diff
					}
					vara /= float64(i)
					for j := i; j < len(r); j++ {
						avgb += r[j].Fitness
					}
					avgb /= float64(len(r) - i)
					for j := i; j < len(r); j++ {
						diff := r[j].Fitness - avgb
						varb += diff * diff
					}
					varb /= float64(len(r) - i)
					reduction := variance - (vara + varb)
					if reduction > max {
						max, cut = reduction, i
					}
				}
				r = r[:cut]
			}
			r.Statistics(s)
		}
		seed++
	}
}

func main() {

}
