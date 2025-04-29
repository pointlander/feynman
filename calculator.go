// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

const (
	// Width is the number of random variables
	Width = 9
)

// Operation is a mathematical operation
type Operation uint

const (
	// OperationNoop is a noop
	OperationNoop Operation = iota
	// OperationAdd adds two numbers
	OperationAdd
	// OperationSubtract subtracts two numbers
	OperationSubtract
	// OperationMultiply multiplies two numbers
	OperationMultiply
	// OperationDivide divides two numbers
	OperationDivide
	// OperationExponentiation raises a number to a number
	OperationExponentiation
	// OperationNegate changes the sign of a number
	OperationNegate
	// OperationNumber is a real number
	OperationNumber
	// OperationVariable is a variable
	OperationVariable
	// OperationModulus computes the modulus of two numbers
	OperationModulus
	// OperationImaginary is an imaginary number
	OperationImaginary
	// OperationNaturalExponentiation raises the natural number to a power
	OperationNaturalExponentiation
	// OperationNatural is the constant e
	OperationNatural
	// OperationPI is the constant pi
	OperationPI
	// OperationNaturalLogarithm os the natural logarithm
	OperationNaturalLogarithm
	// OperationSquareRoot computes the square root of a number
	OperationSquareRoot
	// OperationCosine computes the cosine of a number
	OperationCosine
	// OperationSine computes the sine of a number
	OperationSine
	// OperationTangent computes the tangent of a number
	OperationTangent
	// OperationNotation is E notation operation
	OperationNotation
)

// Value is a value
type Value struct {
	ValueCount    float64
	ValueSum      [2]float64
	ValueVariance [2]float64
	Value         [2]Gaussian
}

// MarkvoValue is a markov model for value
type MarkovValue map[State]*Value

// Source is the source of nodes
type Source struct {
	OperationCount    float64
	OperationSum      [4]float64
	OperationVariance [4]float64
	Operation         [4]Gaussian
	Value             MarkovValue
}

// State is a markov state
type State [2]byte

// Markov is a markov model
type Markov map[State]*Source

// Node is a node in an expression
type Node struct {
	OperationSample [4]float64
	Operation       Operation
	ValueSample     [64][2]float64
	Value           float64
	Variable        string
	Left            *Node
	Right           *Node
}

// Root is the root node
type Root struct {
	Root    *Node
	Fitness float64
}

// Roots is a set of roots
type Roots []Root

// Sample is a sample
type Sample struct {
	Value []float64
}

// Set is a set of samples
type Set struct {
	Set     [Width]Sample
	Fitness float64
}

// Samples is a set of samples
type Samples struct {
	Samples []Set
	Left    *Samples
	Right   *Samples
}

// Gaussian is a gaussian
type Gaussian struct {
	Mean   float64
	Stddev float64
}

// G is a guassian set
type G struct {
	G     [Width]Gaussian
	Left  *G
	Right *G
}

// NewSource creates a new source markov model
func NewSource() Markov {
	source := make(Markov, Width*Width)
	for x := 0; x < Width; x++ {
		for y := 0; y < Width; y++ {
			s := Source{}
			for i := range s.Operation {
				s.Operation[i].Stddev = 1
			}
			s.Value = make(MarkovValue, 3*3)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					value := Value{}
					for k := range value.Value {
						value.Value[k].Stddev = 1
					}
					s.Value[State{byte(i), byte(j)}] = &value
				}
			}
			source[State{byte(x), byte(y)}] = &s
		}
	}
	return source
}

// NewGaussian makes a new gaussian distribution
func NewGaussian() (g G) {
	for i := range g.G {
		g.G[i].Stddev = 1
	}
	return g
}

