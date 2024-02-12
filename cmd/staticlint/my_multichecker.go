// Package main is the entry point for the multichecker.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"honnef.co/go/tools/staticcheck"
)

// main is the entry point for the multichecker.
func main() {
	// myChecks will hold the custom analyzers.
	var myChecks []*analysis.Analyzer

	// Iterate through staticcheck analyzers and add them to myChecks.
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	// Run the multichecker with a combination of standard and custom analyzers.
	multichecker.Main(
		// Append standard analyzers to myChecks.
		append([]*analysis.Analyzer{
			buildssa.Analyzer,
			cgocall.Analyzer,
			composite.Analyzer,
			OsExitAnalyzer, // Custom analyzer.
		}, myChecks...)...,
	)
}
