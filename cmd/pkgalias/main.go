package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sean9999/pkgalias"
)

func main() {

	env := struct {
		OutSteam  io.Writer
		ErrStream io.Writer
	}{
		os.Stdout,
		os.Stderr,
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(env.ErrStream, `usage: 2 args.
		The first is the package you want to shadow (source).
		The second is the path to the package you want to modify (target).`)
		return
	}

	srcPkgName := os.Args[1]

	targetPkgName := pkgalias.PackageNameFromPath(os.Args[2])

	//	all exported symbols we've already defined
	targetSymbols := pkgalias.Symbols(targetPkgName, os.Args[2])

	//	all exported symbols from the package we wish to shadow
	sourceSymbols := pkgalias.Symbols(srcPkgName, pkgalias.ResolvePath(os.Args[1]))

	//	all elements from src minus target
	diff := pkgalias.Difference(sourceSymbols, targetSymbols)

	//	convert that to go statements
	statements := pkgalias.SymbolsToStatements(srcPkgName, diff)

	//	write to output
	for _, s := range statements {
		fmt.Fprintln(env.OutSteam, s)
	}

}
