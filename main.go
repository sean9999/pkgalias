package pkgalias

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"runtime"
	"slices"
	"strings"
)

func Symbols(pkgname, dir string) []string {

	//dir := getDir(pkgname)

	symbols := []string{}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	for pkgName, pkg := range pkgs {

		if strings.HasSuffix(pkgName, "_test") {
			continue
		}

		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.GenDecl:
					if x.Tok == token.CONST || x.Tok == token.VAR {
						for _, spec := range x.Specs {
							vspec := spec.(*ast.ValueSpec)
							for _, name := range vspec.Names {
								if name.IsExported() {
									symbols = append(symbols, name.Name)
								}
							}
						}
					}
				case *ast.FuncDecl:
					if x.Recv == nil && x.Name.IsExported() {
						symbols = append(symbols, x.Name.Name)
					}
				case *ast.TypeSpec:
					if x.Name.IsExported() {
						switch x.Type.(type) {
						case *ast.InterfaceType:
							symbols = append(symbols, x.Name.Name)
						}
					}
				}
				return true
			})
		}
	}

	return symbols
}

func ResolvePath(pkgname string) string {
	//	TODO: panic if dir doesn't exist
	return fmt.Sprintf("%s/src/%s", runtime.GOROOT(), pkgname)
}

func SymbolsToStatements(pkg string, symbols []string) []string {
	statements := make([]string, len(symbols))
	for i, s := range symbols {
		line := fmt.Sprintf("var %s = %s.%s", s, pkg, s)
		statements[i] = line
	}
	return statements
}

func Output(w io.Writer, pkgname string) {
	dir := ResolvePath(pkgname)
	symbols := Symbols(pkgname, dir)
	for _, line := range SymbolsToStatements(pkgname, symbols) {
		fmt.Fprintln(w, line)
	}

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

func PackageNameFromPath(dir string) string {

	fset := token.NewFileSet()
	pkgmap, _ := parser.ParseDir(fset, dir, noTestFiles, parser.AllErrors)
	for name := range pkgmap {
		return name
	}
	return ""

}

func Difference(src, dest []string) []string {

	final := make([]string, 0, len(src))

	for _, str := range src {
		if !slices.Contains(dest, str) {
			final = append(final, str)
		}
	}
	return final
}
