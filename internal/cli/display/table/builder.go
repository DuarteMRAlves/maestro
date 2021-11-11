package table

const defaultPad = 2

type Builder interface {
	WithPadding(pad int) Builder
	WithMinColSize(size int) Builder
	Build() Table
}

type builder struct {
	pad        int
	minColSize int
}

func NewBuilder() Builder {
	return &builder{
		pad: defaultPad,
	}
}

func (b *builder) WithPadding(pad int) Builder {
	b.pad = pad
	return b
}

func (b *builder) WithMinColSize(size int) Builder {
	b.minColSize = size
	return b
}

func (b *builder) Build() Table {
	return new(b)
}
