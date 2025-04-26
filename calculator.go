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
	Width = 8
)

// Operation is a mathematical operation
type Operation uint

const (
	// OperationAdd adds two numbers
	OperationAdd Operation = iota
	// OperationSubtract subtracts two numbers
	OperationSubtract
	// OperationMultiply multiplies two numbers
	OperationMultiply
	// OperationDivide divides two numbers
	OperationDivide
	// OperationModulus computes the modulus of two numbers
	OperationModulus
	// OperationExponentiation raises a number to a number
	OperationExponentiation
	// OperationNegate changes the sign of a number
	OperationNegate
	// OperationNumber is a real number
	OperationNumber
	// OperationVariable is a variable
	OperationVariable
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
	// OperationNoop is a noop
	OperationNoop
)

// Node is a node in an expression
type Node struct {
	Operation Operation
	Value     float64
	Variable  string
	Left      *Node
	Right     *Node
}

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

// NewGaussian makes a new gaussian distribution
func NewGaussian() (g G) {
	for i := range g.G {
		g.G[i].Stddev = 1
	}
	return g
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
