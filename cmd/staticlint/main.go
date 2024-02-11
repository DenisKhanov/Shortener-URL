// Package main is the entry point for the multichecker.
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"honnef.co/go/tools/staticcheck"
)

// main is the entry point for the multichecker.
func main() {
	// mychecks will hold the custom analyzers.
	var mychecks []*analysis.Analyzer

	// Iterate through staticcheck analyzers and add them to mychecks.
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// Run the multichecker with a combination of standard and custom analyzers.
	multichecker.Main(
		// Append standard analyzers to mychecks.
		append([]*analysis.Analyzer{
			buildssa.Analyzer,
			cgocall.Analyzer,
			composite.Analyzer,
			MyAnalyzer, // Custom analyzer.
		}, mychecks...)...,
	)
}

// MyAnalyzer is a custom analyzer that checks for direct os.Exit calls in the main function.
var MyAnalyzer = &analysis.Analyzer{
	Name: "myanalyzer",
	Doc:  "проверяет использование прямого вызова os.Exit в функции main пакета main",
	Run:  run,
}

// run is the main entry point for MyAnalyzer.
func run(pass *analysis.Pass) (interface{}, error) {
	// Пример проверки на прямой вызов os.Exit в функции main
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					if call, ok := n.(*ast.CallExpr); ok {
						if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "os.Exit" {
							pass.Reportf(call.Pos(), "прямой вызов os.Exit в функции main")
						}
					}
					return true
				})
			}
		}
	}
	return nil, nil
}
