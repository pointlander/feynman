// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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
				b := v.Root.Derivative()
				dd := make([]float64, 0, 8)
				for k := range values {
					z := values[k]
					aa := a.Calculate(z)
					bb := b.Calculate(z)
					diff := aa - bb
					dd = append(dd, diff*diff)
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
				min := len(d[0])
				for k := range d {
					/*count, sum := 0.0, 0.0
					for _, v := range d[k] {
						count++
						sum += v.Value
					}
					average := sum / count
					variance := 0.0
					for _, v := range d[k] {
						diff := average - v.Value
						variance += diff * diff
					}
					variance /= float64(len(d[k]))
					max, cut := 0.0, len(d[k])-1
					for i := 1; i < len(d[k])-1; i++ {
						avga, avgb := 0.0, 0.0
						vara, varb := 0.0, 0.0
						for j := 0; j < i; j++ {
							avga += d[k][j].Value
						}
						avga /= float64(i)
						for j := 0; j < i; j++ {
							diff := d[k][j].Value - avga
							vara += diff * diff
						}
						vara /= float64(i)
						for j := i; j < len(d[k]); j++ {
							avgb += d[k][j].Value
						}
						avgb /= float64(len(d[k]) - i)
						for j := i; j < len(d[k]); j++ {
							diff := d[k][j].Value - avgb
							varb += diff * diff
						}
						varb /= float64(len(d[k]) - i)
						reduction := variance - (vara + varb)
						if reduction > max {
							max, cut = reduction, i
						}
					}*/
					cut := len(d[k]) / 2
					for j := 0; j < cut; j++ {
						common[d[k][j].Index]++
					}
					if cut < min {
						min = cut
					}
				}
				sort.Slice(r, func(i, j int) bool {
					return common[r[i].Index] > common[r[j].Index]
				})
				r = r[:min]
			}
			r.Statistics(s)
		}
		seed++
	}
}

func main() {

}
