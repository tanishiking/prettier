package prettier

import (
	"strings"
)

// Empty represents an empty document
func Empty() Doc {
	return &empty{}
}

// Text represents string
// The string must not be empty, and may not contain newlines.
func Text(str string) Doc {
	return &text{
		str: str,
	}
}

// Line represents a single, literal newline.
// which is flattened to a space.
func Line() Doc {
	return &line{flattenToSpace: true}
}

// LineBreak represents a single, literal newline.
// which is flattened to a empty doc.
func LineBreak() Doc {
	return &line{flattenToSpace: false}
}

// LineOrSpace represents a space (if there is enough room),
// or a newline otherwise.
func LineOrSpace() Doc {
	return Group(Line())
}

// Concat concatinate multiple documents.
//
// In practice this method builds documents from the right, so that
// the resulting concatenations are all right-associated.
func Concat(ds []Doc) Doc {
	f := func(a Doc, b Doc) Doc {
		return &concat{a: a, b: b}
	}
	return FoldDocs(f, ds)
}

// Nest represents a "remembered indentation level" for a
// document. Newlines in this document will be followed by at least
// this much indentation (nesting is cumulative).
func Nest(indent uint, doc Doc) Doc {
	return &nest{
		indent: indent,
		doc:    doc,
	}
}

// BracketBy bookend specified Doc between the given Docs.
//
// If the documents (when flattened) all fit on one line, then
// newlines will be collapsed, and spaces will be added,
// and the document will render on one line. If you do not want
// a space, see TightBracketBy
//
// Otherwise, newlines will be used on either side of the document,
// and the requested level of indentation will be added.
func BracketBy(left Doc, right Doc, body Doc, indent uint) Doc {
	return bracketBy(left, right, body, indent, Line())
}

// TightBracketBy bookend specified Doc between the given Docs.
//
// For more details, see BracketBy.
// The difference from BracketBy is that TightBracketBy will
// collapse the documents without a space
// (if flattened docs fit on one line)
func TightBracketBy(left Doc, right Doc, body Doc, indent uint) Doc {
	return bracketBy(left, right, body, indent, LineBreak())
}

func bracketBy(left Doc, right Doc, body Doc, indent uint, ln Doc) Doc {
	return Concat([]Doc{
		left,
		Group(
			Concat([]Doc{
				Nest(
					indent,
					Concat([]Doc{ln, body}),
				),
				ln,
				right,
			}),
		),
	})
}

// Spaces returns multiple spaces.
func Spaces(n uint) Doc {
	if n < 1 {
		return &empty{}
	}
	return &text{
		str: strings.Repeat(" ", int(n)),
	}
}

// Group treats the specified doc as a group that can be compressed.
// The effect of this is to replace newlines with spaces, if there
// is enough room. Otherwise, the Doc will be rendered as-it is.
func Group(doc Doc) Doc {
	flattened, changed := doc.flattenBool()
	if changed {
		return &union{a: flattened, b: doc}
	}
	return flattened
}

// Fill collapse a collection of documents into one document, delimited
// by a specified separator.
func Fill(sep Doc, parts []Doc) Doc {
	flatSep, _ := sep.flattenBool()
	sepGroup := Group(sep)
	if len(parts) == 0 {
		return &empty{}
	} else if len(parts) == 1 {
		return parts[0]
	} else {
		x := parts[0]
		y := parts[1]
		tail := parts[2:]
		flatx, changedx := x.flattenBool()
		flaty, changedy := y.flattenBool()
		if changedx && changedy {
			filling := append([]Doc{flaty}, tail...)
			first := Concat([]Doc{flatx, flatSep, Fill(sep, filling)})
			lazy := lazy(func() Doc {
				return Fill(sep, parts[1:])
			})
			second := Concat([]Doc{x, sep, lazy})
			return &union{a: first, b: second}
		} else if !changedx && changedy {
			filling := append([]Doc{flaty}, tail...)
			first := Concat([]Doc{flatx, flatSep, Fill(sep, filling)})
			lazy := lazy(func() Doc {
				return Fill(sep, parts[1:])
			})
			second := Concat([]Doc{x, sep, lazy})
			return &union{a: first, b: second}
		} else if changedx && !changedy {
			// y == flaty
			filling := append([]Doc{flaty}, tail...)
			filled := Fill(sep, filling)
			first := Concat([]Doc{flatx, flatSep, filled})
			second := Concat([]Doc{x, sep, filled})
			return &union{a: first, b: second}
		} else { // !changedx && !changedy
			// x == flatx
			// y == flaty
			return Concat([]Doc{flatx, sepGroup, Fill(sep, parts[1:])})
		}
	}
}

// FoldDocs combines documents, using the given associative function.
//
// The function `f` must be associative. That is, the expression
// `f(x, f(y, z))` must be equivalent to `f(f(x, y), z)`.
//
// In practice this method builds documents from the right, so that
// the resulting concatenations are all right-associated.
func FoldDocs(f func(a Doc, b Doc) Doc, ds []Doc) Doc {
	if len(ds) == 0 {
		return Empty()
	} else if len(ds) == 1 {
		return ds[0]
	} else {
		x := ds[0]
		xs := ds[1:]
		return f(x, FoldDocs(f, xs))
	}
}

// Intercalate concatenate the given documents together,
// delimited by the given separator.
func Intercalate(sep Doc, ds []Doc) Doc {
	join := []Doc{}
	lastIndex := len(ds) - 1
	for i, d := range ds {
		join = append(join, d)
		if i != lastIndex {
			join = append(join, sep)
		}
	}
	return Concat(join)
}

// lazy creates a Doc which won't be evaluated until needed.
// This is useful in some recursive algorithms.
func lazy(f func() Doc) lazyDoc {
	ch := make(lazyDoc, 1)
	ch <- f
	return ch
}
