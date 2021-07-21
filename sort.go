package slice

import (
	"fmt"
)

const (
	temporary = 1
	permanent = 2
)

// prepareBundles is a step of application bootstrap.
func prepareBundles(bundles []Bundle) ([]Bundle, error) {
	var sorted []Bundle
	marks := map[string]int{}
	for i, b := range bundles {
		if b.Name == "" {
			return nil, fmt.Errorf("bundle with index %d: empty name", i)
		}
		if !visit(b, marks, &sorted) {
			return sorted, fmt.Errorf("bundle cyclic detected") // todo: improve error message
		}
	}
	return sorted, nil
}

// visit
func visit(b Bundle, marks map[string]int, sorted *[]Bundle) bool {
	if marks[b.Name] == permanent {
		return true
	}
	if marks[b.Name] == temporary {
		// acyclic
		return false
	}
	if len(b.Bundles) == 0 {
		marks[b.Name] = permanent
		*sorted = append([]Bundle{b}, *sorted...)
		return true
	}
	marks[b.Name] = temporary
	for _, dep := range b.Bundles {
		if !visit(dep, marks, sorted) {
			return false
		}
	}
	marks[b.Name] = permanent
	*sorted = append([]Bundle{b}, *sorted...)
	return true
}
