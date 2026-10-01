package repository

import "errors"

var ErrConflict = errors.New("record already exists")

type InvalidDataError struct {
	Constraint string
}

func (e *InvalidDataError) Error() string {
	return "data violates constraint " + e.Constraint
}