package lint

import (
	"go/ast"
	"go/token"
	"strings"
)

type nodeWalker func(ast.Node) bool

func (w nodeWalker) Walk(node ast.Node) {
	ast.Walk(w, node)
}

func (w nodeWalker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func isTest(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

func WalkNames(fset *token.FileSet, astfile *ast.File, visit func(id *ast.Ident, thing interface{})) {
	visitList := func(fl *ast.FieldList, thing interface{}) {
		if fl == nil {
			return
		}
		for _, f := range fl.List {
			for _, id := range f.Names {
				visit(id, thing)
			}
		}
	}

	fn := func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.AssignStmt:
			// local variable assignment
			if v.Tok == token.ASSIGN {
				return true
			}
			for _, exp := range v.Lhs {
				if id, ok := exp.(*ast.Ident); ok {
					visit(id, VarObj{})
				}
			}

		case *ast.RangeStmt:
			// local variable assignment in range statement
			if v.Tok == token.ASSIGN {
				return true
			}
			if id, ok := v.Key.(*ast.Ident); ok {
				visit(id, RangeVarObj{})
			}
			if id, ok := v.Value.(*ast.Ident); ok {
				visit(id, RangeVarObj{})
			}

		case *ast.FuncDecl:
			// function declaration
			filename := fset.File(astfile.Pos()).Name()
			if isTest(filename) {
				if strings.HasPrefix(v.Name.Name, "Example") ||
					strings.HasPrefix(v.Name.Name, "Test") ||
					strings.HasPrefix(v.Name.Name, "Benchmark") {
					return true
				}
			}

			var kind ObjKind
			if v.Recv == nil {
				kind = isFunc
			} else {
				kind = isMethod
			}

			// global
			visit(v.Name, FuncObj{ObjKind: kind})

			// local
			visitList(v.Type.Params, ParameterVarObj{ObjKind: kind})
			visitList(v.Type.Results, ResultVarObj{ObjKind: kind})

		case *ast.GenDecl:
			// general declaration (global/local)
			if v.Tok == token.IMPORT {
				return true
			}
			var thing interface{}
			switch v.Tok {
			case token.CONST:
				thing = ConstObj{}
			case token.TYPE:
				thing = TypeObj{}
			case token.VAR:
				thing = VarObj{}
			}
			for _, spec := range v.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					visit(s.Name, thing)
				case *ast.ValueSpec:
					for _, id := range s.Names {
						visit(id, thing)
					}
				}
			}

		case *ast.StructType:
			// struct definition (global/local)
			for _, f := range v.Fields.List {
				for _, id := range f.Names {
					visit(id, StructFieldObj{})
				}
			}

		case *ast.InterfaceType:
			// interface definition (global/local)

			// Do not check interface method names.
			// They are often constrainted by the method names of concrete types.
			for _, x := range v.Methods.List {
				ft, ok := x.Type.(*ast.FuncType)
				if !ok {
					// might be an embedded interface name
					continue
				}
				visitList(ft.Params, ParameterVarObj{ObjKind: isInterfaceMethod})
				visitList(ft.Results, ResultVarObj{ObjKind: isInterfaceMethod})
			}

		}

		return true
	}

	nodeWalker(fn).Walk(astfile)
}
