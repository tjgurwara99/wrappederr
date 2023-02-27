package wrappederr

import (
	"go/ast"
	"go/types"

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
		(*ast.FuncDecl)(nil),
	}
	inspect.Preorder(nodeFilters, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			ok, pos := funcReturnsErr(pass, n)
			if !ok {
				return
			}
			checkFunc(pass, n, pos)
		}
	})
	return nil, nil
}

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func funcReturnsErr(pass *analysis.Pass, n *ast.FuncDecl) (bool, int) {
	funcType := pass.TypesInfo.TypeOf(n.Name)
	signature, _ := funcType.(*types.Signature)
	if n.Body == nil || signature == nil || signature.Results().Len() == 0 {
		return false, 0
	}

	pos := 0
	for pos < signature.Results().Len() {
		if types.Implements(signature.Results().At(pos).Type(), errType) {
			return true, pos
		}
		pos++
	}
	return false, 0
}

func checkFunc(pass *analysis.Pass, n *ast.FuncDecl, pos int) {
	if n.Body == nil {
		return
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilters := []ast.Node{
		(*ast.ReturnStmt)(nil),
	}
	inspect.Preorder(nodeFilters, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ReturnStmt:
			checkReturnStmt(pass, n, pos)
		}
	})
}

func checkReturnStmt(pass *analysis.Pass, n *ast.ReturnStmt, pos int) {
	if len(n.Results) <= pos {
		return
	}
	data := pass.TypesInfo.Types[n.Results[pos]]
	if !types.Implements(data.Type, errType) {
		return
	}
	switch node := n.Results[pos].(type) {
	case *ast.CallExpr:
		if !isWrappedErr(pass, node) {
			pass.Reportf(node.Pos(), "error is not wrapped")
			return
		}
	case *ast.Ident:
		pass.Reportf(node.Pos(), "error is not wrapped")
	}
}

func isWrappedErr(pass *analysis.Pass, n *ast.CallExpr) bool {
	if n.Args == nil {
		return false
	}
	fun, ok := n.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if fun.Sel.Name != "Wrap" {
		return false
	}
	if fun.X.(*ast.Ident).Name != "errors" {
		return false
	}
	return true
}