// Sample samples from the source
func (m Markov) Sample(depth int, state State, rng *rand.Rand) *Node {
	n := Node{}
	depth--
	operation := Operation(0)
	for {
		for i := range m[state].Operation {
			operation <<= 1
			sample := rng.NormFloat64()*m[state].Operation[i].Stddev + m[state].Operation[i].Mean
			if sample > 0 {
				operation |= 1
			}
			n.OperationSample[i] = sample
		}
		if (operation != Operation(state[0]) && operation != Operation(state[1]) && operation > 0 && operation < Width && depth != 0) || ((operation == OperationVariable || operation == OperationNumber) && depth == 0) {
			break
		}
		operation = Operation(0)
	}
	n.Operation = operation
	value, ss := uint64(0), State{}
	for b := 0; b < 64; b++ {
		bits := 0
		for {
			for i := range m[state].Value[ss].Value {
				bits <<= 1
				sample := rng.NormFloat64()*m[state].Value[ss].Value[i].Stddev + m[state].Value[ss].Value[i].Mean
				if sample > 0 {
					bits |= 1
				}
				n.ValueSample[b][i] = sample
			}
			if bits < 3 {
				break
			}
			bits = 0
		}
		ss[0], ss[1] = byte(bits), ss[0]
		if bits == 2 {
			break
		}
		value <<= 1
		if bits == 1 {
			value |= 1
		}
	}
	n.Value = float64(value)
	n.Variable = "x"
	if depth == 0 || operation == OperationVariable || operation == OperationNumber {
		return &n
	}
	next := state
	next[0], next[1] = byte(operation), next[0]
	n.Left = m.Sample(depth, next, rng)
	n.Right = m.Sample(depth, next, rng)
	return &n
}

// Samples generates multiple samples
func (m Markov) Samples(depth int, rng *rand.Rand) Roots {
	root := Roots{}
	for i := 0; i < 1024; i++ {
		root = append(root, Root{
			Root: m.Sample(depth, State{}, rng),
		})
	}
	return root
}

// Reset resets a source
func (m Markov) Reset() {
	for _, s := range m {
		s.OperationCount = 0
		for i := range s.OperationSum {
			s.OperationSum[i] = 0
		}
		for i := range s.OperationVariance {
			s.OperationVariance[i] = 0
		}
		for _, s := range s.Value {
			s.ValueCount = 0
			for i := range s.ValueSum {
				s.ValueSum[i] = 0
			}
			for i := range s.ValueVariance {
				s.ValueVariance[i] = 0
			}
		}
	}
}

