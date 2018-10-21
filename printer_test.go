package prettier_test

import (
	"fmt"

	p "github.com/tanishiking/prettier"
)

func Example() {
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

	fmt.Println(p.Pretty(20, doc))
	// Output: [
	//   foo, bar, baz,
	//   hello, world
	// ]
}
