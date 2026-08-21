package common

import (
	"errors"
	"slices"
)

type SortSlice[t any] struct {
	values []t
	cmp    func(a, b t) int
}

func NewSortSlice[t any](cmpCustom func(a, b t) int) *SortSlice[t] {
	if cmpCustom == nil {
		return nil
	}
	return &SortSlice[t]{
		values: make([]t, 0),
		cmp:    cmpCustom,
	}
}

func (s *SortSlice[t]) Insert(v t) {
	index, _ := slices.BinarySearchFunc(s.values, v, s.cmp)
	s.values = append(s.values, *new(t))
	copy(s.values[index+1:], s.values[index:])
	s.values[index] = v
}

func (s *SortSlice[t]) Get(index int) (*t, error) {
	if index >= len(s.values) {
		return nil, errors.New("Out of bound")
	}
	return new(s.values[index]), nil
}
