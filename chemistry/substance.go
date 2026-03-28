package chemistry

import (
	"fmt"
	"sort"
	"strings"
)

// Substance represents a chemical compound or ion.
// Example: "Fe2(SO4)3" → Composition{"Fe":2, "S":3, "O":12}, Charge: 0
type Substance struct {
	Formula     string      // original input, normalized
	Composition Composition // map[element]count
	Charge      int         // ionic charge: -2, +1, 0, etc.
	State       string      // optional: "(s)", "(aq)", "(g)"
}

// String returns a human-readable representation of the substance.
func (s *Substance) String() string {
	var sb strings.Builder

	// Build formula from composition
	if s.Formula != "" {
		sb.WriteString(s.Formula)
	} else {
		sb.WriteString(compositionToFormula(s.Composition))
	}

	// Add charge if present
	if s.Charge != 0 {
		if s.Charge > 0 {
			sb.WriteString("^")
			sb.WriteString(intToRoman(s.Charge))
			sb.WriteString("+")
		} else {
			sb.WriteString("^")
			sb.WriteString(intToRoman(-s.Charge))
			sb.WriteString("-")
		}
	}

	// Add state if present
	if s.State != "" {
		sb.WriteString(s.State)
	}

	return sb.String()
}

// MolarMass calculates the molar mass of the substance in g/mol.
func (s *Substance) MolarMass() float64 {
	mass := 0.0
	for elem, count := range s.Composition {
		if el := GetElement(elem); el != nil {
			mass += el.AtomicMass * float64(count)
		}
	}
	return mass
}

// compositionToFormula converts a composition map back to a formula string.
// Note: This is a simplified implementation; order may not match IUPAC conventions.
func compositionToFormula(comp Composition) string {
	if len(comp) == 0 {
		return ""
	}

	// Sort elements for consistent output (alphabetically, but H and C first for organics)
	elements := make([]string, 0, len(comp))
	for elem := range comp {
		elements = append(elements, elem)
	}

	sort.Slice(elements, func(i, j int) bool {
		// Special ordering for common elements
		order := map[string]int{"H": 0, "C": 1, "N": 2, "O": 3, "S": 4, "P": 5}
		oi, oki := order[elements[i]]
		oj, okj := order[elements[j]]
		if oki && okj {
			return oi < oj
		}
		if oki {
			return true
		}
		if okj {
			return false
		}
		return elements[i] < elements[j]
	})

	var sb strings.Builder
	for _, elem := range elements {
		sb.WriteString(elem)
		count := comp[elem]
		if count > 1 {
			sb.WriteString(fmt.Sprintf("%d", count))
		}
	}

	return sb.String()
}

// intToRoman converts a small integer to Roman numerals (for charges).
func intToRoman(n int) string {
	if n <= 0 || n > 10 {
		return fmt.Sprintf("%d", n)
	}
	roman := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	return roman[n]
}
