package blocks

import "errors"

var (
	ErrorDataNotFound   error = errors.New("data not found")
	ErrorNotImplemented error = errors.New("not implemented")
)
