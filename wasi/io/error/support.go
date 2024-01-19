package error

// Error implements the error interface.
func (self Error) Error() string {
	return self.ToDebugString()
}
