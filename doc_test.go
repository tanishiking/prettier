package prettier

import (
	"reflect"
	"testing"
)

func TestEmptyDoc(t *testing.T) {
	doc := Empty()
	flat, changed := doc.flattenBool()
	if changed || !reflect.DeepEqual(doc, flat) {
		t.Errorf("flatten(Empty()) should not chage the doc")
	}
}

func TestTextDoc(t *testing.T) {
	doc := Text("test")
	flat, changed := doc.flattenBool()
	if changed || !reflect.DeepEqual(doc, flat) {
		t.Errorf("flatten(text) should not chage the doc")
	}
}

func TestLineDoc(t *testing.T) {
	softline := Line()
	soft, changedsoft := softline.flattenBool()
	if !reflect.DeepEqual(soft, Text(" ")) || !changedsoft {
		t.Errorf("Line() should be flattened to a single space")
	}

	linebreak := LineBreak()
	lb, changedlb := linebreak.flattenBool()
	if !reflect.DeepEqual(lb, Empty()) || !changedlb {
		t.Errorf("Line() should be flattened to empty")
	}
}

func TestConcatDoc(t *testing.T) {
	doc := Concat([]Doc{Line(), Text("test")})
	flat, changed := doc.flattenBool()
	if !reflect.DeepEqual(flat, Concat([]Doc{Text(" "), Text("test")})) || !changed {
		t.Errorf("flatten(Concat) should flatten child nodes")
	}
}

func TestNestDoc(t *testing.T) {
	nestedText := Nest(uint(2), Text("test"))
	flatText, textChanged := nestedText.flattenBool()
	if !reflect.DeepEqual(flatText, nestedText) || textChanged {
		t.Errorf("flatten(Nest) should flatten child node")
	}

	nestedLine := Nest(uint(2), Line())
	flatLine, lineChanged := nestedLine.flattenBool()
	if !reflect.DeepEqual(flatLine, Nest(uint(2), Text(" "))) || !lineChanged {
		t.Errorf("flatten(Nest) should flatten child node")
	}
}

func TestUnionDoc(t *testing.T) {
	doc := Group(Line())
	flat, changed := doc.flattenBool()
	flatline, _ := Line().flattenBool()
	if !changed {
		t.Errorf("union.flattenBool should always return true")
	}
	if !reflect.DeepEqual(flat, flatline) {
		t.Errorf("union.flatten should return optimistic rendering")
	}
}

func TestLazyDocEvaluate(t *testing.T) {
	evaluated := false
	f := func() Doc {
		if !evaluated {
			t.Errorf("specified function should not be called until the lazydoc is evaluated")
		}
		return Text("ok")
	}
	lazydoc := lazy(f)
	evaluated = true
	doc := lazydoc.Evaluated()
	if !reflect.DeepEqual(doc, Text("ok")) {
		t.Errorf("lazydoc should return doc when it is evaluated")
	}
}

func TestLazyDocFlatten(t *testing.T) {
	f := func() Doc {
		return Line()
	}
	lazydoc := lazy(f)
	flatlazy, changedlazy := lazydoc.flattenBool()
	flatline, changedline := Line().flattenBool()
	if !reflect.DeepEqual(flatlazy, flatline) || changedlazy != changedline {
		t.Errorf("lazydoc should be flattened to flattened evaluated doc")
	}
}
