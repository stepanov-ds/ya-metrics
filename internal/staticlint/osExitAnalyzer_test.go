package staticlint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsexit(t *testing.T) {
	tests := []struct {
		files map[string]string
		name  string
	}{
		{
			name: "main_with_os_exit",
			files: map[string]string{
				"main/main.go": `
				package main
				
				import "os"
				
				func main() {
					os.Exit(1) // want "function main has os.Exit()"
				}
				`,
			},
		},
		{
			name: "main_without_os_exit",
			files: map[string]string{
				"main/main.go": `
				package main

				func main() {
				    println("ok")
				}
				`,
			},
		},
		{
			name: "not_main_pkg",
			files: map[string]string{
				"main/main.go": `
				package notmain

				import "os"

				func main() {
				    os.Exit(1)
				}
				`,
			},
		},
		{
			name: "other_func",
			files: map[string]string{
				"main/main.go": `
				package main

				import "os"

				func foo() {
				    os.Exit(1)
				}

				func main() {}
				`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanfunc, err := analysistest.WriteFiles(tt.files)
			if err != nil {
				t.Error("Troubles with test, cannot write temporary files: ", err)
			}
			defer cleanfunc()

			analysistest.RunWithSuggestedFixes(t, dir, Analyzer, "main")
		})
	}
}
