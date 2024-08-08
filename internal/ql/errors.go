package ql

import "fmt"

// ErrNoField is returned when a query specifies a field that does not exist.
// For example, "asdf > 2" would return this error (unless you have registered
// a custom field called "asdf" with the *Parser).
type ErrNoField struct {
	Field string
}

func (e *ErrNoField) Error() string {
	return fmt.Sprintf("no such field: %s", e.Field)
}

// ErrNoOperationForField is returned when a field does not support a given operator.
// This occurs when parsing a query like "name < Drannith", because there is no
// less than operator for the name field.
type ErrNoOperationForField struct {
	Field    string
	Operator operator
}

func (e *ErrNoOperationForField) Error() string {
	return fmt.Sprintf("field %q has no such operation: %s", e.Field, e.Operator)
}
