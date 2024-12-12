package jupiter

import "errors"

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code response")
	ErrRouteNotFound        = errors.New("routes not found")
)
