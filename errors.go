package main

import "errors"

var (
	errValidation = errors.New("validation")
	errNotFound   = errors.New("not found")
)
