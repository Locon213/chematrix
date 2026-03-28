// Package chemistry provides core types for representing chemical elements and compounds.
package chemistry

// Element represents a chemical element from the periodic table.
type Element struct {
	Symbol          string // e.g., "Fe", "O"
	Name            string // e.g., "Iron", "Oxygen"
	AtomicNumber    int    // e.g., 26
	AtomicMass      float64
	OxidationStates []int // common oxidation states, e.g., [-2, +2, +3] for Fe
	Group           int   // periodic table group (for heuristics)
}

// Composition is a map of element symbols to atom counts.
// Example: {"H": 2, "O": 1} for H2O
type Composition map[string]int

// Copy returns a deep copy of the composition.
func (c Composition) Copy() Composition {
	result := make(Composition, len(c))
	for k, v := range c {
		result[k] = v
	}
	return result
}

// Add merges another composition into this one, summing counts for common elements.
func (c Composition) Add(other Composition) {
	for elem, count := range other {
		c[elem] += count
	}
}

// Multiply scales all atom counts by a factor.
func (c Composition) Multiply(factor int) {
	for elem := range c {
		c[elem] *= factor
	}
}

// TotalAtoms returns the total number of atoms in the composition.
func (c Composition) TotalAtoms() int {
	total := 0
	for _, count := range c {
		total += count
	}
	return total
}

// Equals compares two compositions for equality.
func (c Composition) Equals(other Composition) bool {
	if len(c) != len(other) {
		return false
	}
	for elem, count := range c {
		if other[elem] != count {
			return false
		}
	}
	return true
}
