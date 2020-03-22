package slice

import (
	"fmt"
	"reflect"
	"strings"
)

type bundleDIErrors struct {
	bundle interface{}
	list   []bundleDIError
}

// Error implements error interface.
func (p bundleDIErrors) Error() string {
	hash := map[string]bool{}
	for i := 0; i < len(p.list); i++ {
		hash[p.list[i].Error()] = true
	}
	var strs []string
	for k := range hash {
		strs = append(strs, k)
	}
	return fmt.Sprintf("%s: Provide bundle components failed", reflect.TypeOf(p.bundle))
}

type bundleDIError struct {
	err error
}

// Error implements error interface.
func (p bundleDIError) Error() string {
	return fmt.Sprintf("%s", p.err)
}

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
	return fmt.Sprintf("boot failed:\n %s", strings.Join(str, "\n"))
}

type errKernelResolveFailed struct {
	err error
}

// Error implements error interface.
func (e errKernelResolveFailed) Error() string {
	return fmt.Sprintf("kernel resolve failed: %s", e.err)
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
	return fmt.Sprintf("shutdown failed:\n %s", strings.Join(str, "\n"))
}
