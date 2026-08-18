package handler

import "errors"

var errUnauthorized = errors.New("unauthorized")
var errMissingFieldParam = errors.New("missing required query param: field")
