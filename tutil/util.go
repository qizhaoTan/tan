package tutil

import (
	"golang.org/x/exp/constraints"
	"strconv"
)

func Itoa[I constraints.Integer](i I) string {
	return strconv.Itoa(int(i))
}

func Atoi[I constraints.Integer](s string) I {
	i, _ := strconv.Atoi(s)
	return I(i)
}
