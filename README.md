# tlipoca9/errors

## Feature

1. 0 依赖
2. 干净且清晰的错误堆栈

## QuickStart

```bash
go get -u github.com/tlipoca9/errors
```

> see https://github.com/tlipoca9/errorsexamples/blob/main/quickstart/main.go

```go
package main

import (
	"fmt"

	"github.com/tlipoca9/errors"
)

func foo(method string) error {
	switch method {
	case "new":
		return errors.New("foo")
	case "wrap":
		return errors.Wrap(foo("new"), "this is wrapeed foo")
	default:
		panic("invalid method")
	}
}

func bar(method string) error {
	switch method {
	case "new":
		return foo("new")
	case "wrap":
		return errors.Wrap(foo("wrap"), "this is wrapeed bar")
	default:
		panic("invalid method")
	}
}

func baz(method string) error {
	switch method {
	case "new":
		return bar("new")
	case "wrap":
		return errors.Wrap(bar("wrap"), "this is wrapeed baz")
	default:
		panic("invalid method")
	}
}

func main() {
	errors.C.Style = errors.StyleStack

	fmt.Println(baz("new"))

	fmt.Println("=====================================")

	fmt.Println(baz("wrap"))
}

/*
foo
  D:/code/projects/errorsexamples/quickstart/main.go:12 (0x457eda) main.foo()
  D:/code/projects/errorsexamples/quickstart/main.go:23 (0x457f9a) main.bar()
  D:/code/projects/errorsexamples/quickstart/main.go:34 (0x45805a) main.baz()
  D:/code/projects/errorsexamples/quickstart/main.go:45 (0x458133) main.main()
  D:/apps/scoop/apps/go/current/src/runtime/proc.go:267 (0x405231) runtime.main()
=====================================
foo
  D:/code/projects/errorsexamples/quickstart/main.go:12 (0x457eda) main.foo()
this is wrapeed foo
  D:/code/projects/errorsexamples/quickstart/main.go:14 (0x457f15) main.foo()
this is wrapeed bar
  D:/code/projects/errorsexamples/quickstart/main.go:25 (0x457fd5) main.bar()
this is wrapeed baz
  D:/code/projects/errorsexamples/quickstart/main.go:36 (0x458095) main.baz()
  D:/code/projects/errorsexamples/quickstart/main.go:49 (0x4581bb) main.main()
  D:/apps/scoop/apps/go/current/src/runtime/proc.go:267 (0x405231) runtime.main()
*/

```

## Examples

> see https://github.com/tlipoca9/errorsexamples/blob/main/stack/main.go

```bash
==== golang errors.New ====
foo
===========================

==== golang fmt.Errorf(%w) ====
this is wrapeed baz: this is wrapeed bar: this is wrapeed foo: foo
============================

==== pkg errors.New ====
foo
=========================

==== pkg errors.Wrap ====
this is wrapeed baz: this is wrapeed bar: this is wrapeed foo: foo
==========================

==== pkg errors.Wrap(stack trace) ====
main.baz
        D:/code/projects/errorsexamples/stack/main.go:64
main.main
        D:/code/projects/errorsexamples/stack/main.go:92
runtime.main
        D:/apps/scoop/apps/go/current/src/runtime/proc.go:267
runtime.goexit
        D:/apps/scoop/apps/go/current/src/runtime/asm_amd64.s:1650
=====================================

==== tlipoca9 errors.New(normal style) ====
foo
===========================================

==== tlipoca9 errors.Wrap(normal style) ====
this is wrapeed baz: this is wrapeed bar: this is wrapeed foo: foo
============================================

==== tlipoca9 errors.New(stack style: default) ====
foo
  D:/code/projects/errorsexamples/stack/main.go:16 (0x1413f6) main.foo()
  D:/code/projects/errorsexamples/stack/main.go:35 (0x141676) main.bar()
  D:/code/projects/errorsexamples/stack/main.go:54 (0x1418b6) main.baz()
  D:/code/projects/errorsexamples/stack/main.go:109 (0x1420bd) main.main()
  D:/apps/scoop/apps/go/current/src/runtime/proc.go:267 (0xd6311) runtime.main()
===================================================

==== tlipoca9 errors.Wrap(stack style: default) ====
foo
  D:/code/projects/errorsexamples/stack/main.go:16 (0x1413f6) main.foo()
this is wrapeed foo
  D:/code/projects/errorsexamples/stack/main.go:18 (0x141435) main.foo()
this is wrapeed bar
  D:/code/projects/errorsexamples/stack/main.go:37 (0x1416b5) main.bar()
this is wrapeed baz
  D:/code/projects/errorsexamples/stack/main.go:56 (0x1418f5) main.baz()
  D:/code/projects/errorsexamples/stack/main.go:113 (0x142187) main.main()
  D:/apps/scoop/apps/go/current/src/runtime/proc.go:267 (0xd6311) runtime.main()
====================================================

==== tlipoca9 errors.New(stack style: json) ====
[{"message":"foo","frames":[{"function":"main.foo","file":"D:/code/projects/errorsexamples/stack/main.go","line":16},{"function":"main.bar","file":"D:/code/projects/errorsexamples/stack/main.go","line":35},{"function":"main.baz","file":"D:/code/projects/errorsexamples/stack/main.go","line":54},{"function":"main.main","file":"D:/code/projects/errorsexamples/stack/main.go","line":118},{"function":"runtime.main","file":"D:/apps/scoop/apps/go/current/src/runtime/proc.go","line":267}]}]
================================================

==== tlipoca9 errors.Wrap(stack style: json) ====
[{"message":"foo","frames":[{"function":"main.foo","file":"D:/code/projects/errorsexamples/stack/main.go","line":16}]},{"message":"this is wrapeed foo","frames":[{"function":"main.foo","file":"D:/code/projects/errorsexamples/stack/main.go","line":18}]},{"message":"this is wrapeed bar","frames":[{"function":"main.bar","file":"D:/code/projects/errorsexamples/stack/main.go","line":37}]},{"message":"this is wrapeed baz","frames":[{"function":"main.baz","file":"D:/code/projects/errorsexamples/stack/main.go","line":56},{"function":"main.main","file":"D:/code/projects/errorsexamples/stack/main.go","line":122},{"function":"runtime.main","file":"D:/apps/scoop/apps/go/current/src/runtime/proc.go","line":267}]}]
=================================================
```
