// Package exitonmain содержит анализатор, который запрещает использовать os.Exit() в методе main пакета main.
// Он не содержит никаких флагов. Работает как есть.
package exitonmain

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var ExitOnMainAnalyzer = &analysis.Analyzer{
	Name: "exitonmain",
	Doc:  "check for call os.Exit() on func main() in packet main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	isExit := func(v *ast.CallExpr) bool {
		if IsFunctionNamed(typeutil.StaticCallee(pass.TypesInfo, v), "os", "Exit") {
			return true
		}
		return false
	}

	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		inMain := 0
		ast.Inspect(file, func(node ast.Node) bool {
			if inMain == 0 {
				switch x := node.(type) {
				case *ast.FuncDecl:
					if x.Name.Name == "main" {
						inMain++
					} else {
						return false
					}
				}
				return true
			}
			if node == nil {
				inMain--
				return true
			}
			inMain++
			switch x := node.(type) {
			case *ast.ExprStmt:
				if call, ok := x.X.(*ast.CallExpr); ok {
					if isExit(call) {
						pass.Reportf(call.Pos(), `Call os.Exit on function main of package main`)
					}
				}
			case *ast.DeferStmt:
				if isExit(x.Call) {
					pass.Reportf(x.Call.Pos(), `Call os.Exit on function main of package main`)
				}
			case *ast.GoStmt:
				if isExit(x.Call) {
					pass.Reportf(x.Call.Pos(), `Call os.Exit on function main of package main`)
				}
			}
			return true
		})
	}
	return nil, nil
}

func IsFunctionNamed(f *types.Func, pkgPath string, names ...string) bool {
	if f == nil {
		return false
	}
	if f.Pkg() == nil || f.Pkg().Path() != pkgPath {
		return false
	}
	if f.Type().(*types.Signature).Recv() != nil {
		return false
	}
	for _, n := range names {
		if f.Name() == n {
			return true
		}
	}
	return false
}
