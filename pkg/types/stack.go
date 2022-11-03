package types

// Stack struct which contains a list of params
type Stack[T any] struct {
	items []*T
}

// NewEmptyStack returns a new instance of Stack with zero elements
func NewEmptyStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: nil,
	}
}

// Push adds new item to top of existing/empty stack
func (s *Stack[T]) Push(item *T) {
	s.items = append(s.items, item)
}

// Pop removes most recent item(top) from stack
func (s *Stack[T]) Pop() *T {
	if len(s.items) == 0 {
		return nil
	}

	lastItem := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return lastItem
}

// Top read most recent item(top) from stack w/o removing
func (s *Stack[T]) Top() *T {
	if len(s.items) == 0 {
		return nil
	}

	return s.items[len(s.items)-1]
}
