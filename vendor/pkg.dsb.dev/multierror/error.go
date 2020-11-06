// Package multierror provides a mechanism for representing a list of error values as a single error.
package multierror

import (
	"errors"
	"strings"
)

type (
	multiError []error
)

// Append zero or more errors to a single error. Returns nil if the resulting
// slice would be empty.
func Append(err error, errs ...error) error {
	var e multiError
	if errors.As(err, &e) {
		return append(e, errs...)
	}

	sl := make([]error, 0)
	if err != nil {
		sl = append(sl, err)
	}

	for _, err = range errs {
		if err == nil {
			continue
		}

		sl = append(sl, err)
	}

	if len(sl) == 0 {
		return nil
	}

	return multiError(sl)
}

// Error returns a string containing all appended errors separated by semicolons.
func (e multiError) Error() string {
	arr := make([]string, len(e))
	for i, err := range e {
		arr[i] = err.Error()
	}

	return strings.Join(arr, "; ")
}

// Is reports whether any error in the collection whose chain matches target. Returns true on
// the first match.
func (e multiError) Is(target error) bool {
	for _, err := range e {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// As finds the first error in the collection whose chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false. Returns true on
// the first match.
func (e multiError) As(target interface{}) bool {
	for _, err := range e {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}
