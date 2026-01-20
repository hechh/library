package util

import "sync"

type Pair[T any, U any] struct {
	First  T
	Second U
}

type Triple[A any, B any, C any] struct {
	First  A
	Second B
	Three  C
}

type Quad[A any, B any, C any, D any] struct {
	First  A
	Second B
	Three  C
	Four   D
}

type Locker[T any] struct {
	mutex sync.RWMutex
	t     T
}

func NewLocker[T any](t T) *Locker[T] {
	return &Locker[T]{t: t}
}

func (d *Locker[T]) Lock() T {
	d.mutex.Lock()
	return d.t
}

func (d *Locker[T]) Unlock() {
	d.mutex.Unlock()
}

func (d *Locker[T]) RLock() T {
	d.mutex.RLock()
	return d.t
}

func (d *Locker[T]) RUnlock() {
	d.mutex.RUnlock()
}

// ----------map1----------------
type Map1[T1 comparable, V any] map[T1]V

func (d Map1[T1, V]) Size() int {
	return len(d)
}

func (d Map1[T1, V]) Put(t1 T1, value V) {
	d[t1] = value
}

func (d Map1[T1, V]) Get(t1 T1) (V, bool) {
	value, ok := d[t1]
	return value, ok
}

func (d Map1[T1, V]) Copy(t1 T1, f func(V) V) (V, bool) {
	value, ok := d[t1]
	if ok {
		value = f(value)
	}
	return value, ok
}

func (d Map1[T1, V]) Del(t1 T1) (V, bool) {
	value, ok := d[t1]
	delete(d, t1)
	return value, ok
}

// ----------map2----------------
type Map2[T1 comparable, T2 comparable, V any] map[Pair[T1, T2]]V

func (d Map2[T1, T2, V]) Size() int {
	return len(d)
}

func (d Map2[T1, T2, V]) Put(t1 T1, t2 T2, value V) {
	d[Pair[T1, T2]{t1, t2}] = value
}

func (d Map2[T1, T2, V]) Get(t1 T1, t2 T2) (V, bool) {
	value, ok := d[Pair[T1, T2]{t1, t2}]
	return value, ok
}

func (d Map2[T1, T2, V]) Copy(t1 T1, t2 T2, f func(V) V) (V, bool) {
	value, ok := d[Pair[T1, T2]{t1, t2}]
	if ok {
		value = f(value)
	}
	return value, ok
}

func (d Map2[T1, T2, V]) Del(t1 T1, t2 T2) (V, bool) {
	key := Pair[T1, T2]{t1, t2}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}

// ----------map3----------------
type Map3[T1 comparable, T2 comparable, T3 comparable, V any] map[Triple[T1, T2, T3]]V

func (d Map3[T1, T2, T3, V]) Size() int {
	return len(d)
}

func (d Map3[T1, T2, T3, V]) Put(t1 T1, t2 T2, t3 T3, value V) {
	d[Triple[T1, T2, T3]{t1, t2, t3}] = value
}

func (d Map3[T1, T2, T3, V]) Get(t1 T1, t2 T2, t3 T3) (V, bool) {
	value, ok := d[Triple[T1, T2, T3]{t1, t2, t3}]
	return value, ok
}

func (d Map3[T1, T2, T3, V]) Copy(t1 T1, t2 T2, t3 T3, f func(V) V) (V, bool) {
	value, ok := d[Triple[T1, T2, T3]{t1, t2, t3}]
	if ok {
		value = f(value)
	}
	return value, ok
}

func (d Map3[T1, T2, T3, V]) Del(t1 T1, t2 T2, t3 T3) (V, bool) {
	key := Triple[T1, T2, T3]{t1, t2, t3}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}

// ----------map4----------------
type Map4[T1 comparable, T2 comparable, T3 comparable, T4 comparable, V any] map[Quad[T1, T2, T3, T4]]V

func (d Map4[T1, T2, T3, T4, V]) Size() int {
	return len(d)
}

func (d Map4[T1, T2, T3, T4, V]) Put(t1 T1, t2 T2, t3 T3, t4 T4, value V) {
	d[Quad[T1, T2, T3, T4]{t1, t2, t3, t4}] = value
}

func (d Map4[T1, T2, T3, T4, V]) Get(t1 T1, t2 T2, t3 T3, t4 T4) (V, bool) {
	value, ok := d[Quad[T1, T2, T3, T4]{t1, t2, t3, t4}]
	return value, ok
}

func (d Map4[T1, T2, T3, T4, V]) Copy(t1 T1, t2 T2, t3 T3, t4 T4, f func(V) V) (V, bool) {
	value, ok := d[Quad[T1, T2, T3, T4]{t1, t2, t3, t4}]
	if ok {
		value = f(value)
	}
	return value, ok
}

func (d Map4[T1, T2, T3, T4, V]) Del(t1 T1, t2 T2, t3 T3, t4 T4) (V, bool) {
	key := Quad[T1, T2, T3, T4]{t1, t2, t3, t4}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}

// ========================group1==============================
type Group1[T1 comparable, V any] map[T1][]V

func NewGroup1[T1 comparable, V any]() Group1[T1, V] {
	return make(map[T1][]V)
}

