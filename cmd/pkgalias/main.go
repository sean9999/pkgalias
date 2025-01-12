package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sean9999/pkgalias"
	"github.com/spf13/afero"
)

func main() {

	env := struct {
		Filesystem afero.Fs
		Args       []string
		ErrStream  io.Writer
	}{
		afero.NewOsFs(),
		os.Args,
		os.Stderr,
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(env.ErrStream, `usage: 2 args.
		The first is the package you want to shadow (source).
		The second is the path to the package you want to modify (target).`)
		return
	}

	srcPkgName := env.Args[1]

	targetPkgName := pkgalias.PackageNameFromPath(env.Args[2])

	//	all exported symbols we've already defined
	tvars, tfuns, tints := pkgalias.Symbols(targetPkgName, env.Args[2])

	//	all exported symbols from the package we wish to shadow
	svars, sfuns, sints := pkgalias.Symbols(srcPkgName, pkgalias.ResolvePath(env.Args[1]))

	//	all elements from src minus target
	dvars := pkgalias.Difference(svars, tvars)
	dfuns := pkgalias.Difference(sfuns, tfuns)
	dints := pkgalias.Difference(sints, tints)

	//	create or truncate pkgalias.go and open for writing
	filename := fmt.Sprintf("%s/pkgalias.go", env.Args[2])
	_, err := env.Filesystem.Create(filename)
	if err != nil {
		panic(err)
	}
	f, err := env.Filesystem.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	pkgalias.GoCode(f, srcPkgName, targetPkgName, dvars, dfuns, dints)

}
