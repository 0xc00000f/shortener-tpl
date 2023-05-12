/*
Current package is staticcheck, which is a collection of analyzers.
Analyzers which used in this project are:

asmdecl - check for assembly files that don't declare any functions;
assign - check for useless assignments;
analyzers.ExitCheckAnalyzer - check for os.Exit calls in main package;

Some analyzers from simple package: S1005, S1006.
And all analyzers from staticcheck package.

For more information about staticcheck analyzers, see https://staticcheck.io/docs/checks.

To run this linter, use the following command:
go run ./cmd/staticlint/main.go ./...
*/
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"

	"github.com/0xc00000f/shortener-tpl/cmd/staticlint/analyzers"
)

func main() {
	a := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		analyzers.ExitCheckAnalyzer,
	}

	for _, v := range staticcheck.Analyzers {
		a = append(a, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		if v.Analyzer.Name == "S1005" || v.Analyzer.Name == "S1006" {
			a = append(a, v.Analyzer)
		}
	}

	multichecker.Main(
		a...,
	)
}
