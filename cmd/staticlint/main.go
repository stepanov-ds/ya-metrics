// Just use
// go run ./cmd/staticlint ./...
package main

import (
	"github.com/alexkohler/nakedret/v2"
	"github.com/stepanov-ds/ya-metrics/internal/staticlint"
	"github.com/stepanov-ds/ya-metrics/internal/staticlint/analyzersLists"
	"github.com/ultraware/funlen"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var checks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if analyzerslists.StaticcheckAnalyzers[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range simple.Analyzers {
		if analyzerslists.StaticcheckAnalyzers[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range stylecheck.Analyzers {
		if analyzerslists.StaticcheckAnalyzers[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range quickfix.Analyzers {
		if analyzerslists.StaticcheckAnalyzers[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}
	checks = append(checks, analyzerslists.AnalysisAnalyzers...)
	checks = append(checks, nakedret.NakedReturnAnalyzer(&nakedret.NakedReturnRunner{MaxLength: 5, SkipTestFiles: false}))
	checks = append(checks, funlen.NewAnalyzer(250, 40, false))
	checks = append(checks, staticlint.Analyzer)

	multichecker.Main(
		checks...,
	)
}
