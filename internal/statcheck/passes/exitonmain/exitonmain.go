// Package exitonmain содержит анализатор, который запрещает использовать os.Exit() в методе main пакета main.
// Он не содержит никаких флагов. Работает как есть.
package exitonmain

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// ExitOnMainAnalyzer аналайзер
var ExitOnMainAnalyzer = &analysis.Analyzer{
	Name: "exitonmain",
	Doc:  "check for call os.Exit() on func main() in packet main",
	Run:  run,
}

// run запуск аналайзера
func run(pass *analysis.Pass) (any, error) {
	isExit := func(v *ast.CallExpr) bool {
		return IsFunctionNamed(typeutil.StaticCallee(pass.TypesInfo, v), "os", "Exit")
	}

	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true
			}
			ast.Inspect(fn.Body, func(bodyNode ast.Node) bool {
				if callExpr, ok := bodyNode.(*ast.CallExpr); ok && isExit(callExpr) {
					pass.Reportf(callExpr.Pos(), `Call os.Exit on function main of package main`)
				}
				return true
			})

			return true
		})
	}
	return nil, nil
}

// IsFunctionNamed проверка имени функции
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
