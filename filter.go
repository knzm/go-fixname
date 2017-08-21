package main

import (
	"github.com/knzm/go-fixname/lint"
)

type categoryBits uint
type thingBits uint

const (
	AllCategory = categoryBits(0)
	Caps        = categoryBits(1 << iota)
	Underscore
)

const (
	AllThing = thingBits(0)
	Const    = thingBits(1 << iota)
	Var
	Type
	StructField
	Func
)

type Filter struct {
	category categoryBits
	thing    thingBits
}

func (f Filter) byCategory(category lint.Category) bool {
	if f.category == AllCategory {
		return true
	}

	switch category {
	case lint.AllCaps:
		// caps
		return f.category&Caps != 0
	case lint.Underscore:
		// underscore
		return f.category&Underscore != 0
	default:
		return false
	}
}

func (f Filter) byThing(thing interface{}) bool {
	if f.thing == AllThing {
		return true
	}

	// const
	if _, ok := thing.(lint.ConstObj); ok {
		return f.thing&Const != 0
	}

	// var
	if _, ok := thing.(lint.Var); ok {
		return f.thing&Var != 0
	}

	// type
	if _, ok := thing.(lint.TypeObj); ok {
		return f.thing&Type != 0
	}

	// struct field
	if _, ok := thing.(lint.StructFieldObj); ok {
		return f.thing&StructField != 0
	}

	// func
	if _, ok := thing.(lint.FuncObj); ok {
		return f.thing&Func != 0
	}

	return false
}