// Statistics computes the statistics of Roots
func (r Roots) Statistics(m Markov) {
	m.Reset()
	var sum func(State, Markov, *Node)
	sum = func(state State, m Markov, n *Node) {
		s := m[state]
		s.OperationCount++
		for i := range s.OperationSum {
			s.OperationSum[i] += n.OperationSample[i]
		}
		ss := State{}
		s.Value[ss].ValueCount++
	outer:
		for i := range n.ValueSample {
			bits := 0
			for j := range s.Value[ss].ValueSum {
				bits <<= 1
				s.Value[ss].ValueSum[j] += n.ValueSample[i][j]
				if n.ValueSample[i][j] > 0 {
					bits |= 1
				}
			}
			if bits == 2 {
				break outer
			}
			ss[0], ss[1] = byte(bits), ss[0]
		}
		next := state
		next[0], next[1] = byte(n.Operation), next[0]
		if n.Left != nil {
			sum(next, m, n.Left)
		}
		if n.Right != nil {
			sum(next, m, n.Right)
		}
	}
	for _, v := range r {
		sum(State{}, m, v.Root)
	}

	var avg func(State, Markov, *Node)
	avg = func(state State, m Markov, n *Node) {
		s := m[state]
		if s.OperationCount > 2 {
			for i := range s.OperationSum {
				s.Operation[i].Mean = s.OperationSum[i] / s.OperationCount
			}
		}
		ss := State{}
	outer:
		for i := range n.ValueSample {
			s := s.Value[ss]
			if s.ValueCount > 2 {
				bits := 0
				for j := range s.ValueSum {
					bits <<= 1
					s.Value[j].Mean = s.ValueSum[j] / s.ValueCount
					if n.ValueSample[i][j] > 0 {
						bits |= 1
					}
				}
				if bits == 2 {
					break outer
				}
				ss[0], ss[1] = byte(bits), ss[0]
			}
		}

		next := state
		next[0], next[1] = byte(n.Operation), next[0]
		if n.Left != nil {
			avg(next, m, n.Left)
		}
		if n.Right != nil {
			avg(next, m, n.Right)
		}
	}
	for _, v := range r {
		avg(State{}, m, v.Root)
	}

	var variance func(State, Markov, *Node)
	variance = func(state State, m Markov, n *Node) {
		s := m[state]
		if s.OperationCount > 2 {
			for i := range s.Operation {
				diff := s.Operation[i].Mean - n.OperationSample[i]
				s.OperationVariance[i] += diff * diff
			}
		}
		ss := State{}
	outer:
		for i := range n.ValueSample {
			s := s.Value[ss]
			if s.ValueCount > 2 {
				bits := 0
				for j := range s.Value {
					bits <<= 1
					diff := s.Value[j].Mean - n.ValueSample[i][j]
					s.ValueVariance[j] = diff * diff
					if n.ValueSample[i][j] > 0 {
						bits |= 1
					}
				}
				if bits == 2 {
					break outer
				}
				ss[0], ss[1] = byte(bits), ss[0]
			}
		}
		next := state
		next[0], next[1] = byte(n.Operation), next[0]
		if n.Left != nil {
			variance(next, m, n.Left)
		}
		if n.Right != nil {
			variance(next, m, n.Right)
		}
	}
	for _, v := range r {
		variance(State{}, m, v.Root)
	}

	var stddev func(State, Markov, *Node)
	stddev = func(state State, m Markov, n *Node) {
		s := m[state]
		if s.OperationCount > 2 {
			for i := range s.Operation {
				s.Operation[i].Stddev = math.Sqrt(s.OperationVariance[i] / s.OperationCount)
			}
		}
		ss := State{}
	outer:
		for i := range n.ValueSample {
			s := s.Value[ss]
			if s.ValueCount > 2 {
				bits := 0
				for j := range s.Value {
					bits <<= 1
					s.Value[j].Stddev = math.Sqrt(s.ValueVariance[j] / s.ValueCount)
					if n.ValueSample[i][j] > 0 {
						bits |= 1
					}
				}
				if bits == 2 {
					break outer
				}
				ss[0], ss[1] = byte(bits), ss[0]
			}
		}
		next := state
		next[0], next[1] = byte(n.Operation), next[0]
		if n.Left != nil {
			stddev(next, m, n.Left)
		}
		if n.Right != nil {
			stddev(next, m, n.Right)
		}
	}
	for _, v := range r {
		stddev(State{}, m, v.Root)
	}
}

