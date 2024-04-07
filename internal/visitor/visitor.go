package visitor

type Visitor[T comparable] interface {
	Yield(T) bool
	Visited(T) bool
	Done() bool
}

func New[T comparable](yield func(T) bool) Visitor[T] {
	return &visitor[T]{
		yield:   yield,
		visited: make(map[T]struct{}),
	}
}

type visitor[T comparable] struct {
	yield   func(T) bool
	visited map[T]struct{}
	done    bool
}

func (v *visitor[T]) Yield(e T) bool {
	if v.done {
		return false
	}
	if v.Visited(e) {
		return true
	}
	v.visited[e] = struct{}{}
	v.done = !v.yield(e)
	return !v.done
}

func (v *visitor[T]) Visited(e T) bool {
	_, visited := v.visited[e]
	return visited
}

func (v *visitor[T]) Done() bool {
	return v.done
}
