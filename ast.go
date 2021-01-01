package main

type _value struct{}

func (v *_value) valueNode() {}

type Value interface {
	valueNode()
}

type ObjectValue struct {
	_value
	Pair []Pair
}

type Pair struct {
	_value
	Key   *Key
	Value Value
}

type StringValue struct {
	_value
	Token Token
	Value string
}

type NumberValue struct {
	_value
	Token Token
	Value string
}

type Key struct {
	Identifier *Identifier
}

type Identifier struct {
	_value
	Token Token
	Value string
}
