package util

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

func Or[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Index[T any](arr []T, pos int, def T) T {
	ll := len(arr)
	if pos < 0 || pos >= ll {
		return def
	}
	return arr[pos]
}

func Prefix[T any](arr []T, pos int) []T {
	ll := len(arr)
	if pos < 0 || pos >= ll {
		return arr
	}
	return arr[:pos]
}

func Suffix[T any](arr []T, pos int) []T {
	ll := len(arr)
	if pos < 0 || pos >= ll {
		return arr
	}
	return arr[pos:]
}

func Contains[T any](src []T, dst T) bool {
	for i := range src {
		if reflect.DeepEqual(src[i], dst) {
			return true
		}
	}
	return false
}
