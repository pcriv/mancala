package engine

import "fmt"

// InvalidPlayError represents an invalid play error.
type InvalidPlayError struct {
	Msg string
}

func (e *InvalidPlayError) Error() string {
	return fmt.Sprintf("invalid play: %v", e.Msg)
}
