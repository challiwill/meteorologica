package errare

import "fmt"

type CreationError struct {
	thing  string
	reason string
}

func NewCreationError(thing, reason string) CreationError {
	return CreationError{
		thing:  thing,
		reason: reason,
	}
}

func (e CreationError) Error() string {
	return fmt.Sprintf("Failed to create %s: %s", e.thing, e.reason)
}
