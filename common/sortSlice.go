package common

import (
	"errors"
	"fmt"
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

func (s *SortSlice[t]) Remove(v t) (removed bool) {
	fmt.Printf("before %v\n", len(s.values))
	s.values = slices.DeleteFunc(s.values, func(a t) bool { return s.cmp(v, a) == 0 })
	fmt.Printf("after %v\n", len(s.values))
	return true
}

func (s *SortSlice[t]) GetByIndex(index int) (*t, error) {
	if index >= len(s.values) {
		return nil, errors.New("Out of bound")
	}
	return new(s.values[index]), nil
}

func (s *SortSlice[t]) ForEach(callback func(index int, value t) (stop bool, err error)) error {
	var stop bool
	var err error
	for i := range s.values {
		stop, err = callback(i, s.values[i])
		if err != nil {
			return nil
		}
		if stop {
			return nil
		}
	}
	return nil
}

func (s *SortSlice[t]) GetAll() []t {
	return s.values
}
