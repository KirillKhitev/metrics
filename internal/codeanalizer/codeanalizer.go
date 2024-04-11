// codeanalizer служит для поготовки списка анализаторов для проверки кодовй базы программы.
package codeanalizer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
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
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
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
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

// ExitOnMainAnalyzer - кастомный анализатор, осуществляет проверку, что в функции main пакета main нет вызова os.Exit.
var ExitOnMainAnalyzer = &analysis.Analyzer{
	Name: "exitOnMain",
	Doc:  "Проверка на os.Exit в main()",
	Run:  runExitOnMainAnalyzer,
}

// PrepareChecks формирует список анализаторов.
func PrepareChecks() []*analysis.Analyzer {
	var result []*analysis.Analyzer

	result = addCustomAnalyzers(result)
	result = addAnalyzersSA(result)
	result = addAnalyzersST(result)
	result = addAnalyzersPasses(result)

	return result
}

// addCustomAnalyzers добавляет в список кастомные анализаторы.
func addCustomAnalyzers(result []*analysis.Analyzer) []*analysis.Analyzer {
	result = append(result, ExitOnMainAnalyzer)

	return result
}

// addAnalyzersSA добавляет в список все анализаторы класса SA пакета staticcheck.io.
func addAnalyzersSA(result []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, `SA`) {
			result = append(result, v.Analyzer)
		}
	}

	return result
}

// addAnalyzersST добавляет в список анализатор класса SA пакета simple.io.
func addAnalyzersST(result []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range simple.Analyzers {
		if v.Analyzer.Name == "S1001" {
			result = append(result, v.Analyzer)
			return result
		}
	}

	return result
}

// addAnalyzersPasses добавляет в список все стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes.
func addAnalyzersPasses(result []*analysis.Analyzer) []*analysis.Analyzer {
	result = append(result,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	)

	return result
}

func runExitOnMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return true
				}

				ast.Inspect(node, func(nodeMain ast.Node) bool {
					switch xs := nodeMain.(type) {
					case *ast.CallExpr:
						switch fun := xs.Fun.(type) {
						case *ast.SelectorExpr:
							if fun.Sel.Name == "Exit" {
								pass.Reportf(fun.Pos(), "вызов os.Exit в функции main запрещен")
							}
						}
					}

					return true
				})

				return false
			}

			return true
		})
	}

	return nil, nil
}