// Generate generates an equation
func (s *Samples) Generate(depth int, g *G, rng *rand.Rand) *Node {
	if depth == 0 {
		return &Node{
			Operation: OperationNumber,
			Value:     1.0,
		}
	}
	depth--

	generate := func(vv int, samples *Set) *Node {
		sample := rng.NormFloat64()*g.G[2+vv].Stddev + g.G[2+vv].Mean
		samples.Set[2+vv].Value = append(samples.Set[2+vv].Value, sample)
		left := g.Left
		if left == nil {
			left = g
		}
		right := g.Right
		if right == nil {
			right = g
		}
		if s.Left == nil {
			s.Left = &Samples{}
		}
		if s.Right == nil {
			s.Right = &Samples{}
		}
		if sample > 0 {
			s.Left.Samples = append(s.Left.Samples, Set{})
			s.Right.Samples = append(s.Right.Samples, Set{})
			switch vv {
			case 0:
				return &Node{
					Operation: OperationAdd,
					Left:      s.Left.Generate(depth, left, rng),
					Right:     s.Right.Generate(depth, right, rng),
				}
			case 1:

				return &Node{
					Operation: OperationSubtract,
					Left:      s.Left.Generate(depth, left, rng),
					Right:     s.Right.Generate(depth, right, rng),
				}

			case 2:
				return &Node{
					Operation: OperationMultiply,
					Left:      s.Left.Generate(depth, left, rng),
					Right:     s.Right.Generate(depth, right, rng),
				}
			case 3:
				return &Node{
					Operation: OperationDivide,
					Left:      s.Left.Generate(depth, left, rng),
					Right:     s.Right.Generate(depth, right, rng),
				}
			case 4:
				return &Node{
					Operation: OperationExponentiation,
					Left:      s.Left.Generate(depth, left, rng),
					Right:     s.Right.Generate(depth, right, rng),
				}
			case 5:
				return &Node{
					Operation: OperationNegate,
					Left:      s.Left.Generate(depth, left, rng),
				}
			}
		}
		return nil
	}

	x := rng.Perm(3)
	y := rng.Perm(6)
	samples := &s.Samples[len(s.Samples)-1]
	for _, v := range x {
		switch v {
		case 0:
			sample := rng.NormFloat64()*g.G[0].Stddev + g.G[0].Mean
			samples.Set[0].Value = append(samples.Set[0].Value, sample)
			if sample > 0 {
				return &Node{
					Operation: OperationVariable,
					Variable:  "x",
				}
			}
		case 1:
			x := 1
			sample := rng.NormFloat64()*g.G[1].Stddev + g.G[1].Mean
			for sample > 0 {
				x++
				//samples.Set[1].Value = append(samples.Set[1].Value, sample)
				sample = rng.NormFloat64()*g.G[1].Stddev + g.G[1].Mean
				if sample < 0 {
					return &Node{
						Operation: OperationNumber,
						Value:     float64(x),
					}
				}
			}
		case 2:
			for _, vv := range y {
				result := generate(vv, samples)
				if result != nil {
					return result
				}
			}
		}
	}

	return &Node{
		Operation: OperationNumber,
		Value:     1.0,
	}
}

func (c *Calculator[_]) Tree() *Node {
	return c.Rulee(c.AST())
}

func (c *Calculator[U]) Rulee(node *node[U]) *Node {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return nil
}

func (c *Calculator[U]) Rulee1(node *node[U]) *Node {
	node = node.up
	var a *Node
	for node != nil {
		switch node.pegRule {
		case rulee2:
			a = c.Rulee2(node)
		case ruleadd:
			node = node.next
			b := &Node{}
			b.Operation = OperationAdd
			b.Left = a
			b.Right = c.Rulee2(node)
			a = b
		case ruleminus:
			node = node.next
			b := &Node{}
			b.Operation = OperationSubtract
			b.Left = a
			b.Right = c.Rulee2(node)
			a = b
		}
		node = node.next
	}
	return a
}

func (c *Calculator[U]) Rulee2(node *node[U]) *Node {
	node = node.up
	var a *Node
	for node != nil {
		switch node.pegRule {
		case rulee3:
			a = c.Rulee3(node)
		case rulemultiply:
			node = node.next
			b := &Node{}
			b.Operation = OperationMultiply
			b.Left = a
			b.Right = c.Rulee3(node)
			a = b
		case ruledivide:
			node = node.next
			b := &Node{}
			b.Operation = OperationDivide
			b.Left = a
			b.Right = c.Rulee3(node)
			a = b
		case rulemodulus:
			node = node.next
			b := &Node{}
			b.Operation = OperationModulus
			b.Left = a
			b.Right = c.Rulee3(node)
			a = b
		}
		node = node.next
	}
	return a
}

func (c *Calculator[U]) Rulee3(node *node[U]) *Node {
	node = node.up
	var a *Node
	for node != nil {
		switch node.pegRule {
		case rulee4:
			a = c.Rulee4(node)
		case ruleexponentiation:
			node = node.next
			b := &Node{}
			b.Operation = OperationExponentiation
			b.Left = a
			b.Right = c.Rulee4(node)
			a = b
		}
		node = node.next
	}
	return a
}

