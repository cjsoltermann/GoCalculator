package stack

type Stack[T any] struct {
	stack []T
}

type Queue[T any] struct {
	queue []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{make([]T, 0, 64)}
}

func QueueFromSlice[T any](slice []T) *Queue[T] {
	return &Queue[T]{slice}
}

func (q *Queue[T]) Push(val T) {
	q.queue = append(q.queue, val)
}

func (q *Queue[T]) Pop() T {
	val := q.queue[0]
	q.queue = q.queue[1:]
	return val
}

func (q *Queue[T]) Peek() any {
	return q.queue[0]
}

func (q Queue[T]) IsEmpty() bool {
	return len(q.queue) == 0
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{make([]T, 0, 64)}
}

func StackFromSlice[T any](slice []T) *Stack[T] {
	revSlice := make([]T, len(slice), cap(slice))
	for i, s := range slice {
		revSlice[len(revSlice)-i-1] = s
	}
	return &Stack[T]{revSlice}
}

func (s *Stack[T]) Push(val T) {
	s.stack = append(s.stack, val)
}

func (s *Stack[T]) Pop() T {
	val := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return val
}

func (s Stack[T]) Peek() T {
	return s.stack[len(s.stack)-1]
}

func (s Stack[T]) IsEmpty() bool {
	return len(s.stack) == 0
}
