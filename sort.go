package slice

const (
	temporary = 1
	permanent = 2
)

// sortBundles is a step of application bootstrap.
func sortBundles(bundles []Bundle) ([]Bundle, bool) {
	var sorted []Bundle
	marks := map[string]int{}
	for _, b := range bundles {
		if !visit(b, marks, &sorted) {
			return sorted, false
		}
	}
	return sorted, true
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
		*sorted = append(*sorted, b)
		return true
	}
	marks[b.Name] = temporary
	for _, dep := range b.Bundles {
		if !visit(dep, marks, sorted) {
			return false
		}
	}
	marks[b.Name] = permanent
	*sorted = append(*sorted, b)
	return true
}
