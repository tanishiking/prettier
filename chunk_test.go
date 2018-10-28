package prettier

import (
	"testing"
)

func TestEmptyChunkLayout(t *testing.T) {
	chunk := &emptyChunk{}
	actual := chunk.layout()
	expected := ""
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestEmptyChunkFits(t *testing.T) {
	chunk := &emptyChunk{}
	fit := chunk.fits(1)
	if !fit {
		t.Errorf("emptyChunk should always fit if width is larger than or equal to 0")
	}
	notfit := chunk.fits(-1)
	if notfit {
		t.Errorf("emptyChunk should always not fit if width is smaller than 0")
	}
}

func TestTextChunkLayout(t *testing.T) {
	chunk := &textChunk{
		str:       "foo",
		strLength: 3,
		c: &textChunk{
			str:       "bar",
			c:         &emptyChunk{},
			strLength: 3,
		},
	}
	actual := chunk.layout()
	expected := "foobar"
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestTextChunkFits(t *testing.T) {
	chunk := &textChunk{
		str:       "foo",
		strLength: 3,
		c: &textChunk{
			str:       "bar",
			strLength: 3,
			c:         &emptyChunk{},
		},
	}
	fit := chunk.fits(6)
	if !fit {
		t.Errorf(":")
	}
}

func TestLineChunkLayout(t *testing.T) {
	chunk := &lineChunk{
		indent: uint(2),
		c: &textChunk{
			str:       "bar",
			strLength: 3,
			c:         &emptyChunk{},
		},
	}
	actual := chunk.layout()
	expected := "\n  bar"
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
