package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var ExitCheckAnalyzer = &analysis.Analyzer{ //nolint:exhaustruct
	Name: "exitcheck",
	Doc:  "disallows the direct use of os.Exit in the main function of the main package",
	Run:  runExitCheck,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runExitCheck(pass *analysis.Pass) (any, error) { //revive:disable-line:cyclomatic
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" || fn.Recv != nil {
				continue
			}

			for _, stmt := range fn.Body.List {
				call, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}

				expr, ok := call.X.(*ast.CallExpr)
				if !ok || len(expr.Args) != 1 || !isOsExit(expr.Fun) {
					continue
				}

				pass.Reportf(call.Pos(), "direct use of os.Exit in main function is discouraged")
			}
		}
	}

	return nil, nil
}

func isOsExit(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "Exit" && isOsPackage(ident.Obj)
}

func isOsPackage(obj *ast.Object) bool {
	return obj != nil && obj.Kind == ast.Pkg && obj.Name == "os"
}