func (c *Calculator[U]) Rulee4(node *node[U]) *Node {
	node = node.up
	minus := false
	for node != nil {
		switch node.pegRule {
		case rulevalue:
			if minus {
				e := &Node{}
				e.Operation = OperationNegate
				e.Left = c.Rulevalue(node)
				return e
			}
			return c.Rulevalue(node)
		case ruleminus:
			minus = true
		}
		node = node.next
	}
	return nil
}

func (c *Calculator[U]) Rulevalue(node *node[U]) *Node {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulenumber:
			a := &Node{}
			a.Operation = OperationNumber
			value, err := strconv.ParseFloat(strings.TrimSpace(string(c.buffer[node.begin:node.end])), 64)
			if err != nil {
				panic(err)
			}
			a.Value = value
			return a
		case rulevariable:
			a := &Node{}
			a.Operation = OperationVariable
			a.Variable = string(c.buffer[node.begin:node.end])
			return a
		case rulesub:
			return c.Rulesub(node)
		}
		node = node.next
	}
	return nil
}

func (c *Calculator[U]) Rulesub(node *node[U]) *Node {
	node = node.up
	for node != nil {
		switch node.pegRule {
		case rulee1:
			return c.Rulee1(node)
		}
		node = node.next
	}
	return nil
}

// Derivative takes the derivative of the equation
// https://www.cs.utexas.edu/users/novak/asg-symdif.html#:~:text=Introduction,numeric%20calculations%20based%20on%20formulas.
func (n *Node) Derivative() *Node {
	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		switch n.Operation {
		case OperationNoop:
			return n
		case OperationAdd:
			a := &Node{
				Operation: OperationAdd,
				Left:      process(n.Left),
				Right:     process(n.Right),
			}
			return a
		case OperationSubtract:
			a := &Node{
				Operation: OperationSubtract,
				Left:      process(n.Left),
				Right:     process(n.Right),
			}
			return a
		case OperationMultiply:
			left := &Node{
				Operation: OperationMultiply,
				Left:      n.Left,
				Right:     process(n.Right),
			}
			right := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     process(n.Left),
			}
			a := &Node{
				Operation: OperationAdd,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationDivide:
			left := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     process(n.Left),
			}
			right := &Node{
				Operation: OperationMultiply,
				Left:      n.Left,
				Right:     process(n.Right),
			}
			difference := &Node{
				Operation: OperationSubtract,
				Left:      left,
				Right:     right,
			}
			square := &Node{
				Operation: OperationExponentiation,
				Left:      n.Right,
				Right: &Node{
					Operation: OperationNumber,
					Value:     2.0,
				},
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      difference,
				Right:     square,
			}
			return a
		case OperationModulus:
			return n
		case OperationExponentiation:
			one := &Node{
				Operation: OperationNumber,
				Value:     1.0,
			}
			subtract := &Node{
				Operation: OperationSubtract,
				Left:      n.Right,
				Right:     one,
			}
			exp := &Node{
				Operation: OperationExponentiation,
				Left:      n.Left,
				Right:     subtract,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      n.Right,
				Right:     exp,
			}
			a = &Node{
				Operation: OperationMultiply,
				Left:      a,
				Right:     process(n.Left),
			}
			return a
		case OperationNegate:
			a := &Node{
				Operation: OperationNegate,
				Left:      process(n.Left),
			}
			return a
		case OperationVariable:
			a := &Node{
				Operation: OperationNumber,
				Value:     1.0,
			}
			return a
		case OperationImaginary:
			a := &Node{
				Operation: OperationNumber,
				Value:     0.0,
			}
			return a
		case OperationNumber:
			a := &Node{
				Operation: OperationNumber,
				Value:     0.0,
			}
			return a
		case OperationNotation:
			a := &Node{
				Operation: OperationNumber,
				Value:     0.0,
			}
			return a
		case OperationNaturalExponentiation:
			a := &Node{
				Operation: OperationMultiply,
				Left:      n,
				Right:     process(n.Left),
			}
			return a
		case OperationNatural:
			a := &Node{
				Operation: OperationNumber,
				Value:     0.0,
			}
			return a
		case OperationPI:
			a := &Node{
				Operation: OperationNumber,
				Value:     0.0,
			}
			return a
		case OperationNaturalLogarithm:
			a := &Node{
				Operation: OperationDivide,
				Left:      process(n.Left),
				Right:     n.Left,
			}
			return a
		case OperationSquareRoot:
			value2 := &Node{
				Operation: OperationNumber,
				Value:     2.0,
			}
			multiply := &Node{
				Operation: OperationMultiply,
				Left:      value2,
				Right:     n,
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      process(n.Left),
				Right:     multiply,
			}
			return a
		case OperationCosine:
			sin := &Node{
				Operation: OperationSine,
				Left:      n.Left,
			}
			multiply := &Node{
				Operation: OperationMultiply,
				Left:      sin,
				Right:     process(n.Left),
			}
			a := &Node{
				Operation: OperationNegate,
				Left:      multiply,
			}
			return a
		case OperationSine:
			cos := &Node{
				Operation: OperationCosine,
				Left:      n.Left,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      cos,
				Right:     process(n.Left),
			}
			return a
		case OperationTangent:
			value1 := &Node{
				Operation: OperationNumber,
				Value:     1.0,
			}
			value2 := &Node{
				Operation: OperationNumber,
				Value:     2.0,
			}
			exp := &Node{
				Operation: OperationExponentiation,
				Left:      n,
				Right:     value2,
			}
			add := &Node{
				Operation: OperationAdd,
				Left:      value1,
				Right:     exp,
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      add,
				Right:     process(n.Left),
			}
			return a
		}
		return nil
	}
	return process(n)
}

