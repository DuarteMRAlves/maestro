package identifier

import (
	"fmt"
	"math/rand"
	"unicode"
)

const digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Id struct {
	Val string
}

func Rand(size int) (Id, error) {
	if size < 0 {
		return Empty(), fmt.Errorf("size is greater or equal than 0: %d", size)
	}
	if size == 0 {
		return Empty(), nil
	}
	b := make([]byte, size)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return Id{Val: string(b)}, nil
}

func Empty() Id {
	return Id{Val: ""}
}

func (id Id) Clone() Id {
	return Id{Val: id.Val}
}

func (id Id) Size() int {
	return len(id.Val)
}

func (id Id) IsEmpty() bool {
	return id.Val == ""
}

func (id Id) IsValid() bool {
	for _, c := range id.Val {
		if !(unicode.IsUpper(c) || unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}

func (id Id) String() string {
	return fmt.Sprintf("Id{val='%v'}", id.Val)
}
