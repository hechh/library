package datatype

import "github.com/hechh/library/redispool"

type HashType0[R any, I any] struct {
	*DataType[R, I]
	key       string
	fieldFunc func() string
	i         I
}
type HashType1[R any, I any, A any] struct {
	*HashType0[R, I]
	fieldFunc func(A) string
	a         A
}
type HashType2[R any, I any, A any, B any] struct {
	*HashType1[R, I, A]
	fieldFunc func(A, B) string
	b         B
}
type HashType3[R any, I any, A any, B any, C any] struct {
	*HashType2[R, I, A, B]
	fieldFunc func(A, B, C) string
	c         C
}

func (d *HashType0[R, I]) GetClient() redispool.IClient { return d.DataType.GetClient(d.i) }
func (d *HashType0[R, I]) GetKey() string               { return d.key }
func (d *HashType0[R, I]) GetField() string             { return d.fieldFunc() }
func (d *HashType1[R, I, A]) GetField() string          { return d.fieldFunc(d.a) }
func (d *HashType2[R, I, A, B]) GetField() string       { return d.fieldFunc(d.a, d.b) }
func (d *HashType3[R, I, A, B, C]) GetField() string    { return d.fieldFunc(d.a, d.b, d.c) }

func H0[R any, I any](h *DataType[R, I], key string, f func() string, i I) *HashType0[R, I] {
	return &HashType0[R, I]{
		DataType:  h,
		key:       key,
		fieldFunc: f,
		i:         i,
	}
}
func H1[R any, I any, A any](h *DataType[R, I], key string, f func(A) string, i I, a A) *HashType1[R, I, A] {
	return &HashType1[R, I, A]{
		HashType0: H0(h, key, nil, i),
		fieldFunc: f,
		a:         a,
	}
}
func H2[R any, I any, A any, B any](h *DataType[R, I], key string, f func(A, B) string, i I, a A, b B) *HashType2[R, I, A, B] {
	return &HashType2[R, I, A, B]{
		HashType1: H1(h, key, nil, i, a),
		fieldFunc: f,
		b:         b,
	}
}
func H3[R any, I any, A any, B any, C any](h *DataType[R, I], key string, f func(A, B, C) string, i I, a A, b B, c C) *HashType3[R, I, A, B, C] {
	return &HashType3[R, I, A, B, C]{
		HashType2: H2(h, key, nil, i, a, b),
		fieldFunc: f,
		c:         c,
	}
}
