// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

type Type uint8

const (
	// TypeAdd adds number together
	TypeAdd Type = iota
	// TypeSubtract subtracts one number from another
	TypeSubtract
	// TypeMultiply multiplies two numbers
	TypeMultiply
	// TypeDivide divides two numbers
	TypeDivide
	// TypeModules computes the modulus
	TypeModulus
	// TypeExponentiation raises a number to a power
	TypeExponentiation
	// TypeNegation negates a number
	TypeNegation
	// TypeExpression is an expression
	TypeExpression
	// TypeNumber is a number
	TypeNumber
	// TypeVariable is a variable
	TypeVariable
)

// Symbol is a symbol
type Symbol struct {
	Symbol string
	Type   int
}

var Symbols = [...]Symbol{
	{"+", 0},
	{"-", 0},
	{"*", 0},
	{"/", 0},
	{"%", 0},
	{"^", 0},
	{"-", 1},
	{"()", 2},
}

// Generate generates an equation
func Generate(rng *rand.Rand) string {
	x := rng.Perm(3)
	y := rng.Perm(8)
	for _, v := range x[:2] {
		switch v {
		case 0:
			if rng.NormFloat64() > 0 {
				return "x"
			}
		case 1:
			x := 0
			for rng.NormFloat64() > 0 {
				x++
			}
			return fmt.Sprintf("%d", x)
		case 2:
			for _, vv := range y[:7] {
				if Symbols[vv].Type == 0 {
					if rng.NormFloat64() > 0 {
						return Generate(rng) + Symbols[vv].Symbol + Generate(rng)
					}
				} else if Symbols[vv].Type == 1 {
					if rng.NormFloat64() > 0 {
						return Symbols[vv].Symbol + Generate(rng)
					}
				} else {
					if rng.NormFloat64() > 0 {
						return "(" + Generate(rng) + ")"
					}
				}
			}
		}
	}

	switch x[2] {
	case 0:
		if rng.NormFloat64() > 0 {
			return "x"
		}
	case 1:
		x := 0
		for rng.NormFloat64() > 0 {
			x++
		}
		return fmt.Sprintf("%d", x)
	case 2:
		vv := y[7]
		if Symbols[vv].Type == 0 {
			if rng.NormFloat64() > 0 {
				return Generate(rng) + Symbols[vv].Symbol + Generate(rng)
			}
		} else if Symbols[vv].Type == 1 {
			if rng.NormFloat64() > 0 {
				return Symbols[vv].Symbol + Generate(rng)
			}
		} else {
			if rng.NormFloat64() > 0 {
				return "(" + Generate(rng) + ")"
			}
		}
	}
	return "0"
}

// Node is a node in an expression
type Node struct {
	Type     Type
	Left     *Node
	Right    *Node
	Value    *big.Int
	Variable string
	Count    int
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
			b.Type = TypeAdd
			b.Left = a
			b.Right = c.Rulee2(node)
			a = b
		case ruleminus:
			node = node.next
			b := &Node{}
			b.Type = TypeSubtract
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
			b.Type = TypeMultiply
			b.Left = a
			b.Right = c.Rulee3(node)
			a = b
		case ruledivide:
			node = node.next
			b := &Node{}
			b.Type = TypeDivide
			b.Left = a
			b.Right = c.Rulee3(node)
			a = b
		case rulemodulus:
			node = node.next
			b := &Node{}
			b.Type = TypeModulus
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
			b.Type = TypeExponentiation
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
	minus := 0
	for node != nil {
		switch node.pegRule {
		case rulevalue:
			if minus > 0 {
				e := &Node{}
				e.Type = TypeNegation
				e.Count = minus
				e.Left = c.Rulevalue(node)
				return e
			}
			return c.Rulevalue(node)
		case ruleminus:
			minus++
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
			a.Type = TypeNumber
			a.Value = big.NewInt(0)
			a.Value.SetString(string(c.buffer[node.begin:node.end]), 10)
			return a
		case rulevariable:
			a := &Node{}
			a.Type = TypeVariable
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
			e := &Node{
				Type: TypeExpression,
				Left: c.Rulee1(node),
			}
			return e
		}
		node = node.next
	}
	return nil
}

func (n *Node) Calculate() *big.Int {
	var a *big.Int
	switch n.Type {
	case TypeNumber:
		a = n.Value
	case TypeNegation:
		a = n.Left.Calculate()
		a = a.Neg(a)
	case TypeAdd:
		a = n.Left.Calculate()
		a.Add(a, n.Right.Calculate())
	case TypeSubtract:
		a = n.Left.Calculate()
		a.Sub(a, n.Right.Calculate())
	case TypeMultiply:
		a = n.Left.Calculate()
		a.Mul(a, n.Right.Calculate())
	case TypeDivide:
		a = n.Left.Calculate()
		a.Div(a, n.Right.Calculate())
	case TypeModulus:
		a = n.Left.Calculate()
		a.Mod(a, n.Right.Calculate())
	case TypeExponentiation:
		a = n.Left.Calculate()
		a.Exp(a, n.Right.Calculate(), nil)
	case TypeExpression:
		a = n.Left.Calculate()
	}
	return a
}

func (n *Node) String() string {
	var a string
	switch n.Type {
	case TypeNumber:
		a = n.Value.String()
	case TypeVariable:
		a = n.Variable
	case TypeNegation:
		a = n.Left.String()
		minus := ""
		for range n.Count {
			minus += "-"
		}
		a = minus + a
	case TypeAdd:
		a = n.Left.String()
		a = a + "+" + n.Right.String()
	case TypeSubtract:
		a = n.Left.String()
		a = a + "-" + n.Right.String()
	case TypeMultiply:
		a = n.Left.String()
		a = a + "*" + n.Right.String()
	case TypeDivide:
		a = n.Left.String()
		a = a + "/" + n.Right.String()
	case TypeModulus:
		a = n.Left.String()
		a = a + "%" + n.Right.String()
	case TypeExponentiation:
		a = n.Left.String()
		a = a + "^" + n.Right.String()
	case TypeExpression:
		a = "(" + n.Left.String() + ")"
	}
	return a
}
