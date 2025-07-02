// Package analyzerslists provides the lists of analyzers
package analyzerslists

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/gofix"
	"golang.org/x/tools/go/analysis/passes/hostport"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
)

var AnalysisAnalyzers = []*analysis.Analyzer{
	// appends: check for only one variable in append
	appends.Analyzer,
	// asmdecl: check mismatches between assembly files and Go declarations.
	asmdecl.Analyzer,
	// assign: check useless asignments
	assign.Analyzer,
	// atomic: check for common mistakes using the sync/atomic package
	atomic.Analyzer,
	// atomicalign: check for non-64-bits-aligned arguments to sync/atomic function
	atomicalign.Analyzer,
	// bools: check common mistakes involving boolean operators.
	bools.Analyzer,
	// buildssa: build SSA-form IR for later passes
	buildssa.Analyzer,
	// buildtag: check build tags.
	buildtag.Analyzer,
	// cgocall: check some violations of the cgo pointer passing rules.
	cgocall.Analyzer,
	// composite: check for unkeyed composite literals.
	composite.Analyzer,
	// copylock: check for locks erroneously passed by value.
	copylock.Analyzer,
	// ctrlflow: build a control-flow graph
	ctrlflow.Analyzer,
	// deepequalerrors: check for the use of reflect.DeepEqual with error values.
	deepequalerrors.Analyzer,
	// defers: report common mistakes in defer statements
	defers.Analyzer,
	// directive: check known Go toolchain directives.
	directive.Analyzer,
	// errorsas: check that the second argument to errors.As is a pointer to a type implementing error.
	errorsas.Analyzer,
	// fieldalignment: detects structs that would use less memory if their fields were sorted.
	fieldalignment.Analyzer,
	// findcall: serves as a trivial example and test of the Analysis API.
	findcall.Analyzer,
	// framepointer: reports assembly code that clobbers the frame pointer before saving it.
	framepointer.Analyzer,
	// gofix: check go:fix directives.
	gofix.Analyzer,
	// hostport: check format of addresses passed to net.Dial
	hostport.Analyzer,
	// httpmux: report using Go 1.22 enhanced ServeMux patterns in older Go versions.
	httpmux.Analyzer,
	// httpresponse: check for mistakes using HTTP responses.
	httpresponse.Analyzer,
	// ifaceassert: detect impossible interface-to-interface type assertions
	ifaceassert.Analyzer,
	// inspect: optimize AST traversal for later passes
	inspect.Analyzer,
	// loopclosure: inspect statements to find function literals that may be run outside of
	// the current loop iteration.
	loopclosure.Analyzer,
	// lostcancel: check cancel func returned by context.WithCancel is called
	lostcancel.Analyzer,
	// nilfunc: check for useless comparisons between functions and nil
	nilfunc.Analyzer,
	// nilness: check for redundant or impossible nil comparisons
	nilness.Analyzer,
	// pkgfact: gather name/value pairs from constant declarations
	pkgfact.Analyzer,
	// printf: check consistency of Printf format strings and arguments
	printf.Analyzer,
	// reflectvaluecompare: check for comparing reflect.Value values with == or reflect.DeepEqual
	reflectvaluecompare.Analyzer,
	// shadow: check for possible unintended shadowing of variables
	shadow.Analyzer,
	// shift: check for shifts that equal or exceed the width of the integer
	shift.Analyzer,
	// sigchanyzer: check for unbuffered channel of os.Signal
	sigchanyzer.Analyzer,
	// slog: check for invalid structured logging calls
	slog.Analyzer,
	// sortslice: check the argument type of sort.Slice
	sortslice.Analyzer,
	// stdmethods: check signature of methods of well-known interfaces
	stdmethods.Analyzer,
	// stdversion: report uses of too-new standard library symbols
	stdversion.Analyzer,
	// stringintconv: check for string(int) conversions
	stringintconv.Analyzer,
	// structtag: check that struct field tags conform to reflect.StructTag.Get
	structtag.Analyzer,
	// testinggoroutine: report calls to (*testing.T).Fatal from goroutines started by a test
	testinggoroutine.Analyzer,
	// tests: check for common mistaken usages of tests and examples
	tests.Analyzer,
	// timeformat: check for calls of (time.Time).Format or time.Parse with 2006-02-01
	timeformat.Analyzer,
	// unmarshal: report passing non-pointer or non-interface values to unmarshal
	unmarshal.Analyzer,
	// unreachable: check for unreachable code
	unreachable.Analyzer,
	// unsafeptr: check for invalid conversions of uintptr to unsafe.Pointer
	unsafeptr.Analyzer,
	// unusedwrite: reports instances of writes to struct fields and arrays
	// that are never read.
	unusedwrite.Analyzer,
	// usesgenerics: detect whether a package uses generics features
	usesgenerics.Analyzer,
	// waitgroup: check for misuses of sync.WaitGroup
	waitgroup.Analyzer,
}
