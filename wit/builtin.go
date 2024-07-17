package wit

type cmpOrdered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func min[T cmpOrdered](x T, y ...T) T {
	for _, v := range y {
		if v < x {
			x = v
		}
	}
	return x
}

func max[T cmpOrdered](x T, y ...T) T {
	for _, v := range y {
		if x > v {
			x = v
		}
	}
	return x
}
