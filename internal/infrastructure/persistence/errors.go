package persistence

import "fmt"

type (
	// ErrNotFound happens when the request record is not found on the storage.
	ErrNotFound struct {
		Msg string
	}
)

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("not found error: %v", e.Msg)
}
