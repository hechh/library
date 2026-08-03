package datatype

import "github.com/hechh/library/redispool"

type StringType0[R any, I any] struct {
	*DataType[R, I]
	keyFunc func() string
	i       I
}
type StringType1[R any, I any, A any] struct {
	*StringType0[R, I]
	keyFunc func(A) string
	a       A
}
type StringType2[R any, I any, A any, B any] struct {
	*StringType1[R, I, A]
	keyFunc func(A, B) string
	b       B
}
type StringType3[R any, I any, A any, B any, C any] struct {
	*StringType2[R, I, A, B]
	keyFunc func(A, B, C) string
	c       C
}

func (d *StringType0[R, I]) GetField() string             { return "" }
func (d *StringType0[R, I]) GetClient() redispool.IClient { return d.DataType.GetClient(d.i) }
func (d *StringType0[R, I]) GetKey() string               { return d.keyFunc() }
func (d *StringType1[R, I, A]) GetKey() string            { return d.keyFunc(d.a) }
func (d *StringType2[R, I, A, B]) GetKey() string         { return d.keyFunc(d.a, d.b) }
func (d *StringType3[R, I, A, B, C]) GetKey() string      { return d.keyFunc(d.a, d.b, d.c) }

func S0[R any, I any](t *DataType[R, I], f func() string, i I) *StringType0[R, I] {
	return &StringType0[R, I]{
		DataType: t,
		keyFunc:  f,
		i:        i,
	}
}
func S1[R any, I any, A any](t *DataType[R, I], f func(A) string, i I, a A) *StringType1[R, I, A] {
	return &StringType1[R, I, A]{
		StringType0: S0(t, nil, i),
		keyFunc:     f,
		a:           a,
	}
}
func S2[R any, I any, A any, B any](t *DataType[R, I], f func(A, B) string, i I, a A, b B) *StringType2[R, I, A, B] {
	return &StringType2[R, I, A, B]{
		StringType1: S1(t, nil, i, a),
		keyFunc:     f,
		b:           b,
	}
}
func S3[R any, I any, A any, B any, C any](t *DataType[R, I], f func(A, B, C) string, i I, a A, b B, c C) *StringType3[R, I, A, B, C] {
	return &StringType3[R, I, A, B, C]{
		StringType2: S2(t, nil, i, a, b),
		keyFunc:     f,
		c:           c,
	}
}
