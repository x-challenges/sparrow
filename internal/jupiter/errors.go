package jupiter

import "errors"

var (
	// Not found status code
	ErrNotFoundStatusCode = errors.New("not found status code response")

	// Unexpected status code response
	ErrUnexpectedStatusCode = errors.New("unexpected status code response")

	// Route not found
	ErrRouteNotFound = errors.New("routes not found")
)
