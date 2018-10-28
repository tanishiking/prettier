package prettier

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestEmpty(t *testing.T) {
	e := Empty()
	expected := &empty{}
	if !reflect.DeepEqual(e, expected) {
		t.Errorf("expected: %v, actual: %v", expected, e)
	}
}

func TestText(t *testing.T) {
	doc := Text("test")
	expected := &text{str: "test", length: len([]rune("test"))}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestTextWithLength(t *testing.T) {
	doc := TextWithLength("test", 1)
	expected := &text{str: "test", length: 1}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestLine(t *testing.T) {
	doc := Line()
	expected := &line{
		flattenToSpace: true,
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestLineBreak(t *testing.T) {
	doc := LineBreak()
	expected := &line{
		flattenToSpace: false,
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestConcat(t *testing.T) {
	// concat should be right associated
	doc := Concat([]Doc{Text("a"), Text("b"), Text("c"), Text("d")})
	expected := &concat{
		a: Text("a"),
		b: &concat{
			a: Text("b"),
			b: &concat{
				a: Text("c"),
				b: Text("d"),
			},
		},
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestNest(t *testing.T) {
	doc := Nest(uint(4), Text("test"))
	expected := &nest{
		indent: uint(4),
		doc:    Text("test"),
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestBracketBy(t *testing.T) {
	sep := Concat([]Doc{Text(","), LineOrSpace()})
	ds := []Doc{
		Text("foo"),
		Text("bar"),
		Text("baz"),
		Text("hello"),
		Text("world"),
	}
	bracketby := BracketBy(
		Text("["),
		Text("]"),
		Intercalate(sep, ds),
		uint(2),
	)
	actual40 := "[ foo, bar, baz, hello, world ]"
	if Pretty(40, bracketby) != actual40 {
		t.Errorf("expected: %v, actual: %v", actual40, Pretty(40, bracketby))
	}
	actual10 := `[
  foo,
  bar,
  baz,
  hello,
  world
]`
	if Pretty(10, bracketby) != actual10 {
		t.Errorf("expected: %v, actual: %v", actual10, Pretty(10, bracketby))
	}
}

func TestTightBracketBy(t *testing.T) {
	sep := Concat([]Doc{Text(","), LineOrSpace()})
	ds := []Doc{
		Text("foo"),
		Text("bar"),
		Text("baz"),
		Text("hello"),
		Text("world"),
	}
	bracketby := TightBracketBy(
		Text("["),
		Text("]"),
		Intercalate(sep, ds),
		uint(2),
	)
	actual40 := "[foo, bar, baz, hello, world]"
	if Pretty(40, bracketby) != actual40 {
		t.Errorf("expected: %v, actual: %v", actual40, Pretty(40, bracketby))
	}
	actual10 := `[
  foo,
  bar,
  baz,
  hello,
  world
]`
	if Pretty(10, bracketby) != actual10 {
		t.Errorf("expected: %v, actual: %v", actual10, Pretty(10, bracketby))
	}
}

func TestSpaces(t *testing.T) {
	empty := Spaces(uint(0))
	if !reflect.DeepEqual(empty, Empty()) {
		t.Errorf("Spaces should retrn Empty when specified number is 0")
	}

	doc := Spaces(uint(3))
	expected := Text("   ")
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestGroup(t *testing.T) {
	empty := Group(Empty())
	if !reflect.DeepEqual(empty, Empty()) {
		t.Errorf("Group shouldn't construct union if nothing is flattened")
	}

	flatline, _ := Line().flattenBool()
	lineOrSpace := Group(Line())
	lineOrSpaceExpected := &union{
		a: flatline,
		b: Line(),
	}
	if !reflect.DeepEqual(lineOrSpace, lineOrSpaceExpected) {
		t.Errorf("Group should construct union of flatened doc and doc as it is.")
	}
}

func TestGroupLaw(t *testing.T) {
	// group(x) = (x' | x) where x' is flatten(x)
	//
	// (a | b)*c == (a*c | b*c) so, if flatten(c) == c we have:
	// c * (a | b) == (a*c | b*c)
	//
	// b.grouped + flatten(c) == (b + flatten(c)).grouped
	// flatten(c) + b.grouped == (flatten(c) + b).grouped
	cfg := &quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			args[0] = reflect.ValueOf(generateRandomDoc())
			args[1] = reflect.ValueOf(generateRandomDoc())
		},
	}
	f := func(b, c Doc) bool {
		flatc, _ := c.flattenBool()
		if !reflect.DeepEqual(
			Pretty(80, Concat([]Doc{Group(b), flatc})),
			Pretty(80, Group(Concat([]Doc{b, flatc}))),
		) {
			return false
		}
		if !reflect.DeepEqual(
			Pretty(80, Concat([]Doc{flatc, Group(b)})),
			Pretty(80, Group(Concat([]Doc{flatc, b}))),
		) {
			return false
		}
		return true
	}
	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Group violate laws")
	}
}

func TestFill(t *testing.T) {
	ds := []Doc{Text("a"), Text("b"), Text("c")}
	sep := Concat([]Doc{Text(","), Line()})
	if Pretty(10, Fill(sep, ds)) != "a, b, c" {
		t.Errorf("expected: \"%v\", actual: %v", "a, b, c", Pretty(10, Fill(sep, ds)))
	}
	if Pretty(6, Fill(sep, ds)) != "a, b,\nc" {
		t.Errorf("expected: \"%v\", actual: %v", "a, b,\nc", Pretty(5, Fill(sep, ds)))
	}
	if Pretty(0, Fill(sep, ds)) != "a,\nb,\nc" {
		t.Errorf("expected: \"%v\", actual: %v", "a,\nb,\nc", Pretty(0, Fill(sep, ds)))
	}
}

func TestFoldDocs(t *testing.T) {
	ds := []Doc{Text("a"), Text("b"), Text("c")}
	f := func(a Doc, b Doc) Doc {
		return Concat([]Doc{a, Text(","), b})
	}
	doc := FoldDocs(f, ds)
	expected := &concat{
		a: Text("a"),
		b: &concat{
			a: Text(","),
			b: &concat{
				a: Text("b"),
				b: &concat{
					a: Text(","),
					b: Text("c"),
				},
			},
		},
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func TestIntercalate(t *testing.T) {
	ds := []Doc{Text("a"), Text("b"), Text("c")}
	sep := Text(",")
	doc := Intercalate(sep, ds)
	expected := &concat{
		a: Text("a"),
		b: &concat{
			a: Text(","),
			b: &concat{
				a: Text("b"),
				b: &concat{
					a: Text(","),
					b: Text("c"),
				},
			},
		},
	}
	if !reflect.DeepEqual(doc, expected) {
		t.Errorf("expected: %v, actual: %v", expected, doc)
	}
}

func generateRandomDoc() Doc {
	sources := []Doc{
		Empty(),
		Text("test"),
		Line(),
		LineBreak(),
		LineOrSpace(),
		Concat([]Doc{Text("test"), Empty()}),
		Nest(uint(4), Text("test")),
	}
	i := rand.Intn(len(sources))
	return sources[i]
}
