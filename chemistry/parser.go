// Package chemistry provides core types for representing chemical elements and compounds.
package chemistry

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// ParseFormula converts a chemical formula string into a Substance.
//
// Supports:
//   - Elements: "H", "Fe", "Og"
//   - Numbers: "H2O", "Fe2(SO4)3"
//   - Parentheses: "Ca(OH)2", "Al2(SO4)3"
//   - Ionic charges: "SO4^2-", "NH4+", "Fe^3+"
//   - Hydrates: "CuSO4·5H2O" (dot notation)
//
// Returns error for invalid syntax or unknown elements.
//
// Examples:
//
//	ParseFormula("H2SO4") → Substance{Composition: {"H":2,"S":1,"O":4}, Charge:0}
//	ParseFormula("OH^-")  → Substance{Composition: {"O":1,"H":1}, Charge:-1}
func ParseFormula(formula string) (*Substance, error) {
	if formula == "" {
		return nil, fmt.Errorf("parse error: empty formula")
	}

	// Extract state symbol if present (e.g., "(s)", "(aq)")
	state := ""
	formula, state = extractState(formula)

	// Extract charge if present (e.g., "^2-", "^+")
	charge := 0
	formula, charge = extractCharge(formula)

	// Parse the main formula
	composition, err := parseFormulaBody(formula)
	if err != nil {
		return nil, err
	}

	return &Substance{
		Formula:     normalizeFormula(formula, charge, state),
		Composition: composition,
		Charge:      charge,
		State:       state,
	}, nil
}

// extractState removes and returns the state symbol from the end of the formula.
func extractState(formula string) (string, string) {
	if len(formula) < 2 || formula[len(formula)-1] != ')' {
		return formula, ""
	}

	// Find the matching opening parenthesis
	for i := len(formula) - 2; i >= 0; i-- {
		if formula[i] == '(' {
			state := formula[i:len(formula)]
			// Validate common state symbols
			validStates := map[string]bool{"(s)": true, "(l)": true, "(g)": true, "(aq)": true}
			if validStates[state] {
				return formula[:i], state
			}
			return formula, ""
		}
	}
	return formula, ""
}

// extractCharge removes and returns the charge from the end of the formula.
// Supports formats: "^2-", "^+", "^3+", "^-", "Na+", "Cl-", etc.
func extractCharge(formula string) (string, int) {
	if formula == "" {
		return formula, 0
	}

	// Look for '^' which indicates explicit charge notation
	caretIndex := -1
	for i := len(formula) - 1; i >= 0; i-- {
		if formula[i] == '^' {
			caretIndex = i
			break
		}
	}

	if caretIndex >= 0 {
		// Found ^, parse charge from caretIndex+1 to end
		chargePart := formula[caretIndex+1:]
		charge, err := parseChargeValue(chargePart)
		if err != nil {
			return formula, 0
		}
		return formula[:caretIndex], charge
	}

	// No ^ found, check for simple ionic format like "Na+" or "Cl-" at the very end
	lastChar := formula[len(formula)-1]
	if lastChar == '+' || lastChar == '-' {
		// Look for the last non-digit position
		i := len(formula) - 2
		for i >= 0 && unicode.IsDigit(rune(formula[i])) {
			i--
		}

		if i >= 0 {
			// Check if the character before the sign is a letter (element symbol end)
			prevRune := rune(formula[i])
			if unicode.IsLetter(prevRune) {
				// Only treat as charge if it's directly after a letter with no number
				if i == len(formula)-2 || !unicode.IsDigit(rune(formula[i+1])) {
					// Check if adding this charge makes sense chemically
					baseFormula := formula[:i+1]
					// For simple ions like Na+, Cl-, OH-, etc.
					if isLikelySimpleIon(baseFormula) {
						charge := 1
						if lastChar == '-' {
							charge = -1
						}
						return baseFormula, charge
					}
				}
			}
		}
	}

	return formula, 0
}

// isLikelySimpleIon checks if a formula is likely to be a simple ion.
func isLikelySimpleIon(formula string) bool {
	// Common simple cations
	cations := map[string]bool{
		"Na": true, "K": true, "Li": true, "Ca": true, "Mg": true,
		"Ba": true, "Sr": true, "Al": true, "Zn": true, "Fe": true,
		"Cu": true, "Ag": true, "Pb": true, "Sn": true, "Ni": true,
		"Co": true, "Mn": true, "Cr": true, "Hg": true, "NH4": true,
		"H": true,
	}
	// Common simple anions
	anions := map[string]bool{
		"Cl": true, "Br": true, "I": true, "F": true, "O": true,
		"S": true, "OH": true, "CN": true, "NO3": true, "NO2": true,
		"SO4": true, "SO3": true, "CO3": true, "PO4": true, "ClO3": true,
		"ClO4": true, "MnO4": true, "CrO4": true, "Cr2O7": true,
		"HCO3": true, "HSO4": true, "H2PO4": true, "HPO4": true,
	}

	_, isCation := cations[formula]
	_, isAnion := anions[formula]
	return isCation || isAnion
}

