# pkgalias

pkgalias allows you to generate an alias of a Go package. This is useful for "shadowing" an entire package.

It generates code of the form:

```go
var SomeFunc = otherpage.SomeFunc
var SomeVar = otherpackage.SomeVar
```

Whichever exported symbols (variables, constants, and functions) you've already defined will not be included.

This let's your users take advantage of the full power of the source package, while benefiting from those exported symbols which you've chosen to modify.

The included binary is meant to be used for code generation. Ex:

```sh
$ go install github.com/sean9999/pkgalias/cmd/pkgalias@latest
```

```go
package mypackage

//  this package shadows "fmt", from the standard library, giving it a Shoutln() function
//  and panicking if you use Println()

//go:generate pkgalias "fmt" .

func Println(..._ any) {
    panic("don't do this")
}

func Shoutln(s string) {
    s = strings.ToUpper(s)
	fmt.Println(s)
}
```
