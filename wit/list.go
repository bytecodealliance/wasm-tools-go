package wit

type List[E, C any] []*E

func (list *List[E, C]) Element(i int) *E {
	return remake(list, i)
}

func (list *List[E, C]) Visitor(ctx C) (any, error) {
	return &ListVisitor[E, C]{list, ctx}, nil
}

type ListVisitor[E, C any] struct {
	*List[E, C]
	ctx C
}

func (d *ListVisitor[E, C]) DecodeElement(i int) (*E, error) {
	return d.Element(i), nil
}

// remake returns the value of slice s at index i,
// reallocating the slice if necessary. s must be a slice
// of pointers, because the underlying backing to s might
// change when reallocated.
// If the value at s[i] is nil, a new *E will be allocated.
func remake[S ~[]*E, E any](s *S, i int) *E {
	if i < 0 {
		return nil
	}
	if i >= len(*s) {
		*s = append(*s, make([]*E, i-len(*s))...)
	}
	if (*s)[i] == nil {
		(*s)[i] = new(E)
	}
	return (*s)[i]
}
