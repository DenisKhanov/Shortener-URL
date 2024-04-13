// Package models defines common models and errors for the application.

package models

import "errors"

// ErrURLFound is an error indicating that a short URL is found in the database.
var ErrURLFound = errors.New("short URL found in database")

// ErrURLDeleted is an error indicating that a short URL is marked as deleted.
var ErrURLDeleted = errors.New("short URL marked as deleted")
