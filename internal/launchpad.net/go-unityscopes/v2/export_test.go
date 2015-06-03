package scopes

// This file exports certain private functions for use by tests.

func NewTestingResult() *Result {
	return newTestingResult()
}
