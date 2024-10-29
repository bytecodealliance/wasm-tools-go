package cm

import "fmt"

// Error is an interface that adapts a Component Model error to the Go error interface.
// The Err method returns the underlying error.
type Error[Err any] interface {
	error

	// Err returns the underlying error value.
	Err() Err
}

var (
	_ Error[string] = &cmError[string]{}
	_ error         = &cmError[string]{}
	_ error         = Error[string](nil)
)

// cmError wraps T to act as an error.
type cmError[Err any] struct {
	err Err
}

func (err cmError[Err]) Err() Err {
	return err.err
}

func (err cmError[Err]) Error() string {
	switch v := any(err.err).(type) {
	case string:
		return v
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	}
	return fmt.Sprintf("error: %v", err.err)
}
