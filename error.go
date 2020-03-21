package slice

import (
	"fmt"
	"strings"
)

type errBootFailed []error

// Append appends boot error if it not nil.
func (e errBootFailed) Append(err error) errBootFailed {
	if err != nil {
		return append(e, err)
	}
	return e
}

// Error implements error interface.
func (e errBootFailed) Error() string {
	var str []string
	for _, err := range e {
		str = append(str, err.Error())
	}
	return fmt.Sprintf("boot failed: %s", strings.Join(str, ", "))
}

type errShutdownFailed []error

// Append appends shutdown error if it not nil.
func (e errShutdownFailed) Append(err error) errShutdownFailed {
	if err != nil {
		return append(e, err)
	}
	return e
}

// Error implements error interface.
func (e errShutdownFailed) Error() string {
	var str []string
	for _, err := range e {
		str = append(str, err.Error())
	}
	return fmt.Sprintf("shutdown failed: %s", strings.Join(str, ", "))
}

type errKernelResolveFailed struct {
	err error
}

// Error implements error interface.
func (e errKernelResolveFailed) Error() string {
	return fmt.Sprintf("kernel resolve failed: %s", e.err)
}

type provideErrors []provideError

// Error implements error interface.
func (p provideErrors) Error() string {
	var str []string
	for i := 0; len(p) > 0; i++ {
		str = append(str, p[i].Error())
	}
	return strings.Join(str, ", ")
}

type provideError struct {
	err error
}

// Error implements error interface.
func (p provideError) Error() string {
	return fmt.Sprintf("provide failed: %s", p.err)
}
