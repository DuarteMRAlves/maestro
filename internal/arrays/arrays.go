package arrays

type MapFn[T any, S any] func(T) S

func Map[T any, S any](f MapFn[T, S], elements ...T) []S {
	res := make([]S, 0, len(elements))
	for _, el := range elements {
		res = append(res, f(el))
	}
	return res
}

type FilterFn[T any] func(T) bool

func Filter[T any](f FilterFn[T], elements ...T) []T {
	var res []T
	for _, el := range elements {
		if f(el) {
			res = append(res, el)
		}
	}
	return res
}
