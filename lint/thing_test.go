package lint

import (
	"fmt"
	"testing"
)

func TestObjString(t *testing.T) {
	testData := []struct {
		obj      interface{}
		expected string
	}{
		{
			obj:      VarObj{},
			expected: "var",
		},
		{
			obj:      RangeVarObj{},
			expected: "range var",
		},
		{
			obj:      ConstObj{},
			expected: "const",
		},
		{
			obj:      TypeObj{},
			expected: "type",
		},
		{
			obj:      StructFieldObj{},
			expected: "struct field",
		},
		{
			obj:      FuncObj{},
			expected: "func",
		},
		{
			obj:      MethodObj{},
			expected: "method",
		},
		{
			obj:      ParameterVarObj{ObjKind: isFunc},
			expected: "func parameter",
		},
		{
			obj:      ParameterVarObj{ObjKind: isMethod},
			expected: "method parameter",
		},
		{
			obj:      ParameterVarObj{ObjKind: isInterfaceMethod},
			expected: "interface method parameter",
		},
		{
			obj:      ResultVarObj{ObjKind: isFunc},
			expected: "func result",
		},
		{
			obj:      ResultVarObj{ObjKind: isMethod},
			expected: "method result",
		},
		{
			obj:      ResultVarObj{ObjKind: isInterfaceMethod},
			expected: "interface method result",
		},
	}

	for _, tt := range testData {
		actual := tt.obj.(fmt.Stringer).String()
		if tt.expected != actual {
			t.Errorf("obj: %v, expected: %q, got: %q", tt.obj, tt.expected, actual)
		}
	}
}

func TestVarObjIsVar(t *testing.T) {
	var obj interface{} = VarObj{}
	if _, ok := obj.(Var); !ok {
		t.Errorf("%v is not a Var", obj)
	}
}

func TestRangeVarObjIsVar(t *testing.T) {
	var obj interface{} = RangeVarObj{}
	if _, ok := obj.(Var); !ok {
		t.Error("%v is not a Var", obj)
	}
}

func TestParameterVarObjIsVar(t *testing.T) {
	var obj interface{} = ParameterVarObj{}
	if _, ok := obj.(Var); !ok {
		t.Error("%v is not a Var", obj)
	}
}

func TestResultVarObjIsVar(t *testing.T) {
	var obj interface{} = ResultVarObj{}
	if _, ok := obj.(Var); !ok {
		t.Error("%v is not a Var", obj)
	}
}

func TestParameterVarObjKind(t *testing.T) {
	tableData := []struct {
		name     string
		objKind  ObjKind
		expected [3]bool
	}{
		{
			name:     "isFunc",
			objKind:  isFunc,
			expected: [3]bool{true, false, false},
		},
		{
			name:     "isMethod",
			objKind:  isMethod,
			expected: [3]bool{false, true, false},
		},
		{
			name:     "isInterfaceMethod",
			objKind:  isInterfaceMethod,
			expected: [3]bool{false, false, true},
		},
	}
	for _, tt := range tableData {
		obj := ParameterVarObj{ObjKind: tt.objKind}
		result := [3]bool{
			obj.OfFunc(),
			obj.OfMethod(),
			obj.OfInterfaceMethod(),
		}
		if result != tt.expected {
			t.Error("Test: %s, expected: %v, got: %v", tt.name, tt.expected, result)
		}
	}
}

func TestResultVarObjKind(t *testing.T) {
	tableData := []struct {
		name     string
		objKind  ObjKind
		expected [3]bool
	}{
		{
			name:     "isFunc",
			objKind:  isFunc,
			expected: [3]bool{true, false, false},
		},
		{
			name:     "isMethod",
			objKind:  isMethod,
			expected: [3]bool{false, true, false},
		},
		{
			name:     "isInterfaceMethod",
			objKind:  isInterfaceMethod,
			expected: [3]bool{false, false, true},
		},
	}
	for _, tt := range tableData {
		obj := ResultVarObj{ObjKind: tt.objKind}
		result := [3]bool{
			obj.OfFunc(),
			obj.OfMethod(),
			obj.OfInterfaceMethod(),
		}
		if result != tt.expected {
			t.Error("Test: %s, expected: %v, got: %v", tt.name, tt.expected, result)
		}
	}
}
