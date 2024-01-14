package cm

var (
	_ Result[struct{}, struct{}] = &UntypedResult{}
	_ Result[struct{}, struct{}] = &UnsizedResult[struct{}, struct{}]{}
	_ Result[string, bool]       = &SizedResult[Shape[string], string, bool]{}
	_ Result[string, bool]       = &OKSizedResult[string, bool]{}
	_ Result[bool, string]       = &ErrSizedResult[bool, string]{}
)