var numeric = map[Operation]bool{
	OperationNumber:    true,
	OperationImaginary: true,
	OperationNotation:  true,
}

func isNumeric(operation Operation) bool {
	return numeric[operation]
}

// Simplify simplifies an expression
func (n *Node) Simplify() *Node {
	var process func(n *Node) *Node
	process = func(n *Node) *Node {
		if n == nil {
			return nil
		}
		switch n.Operation {
		case OperationNoop:
			return n
		case OperationAdd:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				return right
			} else if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			a := &Node{
				Operation: OperationAdd,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationSubtract:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNegate,
					Left:      right,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				return left
			}
			a := &Node{
				Operation: OperationSubtract,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationMultiply:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				return right
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationMultiply,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationDivide:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     math.Inf(1),
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationDivide,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationModulus:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationModulus,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationExponentiation:
			left, right := process(n.Left), process(n.Right)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     1.0,
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationNumber,
					Value:     1.0,
				}
				return a
			} else if isNumeric(right.Operation) && right.Equals(1) {
				return left
			}
			a := &Node{
				Operation: OperationExponentiation,
				Left:      left,
				Right:     right,
			}
			return a
		case OperationNegate:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			}
			a := &Node{
				Operation: OperationNegate,
				Left:      left,
			}
			return a
		case OperationVariable:
			return n
		case OperationImaginary:
			return n
		case OperationNumber:
			return n
		case OperationNotation:
			return n
		case OperationNaturalExponentiation:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     1.0,
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationVariable,
					Value:     math.E,
				}
				return a
			}
			a := &Node{
				Operation: OperationNaturalExponentiation,
				Left:      left,
			}
			return a
		case OperationNatural:
			return n
		case OperationPI:
			return n
		case OperationNaturalLogarithm:
			left := process(n.Left)
			if left.Operation == OperationNatural {
				return left
			}
			a := &Node{
				Operation: OperationNaturalLogarithm,
				Left:      left,
			}
			return a
		case OperationSquareRoot:
			left := process(n.Left)
			if isNumeric(left.Operation) && left.Equals(0) {
				a := &Node{
					Operation: OperationNumber,
					Value:     0.0,
				}
				return a
			} else if isNumeric(left.Operation) && left.Equals(1) {
				a := &Node{
					Operation: OperationNumber,
					Value:     1.0,
				}
				return a
			}
			a := &Node{
				Operation: OperationSquareRoot,
				Left:      left,
			}
			return a
		case OperationCosine:
			a := &Node{
				Operation: OperationCosine,
				Left:      process(n.Left),
			}
			return a
		case OperationSine:
			a := &Node{
				Operation: OperationSine,
				Left:      process(n.Left),
			}
			return a
		case OperationTangent:
			a := &Node{
				Operation: OperationTangent,
				Left:      process(n.Left),
			}
			return a
		}
		return nil
	}
	return process(n)
}

