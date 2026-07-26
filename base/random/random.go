package random

import (
	"math/rand"

	"github.com/hechh/library/base/templ"
)

// Rand 生成一个随机数，范围为 [min, max]
// 整数类型通过 uint64 计算范围避免溢出。
func Rand[T templ.Number](min, max T) T {
	switch any(min).(type) {
	case float64:
		return min + T(rand.Float64())*(max-min)
	case float32:
		return min + T(rand.Float32())*(max-min)
	default:
		return randInt(min, max)
	}
}

// randInt 安全的整数范围随机，使用 uint64 避免 imax-imin+1 溢出
func randInt[T templ.Number](min, max T) T {
	imin, imax := int64(min), int64(max)
	if imax == imin {
		return min
	}
	if imax < imin {
		imin, imax = imax, imin
	}
	span := uint64(imax - imin)
	return T(int64(uint64(imin) + uint64(rand.Int63n(int64(span+1)))))
}

func RandN[T templ.Number](min, max T, num int) (rets []T) {
	tmps := map[T]struct{}{}
	for range num {
		rets = append(rets, randN(min, max, tmps, 10))
	}
	return
}

func randN[T templ.Number](min, max T, tmps map[T]struct{}, times int) T {
	val := Rand(min, max)
	if _, ok := tmps[val]; ok {
		if times <= 0 {
			return val
		}
		return randN(min, max, tmps, times-1)
	}
	tmps[val] = struct{}{}
	return val
}

// Weight 根据权重函数选择一个索引。
func Weight[T templ.Number](size int, f func(int) T) int {
	var total T
	weights := make([]T, 0, size)
	for i := range size {
		w := f(i)
		total += w
		weights = append(weights, w)
	}

	rd := Rand(T(0), total)
	var temp T
	for i := range size {
		temp += weights[i]
		if temp >= rd {
			return i
		}
	}
	return size - 1
}

// WeightN 根据权重函数选择 count 个不重复的索引。
// 如果 count >= size，直接返回所有索引的随机排列。
func WeightN[T templ.Number](size int, count int, f func(int) T) (rets []int) {
	if size <= 0 || count <= 0 {
		return nil
	}
	if count >= size {
		idxs := make([]int, size)
		for i := range size {
			idxs[i] = i
		}
		Shuffle(idxs)
		return idxs
	}

	var sum T
	weights := make([]T, size)
	for i := range size {
		w := f(i)
		sum += w
		weights[i] = w
	}

	tmps := map[int]struct{}{}
	for range count {
		rd := Rand(T(0), sum)
		var temp T
		for j := range size {
			if _, ok := tmps[j]; ok {
				continue
			}
			temp += weights[j]
			if temp >= rd {
				rets = append(rets, j)
				tmps[j] = struct{}{}
				sum -= weights[j]
				break
			}
		}
	}
	return
}

// Shuffle 随机打乱切片
func Shuffle[T any](arr []T) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}
