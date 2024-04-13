package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

// OsExitAnalyzer checks for direct os.Exit in main function.
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "OsExitAnalyzer",
	Doc:  "check for direct os.Exit in main function",
	Run:  runOsExitMain,
}

func runOsExitMain(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fd, okFD := decl.(*ast.FuncDecl)
			if !okFD || fd.Name.Name != "main" {
				continue
			}

			filePath := pass.Fset.Position(file.Package).Filename
			if strings.Contains(filePath, "go-build") {
				continue
			}

			ast.Inspect(fd, func(n ast.Node) bool {
				ce, okCE := n.(*ast.CallExpr)
				if !okCE {
					return true
				}

				se, okSE := ce.Fun.(*ast.SelectorExpr)
				if !okSE {
					return true
				}

				if ident, ok := se.X.(*ast.Ident); ok && ident.Name == "os" && se.Sel.Name == "Exit" {
					pass.Reportf(ce.Pos(), "avoid using os.Exit directly in main function")
				}

				return true
			})
		}
	}

	return nil, nil
}