func (d Group1[T1, V]) Size() int {
	return len(d)
}

func (d Group1[T1, V]) Put(t1 T1, value V) {
	d[t1] = append(d[t1], value)
}

func (d Group1[T1, V]) Get(t1 T1) ([]V, bool) {
	value, ok := d[t1]
	return value, ok
}

func (d Group1[T1, V]) Copy(t1 T1, f func(V) V) (rets []V, ok bool) {
	if values, ok := d[t1]; ok {
		rets = make([]V, len(values))
		for i, item := range values {
			rets[i] = f(item)
		}
	}
	return
}

func (d Group1[T1, V]) Del(t1 T1) ([]V, bool) {
	value, ok := d[t1]
	if ok {
		delete(d, t1)
	}
	return value, ok
}

// ========================group2==============================
type Group2[T1 comparable, T2 comparable, V any] map[Pair[T1, T2]][]V

func NewGroup2[T1 comparable, T2 comparable, V any]() Group2[T1, T2, V] {
	return make(map[Pair[T1, T2]][]V)
}

func (d Group2[T1, T2, V]) Size() int {
	return len(d)
}

func (d Group2[T1, T2, V]) Put(t1 T1, t2 T2, value V) {
	key := Pair[T1, T2]{t1, t2}
	d[key] = append(d[key], value)
}

func (d Group2[T1, T2, V]) Get(t1 T1, t2 T2) ([]V, bool) {
	key := Pair[T1, T2]{t1, t2}
	value, ok := d[key]
	return value, ok
}

func (d Group2[T1, T2, V]) Copy(t1 T1, t2 T2, f func(V) V) (rets []V, ok bool) {
	if values, ok := d[Pair[T1, T2]{t1, t2}]; ok {
		rets = make([]V, len(values))
		for i, item := range values {
			rets[i] = f(item)
		}
	}
	return
}

func (d Group2[T1, T2, V]) Del(t1 T1, t2 T2) ([]V, bool) {
	key := Pair[T1, T2]{t1, t2}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}

// ========================group3==============================
type Group3[T1 comparable, T2 comparable, T3 comparable, V any] map[Triple[T1, T2, T3]][]V

func NewGroup3[T1 comparable, T2 comparable, T3 comparable, V any]() Group3[T1, T2, T3, V] {
	return make(map[Triple[T1, T2, T3]][]V)
}

func (d Group3[T1, T2, T3, V]) Size() int {
	return len(d)
}

func (d Group3[T1, T2, T3, V]) Put(t1 T1, t2 T2, t3 T3, value V) {
	key := Triple[T1, T2, T3]{t1, t2, t3}
	d[key] = append(d[key], value)
}

func (d Group3[T1, T2, T3, V]) Get(t1 T1, t2 T2, t3 T3) ([]V, bool) {
	key := Triple[T1, T2, T3]{t1, t2, t3}
	value, ok := d[key]
	return value, ok
}

func (d Group3[T1, T2, T3, V]) Copy(t1 T1, t2 T2, t3 T3, f func(V) V) (rets []V, ok bool) {
	key := Triple[T1, T2, T3]{t1, t2, t3}
	if values, ok := d[key]; ok {
		rets = make([]V, len(values))
		for i, item := range values {
			rets[i] = f(item)
		}
	}
	return
}

func (d Group3[T1, T2, T3, V]) Del(t1 T1, t2 T2, t3 T3) ([]V, bool) {
	key := Triple[T1, T2, T3]{t1, t2, t3}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}

// ========================group4==============================
type Group4[T1 comparable, T2 comparable, T3 comparable, T4 comparable, V any] map[Quad[T1, T2, T3, T4]][]V

func NewGroup4[T1 comparable, T2 comparable, T3 comparable, T4 comparable, V any]() Group4[T1, T2, T3, T4, V] {
	return make(map[Quad[T1, T2, T3, T4]][]V)
}

func (d Group4[T1, T2, T3, T4, V]) Size() int {
	return len(d)
}

func (d Group4[T1, T2, T3, T4, V]) Put(t1 T1, t2 T2, t3 T3, t4 T4, value V) {
	key := Quad[T1, T2, T3, T4]{t1, t2, t3, t4}
	d[key] = append(d[key], value)
}

func (d Group4[T1, T2, T3, T4, V]) Get(t1 T1, t2 T2, t3 T3, t4 T4) ([]V, bool) {
	key := Quad[T1, T2, T3, T4]{t1, t2, t3, t4}
	value, ok := d[key]
	return value, ok
}

func (d Group4[T1, T2, T3, T4, V]) Copy(t1 T1, t2 T2, t3 T3, t4 T4, f func(V) V) (rets []V, ok bool) {
	key := Quad[T1, T2, T3, T4]{t1, t2, t3, t4}
	if values, ok := d[key]; ok {
		rets = make([]V, len(values))
		for i, item := range values {
			rets[i] = f(item)
		}
	}
	return
}

func (d Group4[T1, T2, T3, T4, V]) Del(t1 T1, t2 T2, t3 T3, t4 T4) ([]V, bool) {
	key := Quad[T1, T2, T3, T4]{t1, t2, t3, t4}
	value, ok := d[key]
	if ok {
		delete(d, key)
	}
	return value, ok
}
