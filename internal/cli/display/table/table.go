package table

import (
	"fmt"
	"strings"
)

type Table interface {
	fmt.Stringer
	AddColumn(name string, values []string) error
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

func (t *table) AddColumn(name string, values []string) error {
	if len(t.data) != 0 && len(t.data[0].values) != len(values) {
		return fmt.Errorf(
			"size mismatch between previous columns and new values: "+
				"expected %d rows but received %d",
			len(t.data[0].values),
			len(values))
	}
	t.data = append(t.data, newColumn(name, values, t.minColSize))
	return nil
}

func (t *table) String() string {
	sb := &strings.Builder{}
	t.writeTitle(sb)
	t.writeContent(sb)
	return sb.String()
}

func (t *table) writeTitle(sb *strings.Builder) {
	numCols := len(t.data)
	if numCols == 0 {
		return
	}
	// Add Padding and last \n
	titleLen := (numCols-1)*t.pad + 1
	for _, c := range t.data {
		titleLen += c.size
	}
	sb.Grow(titleLen)
	for _, c := range t.data {
		pad := c.size + t.pad - len(c.name)
		sb.WriteString(c.name)
		for i := 0; i < pad; i++ {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
}

func (t *table) writeContent(sb *strings.Builder) {
	numCols := len(t.data)
	if numCols == 0 {
		return
	}
	numRows := len(t.data[0].values)
	if numRows == 0 {
		return
	}
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		for _, col := range t.data {
			val := col.values[rowIdx]
			pad := col.size + t.pad - len(val)
			sb.WriteString(val)
			for i := 0; i < pad; i++ {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
}
