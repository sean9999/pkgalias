package pkgalias

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

//go:embed alias.tmpl
var t string

// Symbols returns all exported symbols from "pkgname" by scanning the directory "dir"
func Symbols(pkgname, dir string) (variables, functions, interfaces []string) {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, noTestFiles, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	pkg := pkgs[pkgname]

	for _, file := range pkg.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.GenDecl:
				if x.Tok == token.CONST || x.Tok == token.VAR {
					for _, spec := range x.Specs {
						vspec := spec.(*ast.ValueSpec)
						for _, name := range vspec.Names {
							if name.IsExported() {
								variables = append(variables, name.Name)
							}
						}
					}
				}
			case *ast.FuncDecl:
				if x.Recv == nil && x.Name.IsExported() {
					functions = append(functions, x.Name.Name)
				}
			case *ast.TypeSpec:
				if x.Name.IsExported() {
					switch x.Type.(type) {
					case *ast.InterfaceType:
						interfaces = append(interfaces, x.Name.Name)
					}
				}
			}
			return true
		})
	}

	return variables, functions, interfaces
}

// ResolvePath resolves a package name to the location on disk where it's source code can be found
func ResolvePath(pkgname string) string {

	gopath, _ := os.LookupEnv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	pkgDir := fmt.Sprintf("%s/src/%s", runtime.GOROOT(), pkgname)
	f, err := os.Stat(pkgDir)
	if err == nil && f.IsDir() {
		return pkgDir
	}
	pkgDir = fmt.Sprintf("%s/pkg/mod/%s", gopath, pkgname)
	f, err = os.Stat(pkgDir)
	if err == nil && f.IsDir() {
		return pkgDir
	}
	panic(fmt.Sprintf("could not find a path for %q", pkgname))

}

// SymbolsToStatements produces syntactically valid go statements
func SymbolsToStatements(pkg string, symbols []string) []string {
	statements := make([]string, len(symbols))
	for i, s := range symbols {
		line := fmt.Sprintf("var %s = %s.%s", s, pkg, s)
		statements[i] = line
	}
	return statements
}

func noTestFiles(f fs.FileInfo) bool {
	if f.IsDir() {
		return false
	}
	if f.Name() == "" {
		return false
	}
	if strings.HasSuffix(f.Name(), "_test.go") {
		return false
	}
	return strings.HasSuffix(f.Name(), ".go")
}

// PackageNameFromPath takes a path to a directory and if that directory represents a go package, returns it's name.
// If there are multiple packages, it returns the first one found.
// It ignores test packages.
func PackageNameFromPath(dir string) string {

	fset := token.NewFileSet()
	pkgmap, _ := parser.ParseDir(fset, dir, noTestFiles, parser.AllErrors)
	for name := range pkgmap {
		return name
	}
	return ""

}

// Difference returns those elements from dest which are not present in src.
func Difference(src, dest []string) []string {

	final := make([]string, 0, len(src))
	for _, str := range src {
		if !slices.Contains(dest, str) {
			final = append(final, str)
		}
	}
	return final
}

func GoCode(f afero.File, srcpkg, destpkg string, vars, funcs, interfaces []string) {

	data := struct {
		Src        string
		Dest       string
		Vars       []string
		Funcs      []string
		Interfaces []string
	}{
		srcpkg, destpkg, vars, funcs, interfaces,
	}

	tmpl, err := template.New("pkgalias").Parse(t)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		panic(err)
	}

}
