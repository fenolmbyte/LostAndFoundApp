package app_errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNoChanges = errors.New("no changes detected")
var ErrPermissionDenied = errors.New("permission denied")
var ErrFileNotFound = errors.New("file not found")
var ErrUnauthorized = errors.New("you are not authorized")
