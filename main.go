// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	//"fmt"
	"math"
	"math/cmplx"
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
	for {
		rng := rand.New(rand.NewSource(int64(seed)))
		s := NewSource()
		last := ""
		points := 1
		for {
			r := s.Samples(depth, rng)
			for j, v := range r {
				b := v.Root.Derivative()
				for k := 0; k < points; k++ {
					//z := float64(k + 1)
					z := complex(rng.Float64(), rng.Float64())
					aa := a.CalculateComplex(z)
					//fmt.Println(b.String())
					bb := b.CalculateComplex(z)
					diff := aa - bb
					//fmt.Println(aa, bb, diff)
					if cmplx.IsInf(diff) || cmplx.IsNaN(diff) {
						r[j].Fitness = cmplx.Inf()
					}
					if !(cmplx.IsInf(r[j].Fitness) || cmplx.IsNaN(r[j].Fitness)) {
						r[j].Fitness += diff
					}
				}
			}
			sort.Slice(r, func(i, j int) bool {
				if cmplx.Abs(r[i].Fitness) < cmplx.Abs(r[j].Fitness) {
					return true
				} else if cmplx.Abs(r[i].Fitness) == cmplx.Abs(r[j].Fitness) {
					return math.Abs(cmplx.Phase(r[i].Fitness)) < math.Abs(cmplx.Phase(r[j].Fitness))
				}
				return false
			})
			if last == r[0].Root.String() {
				break
			}
			last = r[0].Root.String()
			if r[0].Fitness == 0 {
				if points > 3 {
					return r[0].Root
				}
				points++
			}
			index := 0
			for _, v := range r {
				if cmplx.IsInf(v.Fitness) {
					break
				}
				index++
			}
			r = r[:index]
			/*if index > 0 {
				count, sum := 0.0, 0.0
				for _, v := range r {
					count++
					sum += cmplx.Abs(v.Fitness)
				}
				average := sum / count
				variance := 0.0
				for _, v := range r {
					diff := average - cmplx.Abs(v.Fitness)
					variance += diff * diff
				}
				variance /= float64(len(r))
				max, cut := 0.0, 0
				for i := 1; i < len(r)-1; i++ {
					avga, avgb := 0.0, 0.0
					vara, varb := 0.0, 0.0
					for j := 0; j < i; j++ {
						avga += cmplx.Abs(r[j].Fitness)
					}
					avga /= float64(i)
					for j := 0; j < i; j++ {
						diff := cmplx.Abs(r[j].Fitness) - avga
						vara += diff * diff
					}
					vara /= float64(i)
					for j := i; j < len(r); j++ {
						avgb += cmplx.Abs(r[j].Fitness)
					}
					avgb /= float64(len(r) - i)
					for j := i; j < len(r); j++ {
						diff := cmplx.Abs(r[j].Fitness) - avgb
						varb += diff * diff
					}
					varb /= float64(len(r) - i)
					reduction := variance - (vara + varb)
					if reduction > max {
						max, cut = reduction, i
					}
				}
				r = r[:cut]
			}*/
			r.Statistics(s)
		}
		seed++
	}
}

func main() {

}
