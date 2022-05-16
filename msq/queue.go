package msq

type Queue[T any] struct {
	ch chan T
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		ch: make(chan T, cap),
	}
}

func (m *Queue[T]) Push(data T) bool {
	select {
	case m.ch <- data:
		return true
	default:
		return false
	}
}

func (m *Queue[T]) Pop() (T, bool) {
	var empty T
	select {
	case data := <-m.ch:
		return data, true
	default:
		return empty, false
	}
}

func (m *Queue[T]) PushChan() chan<- T {
	return m.ch
}

func (m *Queue[T]) PopChan() <-chan T {
	return m.ch
}
