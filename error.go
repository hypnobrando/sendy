package sendy

import "fmt"

type (
	// Error is a struct that implements the native
	// error interface and is instantiated after making
	// and handling an HTTP request.
	Error struct {
		err        error
		statusCode int
	}
)

var _ error = (*Error)(nil)

// Error implements the native error interface.
func (err *Error) Error() string {
	if err.err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Status code: %d", err.statusCode)
}
