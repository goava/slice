package slice

import (
	"fmt"
	"os"
	"reflect"
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

func printStartError(err error) {
	fmt.Printf("%s", err)
	os.Exit(1)
}
