package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "thelper detects tests helpers which is not start with t.Helper() method."
const checksDoc = `coma separated list of enabled checks

Available checks

` + checkTBegin + ` - check t.Helper() begins helper function
` + checkTFirst + ` - check *testing.T is first param of helper function
` + checkTName + `  - check *testing.T param has t name

Also available similar checks for benchmark helpers: ` +
	checkBBegin + `, ` + checkBFirst + `, ` + checkBName + `

`

type enabledChecksValue map[string]struct{}

func (m enabledChecksValue) Enabled(c string) bool {
	_, ok := m[c]
	return ok
}

func (m enabledChecksValue) String() string {
	ss := make([]string, 0, len(m))
	for s := range m {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return strings.Join(ss, ",")
}

func (m enabledChecksValue) Set(s string) error {
	ss := strings.FieldsFunc(s, func(c rune) bool { return c == ',' })
	if len(ss) == 0 {
		return nil
	}

	for k := range m {
		delete(m, k)
	}
	for _, v := range ss {
		switch v {
		case checkTBegin, checkTFirst, checkTName, checkBBegin, checkBFirst, checkBName:
			m[v] = struct{}{}
		default:
			return fmt.Errorf("unknown check name %q (see help for full list)", v)
		}
	}
	return nil
}

const (
	checkTBegin = "t_begin"
	checkTFirst = "t_first"
	checkTName  = "t_name"
	checkBBegin = "b_begin"
	checkBFirst = "b_first"
	checkBName  = "b_name"
)

type thelper struct {
	enabledChecks enabledChecksValue
}

// NewAnalyzer return a new thelper analyzer.
// thelper analyzes Go test codes how they use t.Helper() method.
func NewAnalyzer() *analysis.Analyzer {
	thelper := thelper{}
	thelper.enabledChecks = enabledChecksValue{
		checkTBegin: struct{}{},
		checkTFirst: struct{}{},
		checkTName:  struct{}{},
		checkBBegin: struct{}{},
		checkBFirst: struct{}{},
		checkBName:  struct{}{},
	}

	a := &analysis.Analyzer{
		Name: "thelper",
		Doc:  doc,
		Run:  thelper.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}

	a.Flags.Init("nlreturn", flag.ExitOnError)
	a.Flags.Var(&thelper.enabledChecks, "checks", checksDoc)

	return a
}

func (t thelper) run(pass *analysis.Pass) (interface{}, error) {
	tObj := analysisutil.ObjectOf(pass, "testing", "T")
	if tObj == nil {
		return nil, nil
	}

	bObj := analysisutil.ObjectOf(pass, "testing", "B")
	if bObj == nil {
		return nil, nil
	}

	tHelper, _, _ := types.LookupFieldOrMethod(tObj.Type(), true, tObj.Pkg(), "Helper")
	if tHelper == nil {
		return nil, nil
	}

	bHelper, _, _ := types.LookupFieldOrMethod(bObj.Type(), true, bObj.Pkg(), "Helper")
	if bHelper == nil {
		return nil, nil
	}

	tCheckOpts := checkFuncOpts{
		skipPrefix: "Test",
		varName:    "t",
		tbHelper:   tHelper,
		tbType:     types.NewPointer(tObj.Type()),
		checkBegin: t.enabledChecks.Enabled(checkTBegin),
		checkFirst: t.enabledChecks.Enabled(checkTFirst),
		checkName:  t.enabledChecks.Enabled(checkTName),
	}

	bCheckOpts := checkFuncOpts{
		skipPrefix: "Benchmark",
		varName:    "b",
		tbHelper:   bHelper,
		tbType:     types.NewPointer(bObj.Type()),
		checkBegin: t.enabledChecks.Enabled(checkBBegin),
		checkFirst: t.enabledChecks.Enabled(checkBFirst),
		checkName:  t.enabledChecks.Enabled(checkBName),
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		checkFunc(pass, funcDecl, tCheckOpts)
		checkFunc(pass, funcDecl, bCheckOpts)
	})
	return nil, nil
}

type checkFuncOpts struct {
	skipPrefix string
	varName    string
	tbHelper   types.Object
	tbType     *types.Pointer
	checkBegin bool
	checkFirst bool
	checkName  bool
}

func checkFunc(pass *analysis.Pass, funcDecl *ast.FuncDecl, opts checkFuncOpts) {
	if strings.HasPrefix(funcDecl.Name.Name, opts.skipPrefix) {
		return
	}

	p, pos, ok := searchFuncParam(pass, funcDecl, opts.tbType)
	if !ok {
		return
	}

	if opts.checkFirst {
		if pos != 0 {
			pass.Reportf(funcDecl.Pos(), "parameter %s should be the first", opts.tbType)
		}
	}

	if opts.checkName {
		if len(p.Names) > 0 && p.Names[0].Name != opts.varName {
			pass.Reportf(funcDecl.Pos(), "parameter %s should have name %s", opts.tbType, opts.varName)
		}
	}

	if opts.checkBegin {
		if len(funcDecl.Body.List) == 0 || !isTHelperCall(pass, funcDecl.Body.List[0], opts.tbHelper) {
			pass.Reportf(funcDecl.Pos(), "test helper function should start from %s.Helper()", opts.varName)
		}
	}
}

func searchFuncParam(pass *analysis.Pass, f *ast.FuncDecl, p types.Type) (*ast.Field, int, bool) {
	for i, f := range f.Type.Params.List {
		typeInfo, ok := pass.TypesInfo.Types[f.Type]
		if !ok {
			continue
		}

		if types.Identical(typeInfo.Type, p) {
			return f, i, true
		}
	}
	return nil, 0, false
}

func isTHelperCall(pass *analysis.Pass, s ast.Stmt, tHelper types.Object) bool {
	exprStmt, ok := s.(*ast.ExprStmt)
	if !ok {
		return false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	sel, ok := pass.TypesInfo.Selections[selExpr]
	if !ok {
		return false
	}

	return sel.Obj() == tHelper
}
