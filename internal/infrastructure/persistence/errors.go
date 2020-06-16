package persistence

import "fmt"

type (
	// NotFoundError happens when the request record is not found on the storage.
	NotFoundError struct {
		Msg string
	}
)

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found error: %v", e.Msg)
}