func (n *Node) Calculate(x float64) float64 {
	var a float64
	switch n.Operation {
	case OperationNumber:
		a = n.Value
	case OperationVariable:
		a = x
	case OperationNegate:
		a = -n.Left.Calculate(x)
	case OperationAdd:
		a = n.Left.Calculate(x) + n.Right.Calculate(x)
	case OperationSubtract:
		a = n.Left.Calculate(x) - n.Right.Calculate(x)
	case OperationMultiply:
		a = n.Left.Calculate(x) * n.Right.Calculate(x)
	case OperationDivide:
		a = n.Left.Calculate(x) / n.Right.Calculate(x)
	case OperationExponentiation:
		a = math.Pow(n.Left.Calculate(x), n.Right.Calculate(x))
	}
	return a
}

// Equals test if value is equal to x
func (n *Node) Equals(x int64) bool {
	/*if n.Operation == OperationNotation {
		a := big.NewInt(0)
		a.SetString(n.Left.Value, 10)
		b := big.NewInt(10)
		c := big.NewInt(0)
		c.SetString(n.Right.Value, 10)
		b.Exp(b, c, nil)
		a.Mul(a, b)
		return a.Cmp(big.NewInt(x)) == 0
	}*/
	return n.Value == float64(x)
}

// String returns the string form of the equation
func (n *Node) String() string {
	var process func(n *Node) string
	process = func(n *Node) string {
		if n == nil {
			return ""
		}
		switch n.Operation {
		case OperationNoop:
			return "(" + process(n.Left) + "???" + process(n.Right) + ")"
		case OperationAdd:
			return "(" + process(n.Left) + " + " + process(n.Right) + ")"
		case OperationSubtract:
			return "(" + process(n.Left) + " - " + process(n.Right) + ")"
		case OperationMultiply:
			return "(" + process(n.Left) + " * " + process(n.Right) + ")"
		case OperationDivide:
			return "(" + process(n.Left) + " / " + process(n.Right) + ")"
		case OperationModulus:
			return "(" + process(n.Left) + " % " + process(n.Right) + ")"
		case OperationExponentiation:
			return "(" + process(n.Left) + "^" + process(n.Right) + ")"
		case OperationNegate:
			return "-(" + process(n.Left) + ")"
		case OperationVariable:
			return n.Variable
		case OperationImaginary:
			return strconv.FormatFloat(n.Value, 'f', -1, 64) + "i"
		case OperationNumber:
			return strconv.FormatFloat(n.Value, 'f', -1, 64)
		case OperationNotation:
			if n.Left.Operation == OperationImaginary {
				return strconv.FormatFloat(n.Left.Value, 'f', -1, 64) + "e" + process(n.Right) + "i"
			}
			return process(n.Left) + "e" + process(n.Right)
		case OperationNaturalExponentiation:
			return "(e^" + process(n.Left) + ")"
		case OperationNatural:
			return "e"
		case OperationPI:
			return "pi"
		case OperationNaturalLogarithm:
			return "log(" + process(n.Left) + ")"
		case OperationSquareRoot:
			return "sqrt(" + process(n.Left) + ")"
		case OperationCosine:
			return "cos(" + process(n.Left) + ")"
		case OperationSine:
			return "sin(" + process(n.Left) + ")"
		case OperationTangent:
			return "tan(" + process(n.Left) + ")"
		}
		return ""
	}
	return process(n)
}
