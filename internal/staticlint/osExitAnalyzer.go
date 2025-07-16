// Package staticlint implements logics for custom multichecker.
package staticlint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer defines the osexit analyzer.
// It checks for direct calls to os.Exit in the main function of the main package.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "Checks os.Exit calls in main function",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		pos := pass.Fset.Position(funcDecl.Pos())
		filename := pos.Filename

		if strings.Contains(filename, "go-build") {
			return
		}

		if funcDecl.Name.Name == "main" && pass.Pkg.Name() == "main" {
			ast.Inspect(funcDecl, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := selExpr.X.(*ast.Ident); ok {
						if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
							pass.Reportf(callExpr.Lparen, "function main has os.Exit() call")
						}
					}
				}
				return true
			})
		}
	})

	return nil, nil
}
