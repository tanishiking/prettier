/*
Package prettier is an implementation of Wadler's classic "A Prettier Printer"
see http://homepages.inf.ed.ac.uk/wadler/papers/prettier/prettier.pdf
*/
package prettier

import (
	"fmt"
)

type document struct {
	col uint
	doc Doc
}

// Pretty renders the given Doc as a string, limiting line lengths to
// `width` or shorter when possible.
//
// Note that this function does not guarantee there are no lines
// longer than `width` -- it just attempts to keep lines within this
// length when possible.
func Pretty(width int, doc Doc) string {
	return best(width, uint(0), doc).layout()
}

func best(width int, k uint, d Doc) chunk {
	return be(
		width,
		k,
		[]*document{
			&document{col: uint(0), doc: d},
		},
	)
}

func be(width int, k uint, x []*document) chunk {
	if len(x) == 0 {
		return &emptyChunk{}
	} else if _, ok := x[0].doc.(*empty); ok {
		return be(width, k, x[1:])
	} else if v, ok := x[0].doc.(*concat); ok {
		i := x[0].col
		return be(
			width,
			k,
			append([]*document{
				&document{col: i, doc: v.a},
				&document{col: i, doc: v.b},
			}, x[1:]...),
		)
	} else if v, ok := x[0].doc.(*text); ok {
		s := v.str
		chunk := be(width, k+uint(v.length), x[1:])
		return &textChunk{
			str: s,
			c:   chunk,
		}
	} else if v, ok := x[0].doc.(*nest); ok {
		i := x[0].col
		indent := v.indent
		return be(
			width,
			k,
			append([]*document{&document{col: i + indent, doc: v.doc}}, x[1:]...),
		)
	} else if _, ok := x[0].doc.(*line); ok {
		i := x[0].col
		chunk := be(width, i, x[1:])
		return &lineChunk{
			indent: i,
			c:      chunk,
		}
	} else if v, ok := x[0].doc.(*union); ok {
		i := x[0].col
		// Since it is redundant to caluculate if the first candidate fits
		// if (w - k) < 0, check if w - k < 0 and if it is true,
		// skip the caluculation and caluculate second candidate.
		if (width - int(k) < 0) {
			second := be(
				width,
				k,
				append([]*document{&document{col: i, doc: v.b}}, x[1:]...),
			)
			return second
		}
		first := be(
			width,
			k,
			append([]*document{&document{col: i, doc: v.a}}, x[1:]...),
		)
		// do not evaluate `v.b` until confirm that first doesn't fits
		// in case `v.b` is lazydoc
		if first.fits(width - int(k)) {
			return first
		}
		second := be(
			width,
			k,
			append([]*document{&document{col: i, doc: v.b}}, x[1:]...),
		)
		return second
	} else if v, ok := x[0].doc.(lazyDoc); ok {
		i := x[0].col
		return be(
			width,
			k,
			append([]*document{&document{col: i, doc: v.Evaluated()}}, x[1:]...),
		)
	} else {
		v := x[0].doc
		panic(fmt.Sprintf("Error: %v sould not be here", v))
	}
}
