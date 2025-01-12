# pkgalias

pkgalias allows you to generate an alias of a Go package. This is useful for "shadowing" an entire package. It generates code of the form:

```go
//  func
var SomeFunc = otherpage.SomeFunc

//  var
var SomeVar = otherpackage.SomeVar

//  interface
type SomeInterface interface {
    otherpackage.SomeInterface
}
```

The general use-case for this tool is to support packages that want to act as a drop-in replacements for other packages. The target package need only rewrite the necesary portion to suit the purpose. The rest can be handled by pkgalias.

Whichever exported symbols you've already defined in your target package will not be included. This lets your users take advantage of the full power of the source package, while benefiting from those exported symbols which you've chosen to modify.

The included binary is meant to be used in `//go:generate` style code generation. Ex:

```sh
$ go install github.com/sean9999/pkgalias/cmd/pkgalias@latest
```

```go
package mypackage

/**
 *  This package shadows "fmt", from the standard library, giving it a Shoutln() function.
 *  It panics if you use Println().
 **/

//go:generate pkgalias "fmt" .

func Println(..._ any) {
    panic("don't do this")
}

func Shoutln(s string) {
    s = strings.ToUpper(s)
	fmt.Println(s)
}
```

After building, you'll see a file called `pkgalias.go`. It will import and re-export everything from "fmt" you haven't specifically re-written. In our example, it will include everything but `Println()`.

```go
//  Code generated by pkgalias. DO NOT EDIT.
//  exporting symbols from fmt
package fmt

import (
	"fmt"
)

//	functions
var FormatString = fmt.FormatString
var Fprintf = fmt.Fprintf
var Sprintf = fmt.Sprintf
var Appendf = fmt.Appendf
var Fprint = fmt.Fprint
var Print = fmt.Print
var Sprint = fmt.Sprint
var Append = fmt.Append
var Fprintln = fmt.Fprintln
var Sprintln = fmt.Sprintln
var Appendln = fmt.Appendln
var Scan = fmt.Scan
var Scanln = fmt.Scanln
var Scanf = fmt.Scanf
var Sscan = fmt.Sscan
var Sscanln = fmt.Sscanln
var Sscanf = fmt.Sscanf
var Fscan = fmt.Fscan
var Fscanln = fmt.Fscanln
var Fscanf = fmt.Fscanf

//	interfaces
type State interface {
    fmt.State
}
type Formatter interface {
    fmt.Formatter
}
type Stringer interface {
    fmt.Stringer
}
type GoStringer interface {
    fmt.GoStringer
}
type ScanState interface {
    fmt.ScanState
}
type Scanner interface {
    fmt.Scanner
}
```