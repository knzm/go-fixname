package main

import (
	"testing"

	"github.com/knzm/go-fixname/lint"
)

func TestFilterCategory(t *testing.T) {
	testData := []struct {
		name     string
		category lint.Category
		filter   Filter
		expected bool
	}{
		// AllCategory
		{
			name:     "AllCategory should match AllCaps",
			category: lint.AllCaps,
			filter:   Filter{category: AllCategory},
			expected: true,
		},
		{
			name:     "AllCategory should match Underscore",
			category: lint.Underscore,
			filter:   Filter{category: AllCategory},
			expected: true,
		},
		{
			name:     "AllCategory should match General",
			category: lint.General,
			filter:   Filter{category: AllCategory},
			expected: true,
		},
		// Caps
		{
			name:     "Caps should match AllCaps",
			category: lint.AllCaps,
			filter:   Filter{category: Caps},
			expected: true,
		},
		{
			name:     "Caps should not match Underscore",
			category: lint.Underscore,
			filter:   Filter{category: Caps},
			expected: false,
		},
		{
			name:     "Caps should match General",
			category: lint.General,
			filter:   Filter{category: Caps},
			expected: false,
		},
		// Underscore
		{
			name:     "Underscore should not match AllCaps",
			category: lint.AllCaps,
			filter:   Filter{category: Underscore},
			expected: false,
		},
		{
			name:     "Underscore should match Underscore",
			category: lint.Underscore,
			filter:   Filter{category: Underscore},
			expected: true,
		},
		{
			name:     "Underscore should not match General",
			category: lint.General,
			filter:   Filter{category: Underscore},
			expected: false,
		},
		// Caps | Underscore
		{
			name:     "Caps|Underscore should match AllCaps",
			category: lint.AllCaps,
			filter:   Filter{category: Caps | Underscore},
			expected: true,
		},
		{
			name:     "Caps|Underscore should match Underscore",
			category: lint.Underscore,
			filter:   Filter{category: Caps | Underscore},
			expected: true,
		},
		{
			name:     "Caps|Underscore should not match General",
			category: lint.General,
			filter:   Filter{category: Caps | Underscore},
			expected: false,
		},
	}

	for _, tt := range testData {
		actual := tt.filter.byCategory(tt.category)
		if tt.expected != actual {
			t.Errorf("Test: %s, expected: %v, got: %v", tt.name, tt.expected, actual)
		}
	}
}

func TestFilterThing(t *testing.T) {
	testData := []struct {
		name     string
		thing    interface{}
		filter   Filter
		expected bool
	}{
		// AllThing
		{
			name:     "AllThing should match var",
			thing:    lint.VarObj{},
			filter:   Filter{thing: AllThing},
			expected: true,
		},
		{
			name:     "AllThing should match const",
			thing:    lint.ConstObj{},
			filter:   Filter{thing: AllThing},
			expected: true,
		},
		// Const
		{
			name:     "Const should match const",
			thing:    lint.ConstObj{},
			filter:   Filter{thing: Const},
			expected: true,
		},
		{
			name:     "Const should not match var",
			thing:    lint.VarObj{},
			filter:   Filter{thing: Const},
			expected: false,
		},
		// Var
		{
			name:     "Var should match var",
			thing:    lint.VarObj{},
			filter:   Filter{thing: Var},
			expected: true,
		},
		{
			name:     "Var should match range var",
			thing:    lint.RangeVarObj{},
			filter:   Filter{thing: Var},
			expected: true,
		},
		{
			name:     "Var should match function parameter var",
			thing:    lint.NewFunctionParameterVarObj(),
			filter:   Filter{thing: Var},
			expected: true,
		},
		{
			name:     "Var should match method result var",
			thing:    lint.NewMethodResultVarObj(),
			filter:   Filter{thing: Var},
			expected: true,
		},
		{
			name:     "Var should not match const",
			thing:    lint.ConstObj{},
			filter:   Filter{thing: Var},
			expected: false,
		},
		// Type
		{
			name:     "Type should match type",
			thing:    lint.TypeObj{},
			filter:   Filter{thing: Type},
			expected: true,
		},
		// StructField
		{
			name:     "StructField should match struct field",
			thing:    lint.StructFieldObj{},
			filter:   Filter{thing: StructField},
			expected: true,
		},
		// Func
		{
			name:     "Func should match struct func",
			thing:    lint.FuncObj{},
			filter:   Filter{thing: Func},
			expected: true,
		},
	}

	for _, tt := range testData {
		actual := tt.filter.byThing(tt.thing)
		if tt.expected != actual {
			t.Errorf("Test: %s, expected: %v, got: %v", tt.name, tt.expected, actual)
		}
	}
}
