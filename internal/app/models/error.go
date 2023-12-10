package models

import "errors"

var ErrURLFound = errors.New("short URL found in database")
var ErrURLDeleted = errors.New("short URL mark as deleted")
