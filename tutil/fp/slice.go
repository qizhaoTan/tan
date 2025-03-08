package fp

import (
	"golang.org/x/exp/constraints"
)

func All[S ~[]E, E any](s S, f func(e E) bool) bool {
	for _, e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

func Either[S ~[]E, E any](s S, f func(e E) bool) bool {
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

func Map[S ~[]E, E, F any](s S, f func(e E) F) []F {
	ret := make([]F, 0, len(s))
	for _, e := range s {
		ret = append(ret, f(e))
	}
	return ret
}

func Reduce[S ~[]E, E any, I constraints.Integer | constraints.Float](s S, f func(e E) I) I {
	sum := I(0)
	for _, e := range s {
		sum += f(e)
	}
	return sum
}

func Filter[S ~[]E, E any, F ~func(e E) bool](s S, filter F) (ret S) {
	for _, e := range s {
		if filter(e) {
			ret = append(ret, e)
		}
	}
	return ret
}
