package prettier

import (
	"fmt"
	"strings"
)

type chunk interface {
	layout() string
	fits(width int) bool
	String() string
}

type emptyChunk struct{}

func (e *emptyChunk) layout() string {
	return ""
}

func (e *emptyChunk) fits(width int) bool {
	return width >= 0
}

func (e *emptyChunk) String() string {
	return "Empty"
}

type textChunk struct {
	str       string
	strLength int
	c         chunk
}

func (t *textChunk) layout() string {
	return t.str + t.c.layout()
}

func (t *textChunk) fits(width int) bool {
	if width >= 0 {
		return t.c.fits(width - t.strLength)
	}
	return false
}

func (t *textChunk) String() string {
	return fmt.Sprintf("TextChunk(%v, %v)", t.str, t.c.String())
}

type lineChunk struct {
	indent uint
	c      chunk
}

func (l *lineChunk) layout() string {
	return "\n" + strings.Repeat(" ", int(l.indent)) + l.c.layout()
}

func (l *lineChunk) fits(width int) bool {
	return width >= 0
}

func (l *lineChunk) String() string {
	return fmt.Sprintf("LineChunk(%v, %v)", l.indent, l.c.String())
}
