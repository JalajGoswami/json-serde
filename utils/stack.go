package utils

import "errors"

type Stack[T any] struct {
	list []T
}

var ErrStackEmpty = errors.New("stack is empty")

func (s *Stack[T]) Push(item T) {
	s.list = append(s.list, item)
}

func (s *Stack[T]) Pop() error {
	if len(s.list) == 0 {
		return ErrStackEmpty
	}
	s.list = s.list[:len(s.list)-1]
	return nil
}

func (s *Stack[T]) Top() (item T, isEmpty bool) {
	if len(s.list) == 0 {
		return item, true
	}
	return s.list[len(s.list)-1], false
}
