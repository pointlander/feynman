// Copyright 2025 The Feynman Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/big"
)

type Type uint8

const (
	// TypeNumber is a number
	TypeNumber Type = iota
	// TypeNegation negates a number
	TypeNegation
	// TypeAdd adds number together
	TypeAdd
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
	// TypeExpression is an expression
	TypeExpression
)

// Node is a node in an expression
type Node struct {
	Type  Type
	Left  *Node
	Right *Node
	Value *big.Int
}

func (c *Calculator[_]) Eval() *Node {
	return c.Rulee(c.AST())
}

func (c *Calculator[U]) Rulee(node *node[U]) *Node {
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
	minus := false
	for node != nil {
		switch node.pegRule {
		case rulevalue:
			if minus {
				e := &Node{}
				e.Type = TypeNegation
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
			a.Type = TypeNumber
			a.Value = big.NewInt(0)
			a.Value.SetString(string(c.buffer[node.begin:node.end]), 10)
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

func Calculate(n *Node) *big.Int {
	var a *big.Int
	switch n.Type {
	case TypeNumber:
		a = n.Value
	case TypeNegation:
		a = Calculate(n.Left)
		a = a.Neg(a)
	case TypeAdd:
		a = Calculate(n.Left)
		a.Add(a, Calculate(n.Right))
	case TypeSubtract:
		a = Calculate(n.Left)
		a.Sub(a, Calculate(n.Right))
	case TypeMultiply:
		a = Calculate(n.Left)
		a.Mul(a, Calculate(n.Right))
	case TypeDivide:
		a = Calculate(n.Left)
		a.Div(a, Calculate(n.Right))
	case TypeModulus:
		a = Calculate(n.Left)
		a.Mod(a, Calculate(n.Right))
	case TypeExponentiation:
		a = Calculate(n.Left)
		a.Exp(a, Calculate(n.Right), nil)
	case TypeExpression:
		a = Calculate(n.Left)
	}
	return a
}