// parseChargeValue parses the numeric part of a charge (e.g., "2-" → -2, "+" → +1, "3+" → +3).
func parseChargeValue(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty charge value")
	}

	// Handle formats: "2-", "3+", "+", "-", "2", etc.
	// The sign can be at the beginning or end
	sign := 1
	numStr := s

	// Check for sign at the beginning
	if s[0] == '-' {
		sign = -1
		numStr = s[1:]
	} else if s[0] == '+' {
		sign = 1
		numStr = s[1:]
	}

	// Check for sign at the end
	if len(numStr) > 0 {
		lastChar := numStr[len(numStr)-1]
		if lastChar == '-' {
			sign = -1
			numStr = numStr[:len(numStr)-1]
		} else if lastChar == '+' {
			sign = 1
			numStr = numStr[:len(numStr)-1]
		}
	}

	if numStr == "" {
		return sign, nil // Just sign without number means ±1
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return sign * num, nil
}

// parseFormulaBody parses the main formula (without charge and state).
func parseFormulaBody(formula string) (Composition, error) {
	composition := make(Composition)

	// Handle hydrates (e.g., "CuSO4·5H2O")
	parts := splitHydrate(formula)
	for _, part := range parts {
		// Check if part starts with a number (hydrate coefficient)
		multiplier := 1
		numEnd := 0
		for numEnd < len(part) && unicode.IsDigit(rune(part[numEnd])) {
			numEnd++
		}
		if numEnd > 0 {
			var err error
			multiplier, err = strconv.Atoi(part[:numEnd])
			if err != nil {
				return nil, fmt.Errorf("parse error: invalid hydrate coefficient '%s'", part[:numEnd])
			}
			part = part[numEnd:]
		}

		partComp, err := parseSimpleFormula(part)
		if err != nil {
			return nil, err
		}

		partComp.Multiply(multiplier)
		composition.Add(partComp)
	}

	return composition, nil
}

// splitHydrate splits a hydrate formula by the middle dot (·).
func splitHydrate(formula string) []string {
	// Check for middle dot (·) or regular dot (.)
	for i := 0; i < len(formula); {
		r, size := utf8.DecodeRuneInString(formula[i:])
		if r == '·' || r == '.' {
			// Check if there's a number after the dot
			if i+size < len(formula) {
				nextRune, _ := utf8.DecodeRuneInString(formula[i+size:])
				if unicode.IsDigit(nextRune) {
					return []string{formula[:i], formula[i+size:]}
				}
			}
		}
		i += size
	}
	return []string{formula}
}

// parseSimpleFormula parses a formula without hydrates.
func parseSimpleFormula(formula string) (Composition, error) {
	composition := make(Composition)
	i := 0

	for i < len(formula) {
		r, _ := utf8.DecodeRuneInString(formula[i:])
		if r == utf8.RuneError {
			return nil, fmt.Errorf("parse error: invalid UTF-8 at position %d", i)
		}

		if r == '(' {
			// Handle parentheses
			closeIndex := findMatchingParen(formula, i)
			if closeIndex == -1 {
				return nil, fmt.Errorf("parse error: unmatched parenthesis at position %d", i)
			}

			// Parse content inside parentheses
			innerFormula := formula[i+1 : closeIndex]
			innerComp, err := parseSimpleFormula(innerFormula)
			if err != nil {
				return nil, err
			}

			// Check for multiplier after closing parenthesis
			multiplier := 1
			j := closeIndex + 1
			numStart := j
			for j < len(formula) && unicode.IsDigit(rune(formula[j])) {
				j++
			}
			if j > numStart {
				multiplier, err = strconv.Atoi(formula[numStart:j])
				if err != nil {
					return nil, fmt.Errorf("parse error: invalid number at position %d", numStart)
				}
			}

			innerComp.Multiply(multiplier)
			composition.Add(innerComp)
			i = j
		} else if r == ')' {
			// Unmatched closing parenthesis
			return nil, fmt.Errorf("parse error: unmatched closing parenthesis at position %d", i)
		} else if unicode.IsUpper(r) {
			// Parse element symbol
			elemStart := i
			i++ // Move past the uppercase letter

			// Check for lowercase continuation
			for i < len(formula) && unicode.IsLower(rune(formula[i])) {
				i++
			}

			elemSymbol := formula[elemStart:i]

			// Validate element
			if !IsValidElement(elemSymbol) {
				return nil, fmt.Errorf("parse error: unknown element '%s' at position %d", elemSymbol, elemStart)
			}

			// Parse count
			count := 1
			numStart := i
			for i < len(formula) && unicode.IsDigit(rune(formula[i])) {
				i++
			}
			if i > numStart {
				var err error
				count, err = strconv.Atoi(formula[numStart:i])
				if err != nil {
					return nil, fmt.Errorf("parse error: invalid number at position %d", numStart)
				}
			}

			composition[elemSymbol] += count
		} else {
			return nil, fmt.Errorf("parse error: unexpected character '%c' at position %d", r, i)
		}
	}

	return composition, nil
}

// findMatchingParen finds the index of the closing parenthesis matching the one at openIndex.
func findMatchingParen(s string, openIndex int) int {
	if openIndex >= len(s) || s[openIndex] != '(' {
		return -1
	}

	depth := 0
	for i := openIndex; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1 // No matching parenthesis found
}

// normalizeFormula creates a normalized formula string.
func normalizeFormula(formula string, charge int, state string) string {
	result := formula
	if state != "" {
		result += state
	}
	if charge != 0 {
		if charge == 1 {
			result += "^+"
		} else if charge == -1 {
			result += "^-"
		} else if charge > 0 {
			result += fmt.Sprintf("^%d+", charge)
		} else {
			result += fmt.Sprintf("^%d-", -charge)
		}
	}
	return result
}
