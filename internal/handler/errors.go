package handler

import "errors"

var errUnauthorized = errors.New("unauthorized")
var errMissingFieldParam = errors.New("missing required query param: field")
var errNotActiveOrNotYours = errors.New("not found, not yours, or already left/disqualified")
