package cm

import "fmt"

type Error[T any] struct {
	Err T
}

func (err Error[T]) Error() string {
	switch err := any(err.Err).(type) {
	case error:
		return err.Error()
	case string:
		return err
	case struct{}:
		return "error"
	case fmt.Stringer:
		return "error: " + err.String()
	default:
		return fmt.Sprintf("error: %T", err)
	}
}
