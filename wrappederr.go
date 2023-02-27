package wrappederr

import (
	"go/ast"
	"go/types"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "wrappederr",
	Doc:  "check for all errors in the return statement where the errors are not wrapped",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilters := []ast.Node{
		(*ast.ReturnStmt)(nil),
	}
	inspect.Preorder(nodeFilters, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ReturnStmt:
			checkReturnStmt(pass, n)
		}
	})
	return nil, nil
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func checkReturnStmt(pass *analysis.Pass, n *ast.ReturnStmt) {
	for _, result := range n.Results {
		data := pass.TypesInfo.Types[result]
		if !types.Implements(data.Type, errType) {
			continue
		}
		switch node := result.(type) {
		case *ast.CallExpr:
			if !isWrappedErr(pass, node) {
				pass.Reportf(node.Pos(), "error is not wrapped")
			}
		case *ast.Ident:
			pass.Reportf(node.Pos(), "error is not wrapped")
		}
	}
}

func isWrappedErr(pass *analysis.Pass, n ast.Expr) bool {
	node, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	if node.Args == nil {
		return false
	}
	fun, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if fun.X == nil {
		pass.Reportf(fun.Pos(), "fun x is nil")
		return false
	}
	switch f := fun.X.(type) {
	case *ast.Ident:
		if f.Name != "errors" {
			return false
		}
		if !slices.Contains([]string{"Wrap", "Wrapf", "New", "Errorf"}, fun.Sel.Name) {
			return false
		}
	}
	return true
}
