package model

import "fmt"

var ErrDuplicateKey = fmt.Errorf("duplicate key")
var ErrNotFound = fmt.Errorf("not found")
var ErrForbidden error = fmt.Errorf("forbidden operation")
