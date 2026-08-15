package common

import (
	"slices"
)

type Queue[t any] struct {
	values []t
	index  int
}

func NewQueue[t any](values []t) Queue[t] {
	return Queue[t]{
		values: slices.Clone(values),
	}
}

func (q *Queue[t]) Len() int {
	return len(q.values)
}

func (q *Queue[t]) Reset() {
	q.index = 0
	q.values = q.values[:0]
}

func (q *Queue[t]) Enqueue(v t) {
	q.values = append(q.values, v)
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
