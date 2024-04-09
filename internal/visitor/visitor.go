package visitor

type Visitor[T comparable] struct {
	yield   func(T) bool
	visited map[T]struct{}
	done    bool
}

func New[T comparable](yield func(T) bool) *Visitor[T] {
	return &Visitor[T]{
		yield:   yield,
		visited: make(map[T]struct{}),
	}
}

func (v *Visitor[T]) Yield(e T) bool {
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

func (v *Visitor[T]) Visited(e T) bool {
	_, visited := v.visited[e]
	return visited
}

func (v *Visitor[T]) Done() bool {
	return v.done
}
