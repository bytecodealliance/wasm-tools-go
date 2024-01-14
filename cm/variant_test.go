package cm

var (
	// TODO: nuke Variant2 interface?
	_ Variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ Variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ Variant2[string, bool]       = &SizedVariant2[Shape[string], string, bool]{}
	_ Variant2[bool, string]       = &SizedVariant2[Shape[string], bool, string]{}
)

func zeroPtr[T any]() *T {
	var zero T
	return &zero
}
