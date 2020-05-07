package resources

type (
	// ValidationError represents a validation error
	ValidationError struct {
		Field string `json:"field"`
		Msg   string `json:"message"`
	}

	// ValidationErrors represents a collection of ValidationError
	ValidationErrors struct {
		Errors []ValidationError `json:"errors"`
	}
)

// Add adds an error with the given field and msg to the collection
func (e *ValidationErrors) Add(field string, msg string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Msg: msg})
}

// Any checks if the collection contains any errors
func (e *ValidationErrors) Any() bool {
	return len(e.Errors) > 0
}
