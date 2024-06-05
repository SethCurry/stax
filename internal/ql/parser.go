package ql

import "github.com/SethCurry/stax/internal/bones/predicate"

type Leaf interface {
	Predicate() predicate.Card
}

type Node interface {
	Leaf
	Left() Leaf
	Right() Leaf
}

type LogicNode struct {
	predicator func(...predicate.Card) predicate.Card
	left       Leaf
	right      Leaf
}

func (l *LogicNode) Predicate() predicate.Card {
	return l.predicator(l.left.Predicate(), l.right.Predicate())
}

func (l *LogicNode) Left() Leaf {
	return l.left
}

func (l *LogicNode) Right() Leaf {
	return l.right
}
