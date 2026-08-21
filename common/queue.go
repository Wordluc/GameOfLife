package common

import (
	"slices"
)

type Queue[t any] struct {
	values []t
	index  int
	cmp    func(a, b t) int
}

func NewQueue[t any](values []t, cmp func(a, b t) int) *Queue[t] {
	q := &Queue[t]{
		values: nil,
		cmp:    cmp,
	}
	if values != nil {
		q.Enqueue(values)
	}
	return q
}

func (q *Queue[t]) Len() int {
	return len(q.values)
}

func (q *Queue[t]) Reset() {
	q.index = 0
	q.values = q.values[:0]
}

func (q *Queue[t]) Enqueue(values []t) {
	if q.cmp == nil {
		q.values = append(q.values, values...)
	} else {
		for _, v := range values {
			index, _ := slices.BinarySearchFunc(q.values, v, q.cmp)
			q.values = append(q.values, *new(t))
			copy(q.values[index+1:], q.values[index:])
			q.values[index] = v
		}
	}
}

func (q *Queue[t]) Denqueue() (value t, end bool) {
	if q.index >= len(q.values) {
		return value, true
	}
	value = q.values[q.index]
	q.index++
	return value, false
}

func (q *Queue[t]) GetBack(by int) (value *t) {
	i := q.index - by
	if i >= len(q.values) || i < 0 {
		return
	}
	return new(q.values[i])
}

func (q *Queue[t]) GetLast() (value *t) {
	return new(q.values[len(q.values)-1])
}

func (q *Queue[t]) Clone() *Queue[t] {
	return &Queue[t]{
		values: slices.Clone(q.values),
	}
}
