package prettier

import (
	"fmt"
)

// Doc represents set of layouts.
type Doc interface {
	String() string
	flattenBool() (Doc, bool)
}

type empty struct{}

func (e *empty) String() string {
	return "Empty()"
}

func (e *empty) flattenBool() (Doc, bool) {
	return e, false
}

type text struct {
	str    string
	length int
}

func (t *text) String() string {
	return fmt.Sprintf("Text(%s)", t.str)
}

func (t *text) flattenBool() (Doc, bool) {
	return t, false
}

type line struct {
	flattenToSpace bool
}

func (l *line) String() string {
	return "Line"
}

func (l *line) flattenBool() (Doc, bool) {
	var flattened Doc
	if l.flattenToSpace {
		flattened = Text(" ")
	} else {
		flattened = &empty{}
	}
	return flattened, true
}

// Align sets the nesting at the current position
// type align struct {
// 	doc Doc
// }
//
// func (a *align) String() string {
// 	return fmt.Sprintf("Align(%v)", a.doc.String())
// }
//
// func (a *align) flattenBool() Doc {
// 	flattened, changed := a.doc.flattenBool()
// 	return (
// 		&align{
//
// 		}
// 	)
// }

type concat struct {
	a Doc
	b Doc
}

func (c *concat) String() string {
	return fmt.Sprintf("Concat(%v, %v)", c.a.String(), c.b.String())
}

func (c *concat) flattenBool() (Doc, bool) {
	flata, changeda := c.a.flattenBool()
	flatb, changedb := c.b.flattenBool()
	return Concat([]Doc{flata, flatb}), (changeda || changedb)
}

type nest struct {
	indent uint
	doc    Doc
}

func (n *nest) String() string {
	return fmt.Sprintf("Nest(%v, %v)", n.indent, n.doc.String())
}

func (n *nest) flattenBool() (Doc, bool) {
	flattened, changed := n.doc.flattenBool()
	return &nest{
		indent: n.indent,
		doc:    flattened,
	}, changed
}

type union struct {
	a Doc
	b Doc
}

func (u *union) String() string {
	return fmt.Sprintf("Union(%v, %v)", u.a.String(), u.b.String())
}

func (u *union) flattenBool() (Doc, bool) {
	// invariant `a.flatten == b.flatten`
	return u.a, true
}

type lazyDoc chan func() Doc

func (l lazyDoc) String() string {
	return "LazyDoc()"
}

func (l lazyDoc) flattenBool() (Doc, bool) {
	f := <-l
	l <- f
	doc := f()
	return doc.flattenBool()
}

func (l lazyDoc) Evaluated() Doc {
	f := <-l
	l <- f
	return f()
}
