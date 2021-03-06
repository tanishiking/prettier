# prettier

[![test](https://github.com/tanishiking/prettier/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/tanishiking/prettier/actions/workflows/ci.yml)


## Overview

prettier is an implementation of
[Wadler's "A Prettier Printer"](http://homepages.inf.ed.ac.uk/wadler/papers/prettier/prettier.pdf).

## Install
```sh
$ go get -u github.com/tanishiking/prettier
```

## Usage
```go
import (
    "fmt"

    p "github.com/tanishiking/prettier"
)
j
func main() {
    sep := p.Concat([]p.Doc{p.Text(","), p.LineOrSpace()})
    ds := []p.Doc{
        p.Text("foo"),
        p.Text("bar"),
        p.Text("baz"),
        p.Text("hello"),
        p.Text("world"),
    }
    doc := p.TightBracketBy(
        p.Text("["),
        p.Text("]"),
        p.Intercalate(sep, ds),
        uint(2),
    )

    fmt.Println(p.Pretty(40, doc))
    // [foo, bar, baz, hello, world]

    fmt.Println(p.Pretty(20, doc))
    // [
    //   foo, bar, baz,
    //   hello, world
    // ]

    fmt.Println(p.Pretty(10, doc))
    // [
    //   foo,
    //   bar,
    //   baz,
    //   hello,
    //   world
    // ]
}
```

## License

MIT License
