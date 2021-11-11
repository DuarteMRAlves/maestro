package table

type column struct {
	name   string
	values []string
	size   int
}

func newColumn(name string, values []string, minColSize int) column {
	vCpy := make([]string, 0, len(values))
	size := minColSize
	if len(name) > size {
		size = len(name)
	}
	for _, v := range values {
		vCpy = append(vCpy, v)
		if len(v) > size {
			size = len(v)
		}
	}
	return column{name: name, values: vCpy, size: size}
}
