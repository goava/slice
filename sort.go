package slice

import "reflect"

const (
	temporary = 1
	permanent = 2
)

// sortBundles is a step of application bootstrap.
func sortBundles(bundles []Bundle) ([]Bundle, bool) {
	var sorted []Bundle
	marks := map[reflect.Type]int{}
	for _, b := range bundles {
		if !visit(b, marks, &sorted) {
			return sorted, false
		}
	}
	return sorted, true
}

// visit
func visit(b Bundle, marks map[reflect.Type]int, sorted *[]Bundle) bool {
	typ := reflect.TypeOf(b)
	if marks[typ] == permanent {
		return true
	}
	if marks[typ] == temporary {
		// acyclic
		return false
	}
	dependOn, ok := b.(ComposedBundle)
	if !ok {
		marks[typ] = permanent
		*sorted = append(*sorted, b)
		return true
	}
	marks[typ] = temporary
	deps := dependOn.Bundles()
	for _, dep := range deps {
		if !visit(dep, marks, sorted) {
			return false
		}
	}
	marks[typ] = permanent
	*sorted = append(*sorted, b)
	return true
}
