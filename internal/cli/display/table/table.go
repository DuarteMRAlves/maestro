package table

import (
	"fmt"
	"strings"
)

type Table interface {
	fmt.Stringer
	AddColumn(name string, values []string)
}

type table struct {
	data       []column
	pad        int
	minColSize int
}

func new(b *builder) Table {
	return &table{
		data:       make([]column, 0),
		pad:        b.pad,
		minColSize: b.minColSize,
	}
}

func (t *table) AddColumn(name string, values []string) {
	t.data = append(t.data, newColumn(name, values, t.minColSize))
}

func (t *table) String() string {
	sb := &strings.Builder{}
	t.buildTitle(sb)
	return sb.String()
}

func (t *table) buildTitle(sb *strings.Builder) {
	numCols := len(t.data)
	// Add Padding and last \n
	titleLen := (numCols-1)*t.pad + 1
	for _, c := range t.data {
		titleLen += c.size
	}
	sb.Grow(titleLen)
	for _, c := range t.data {
		pad := c.size - len(c.name)
		sb.WriteString(c.name)
		for i := 0; i < pad; i++ {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
}
