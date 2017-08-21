package lint

import (
	"fmt"
)

type Var interface {
	IsVar()
}

type VarObj struct{}

func (obj VarObj) String() string { return "var" }

func (obj VarObj) IsVar() {}

type RangeVarObj struct {
	VarObj
}

func (obj RangeVarObj) String() string { return "range var" }

type ConstObj struct{}

func (obj ConstObj) String() string { return "const" }

type TypeObj struct{}

func (obj TypeObj) String() string { return "type" }

type StructFieldObj struct{}

func (obj StructFieldObj) String() string { return "struct field" }

type FuncObj struct{}

func (obj FuncObj) String() string { return "func" }

type MethodObj struct{}

func (obj MethodObj) String() string { return "method" }

type ObjKind int

const (
	isFunc = ObjKind(iota + 1)
	isMethod
	isInterfaceMethod
)

func (k ObjKind) String() string {
	switch k {
	case isFunc:
		return "func"
	case isMethod:
		return "method"
	case isInterfaceMethod:
		return "interface method"
	default:
		return fmt.Sprintf("ObjKind(%d)", k)
	}
}

func (k ObjKind) OfFunc() bool            { return k == isFunc }
func (k ObjKind) OfMethod() bool          { return k == isMethod }
func (k ObjKind) OfInterfaceMethod() bool { return k == isInterfaceMethod }

type ParameterVarObj struct {
	VarObj
	ObjKind
}

func (obj ParameterVarObj) String() string {
	return obj.ObjKind.String() + " parameter"
}

type ResultVarObj struct {
	VarObj
	ObjKind
}

func (obj ResultVarObj) String() string {
	return obj.ObjKind.String() + " result"
}

func NewFunctionParameterVarObj() ParameterVarObj {
	return ParameterVarObj{
		ObjKind: isFunc,
	}
}

func NewMethodParameterVarObj() ParameterVarObj {
	return ParameterVarObj{
		ObjKind: isMethod,
	}
}

func NewInterfaceMethodParameterVarObj() ParameterVarObj {
	return ParameterVarObj{
		ObjKind: isInterfaceMethod,
	}
}

func NewFunctionResultVarObj() ResultVarObj {
	return ResultVarObj{
		ObjKind: isFunc,
	}
}

func NewMethodResultVarObj() ResultVarObj {
	return ResultVarObj{
		ObjKind: isMethod,
	}
}

func NewInterfaceMethodResultVarObj() ResultVarObj {
	return ResultVarObj{
		ObjKind: isInterfaceMethod,
	}
}
