// Package chemistry provides core types for representing chemical elements and compounds.
package chemistry

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

// FormulaSuggestion represents a suggested formula with similarity score.
type FormulaSuggestion struct {
	Formula    string  // suggested formula
	Similarity float64 // similarity score (0-1)
	Reason     string  // why this suggestion was made
	Substance  *Substance
}

// EquationSuggestion represents a suggested equation correction.
type EquationSuggestion struct {
	Original     string // original equation
	Suggested    string // suggested correction
	Reason       string // what was fixed
	IsBalanced   bool   // whether the suggestion is balanced
	Coefficients []int  // balancing coefficients if applicable
}

// TypoPattern represents a common typo pattern with its fix.
type TypoPattern struct {
	Pattern     string
	Fix         string
	Description string
}

// SuggestFormulas returns a list of suggested formulas similar to the input.
func SuggestFormulas(input string) []FormulaSuggestion {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	suggestions := make([]FormulaSuggestion, 0)

	// Try to parse as-is first
	if substance, err := ParseFormula(input); err == nil {
		suggestions = append(suggestions, FormulaSuggestion{
			Formula:    input,
			Similarity: 1.0,
			Reason:     "Valid formula",
			Substance:  substance,
		})
	}

	// Generate suggestions based on chemical rules
	ruleSuggestions := generateRuleBasedSuggestions(input)
	suggestions = append(suggestions, ruleSuggestions...)

	// Check for common typos and fixes using pattern matching
	typoSuggestions := detectAndFixTypos(input)
	suggestions = append(suggestions, typoSuggestions...)

	// Generate formulas with similar elements
	similarElementSuggestions := generateSimilarElementFormulas(input)
	suggestions = append(suggestions, similarElementSuggestions...)

	// Remove duplicates and sort by similarity
	suggestions = deduplicateSuggestions(suggestions)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Similarity > suggestions[j].Similarity
	})

	// Return top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions
}

// generateRuleBasedSuggestions generates suggestions based on chemical rules.
func generateRuleBasedSuggestions(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)
	normalized := normalizeInput(input)

	// Rule 1: Check for missing subscripts based on valence rules
	valenceSuggestions := suggestBasedOnValence(normalized)
	suggestions = append(suggestions, valenceSuggestions...)

	// Rule 2: Check for case errors (e.g., "cl" instead of "Cl")
	caseSuggestions := suggestBasedOnCaseErrors(normalized)
	suggestions = append(suggestions, caseSuggestions...)

	// Rule 3: Check for zero/O confusion
	zeroSuggestions := suggestBasedOnZeroConfusion(normalized)
	suggestions = append(suggestions, zeroSuggestions...)

	// Rule 4: Check for incomplete formulas (single element that should be diatomic)
	diatomicSuggestions := suggestBasedOnDiatomicRules(normalized)
	suggestions = append(suggestions, diatomicSuggestions...)

	// Rule 5: Check for common ion patterns
	ionSuggestions := suggestBasedOnIonPatterns(normalized)
	suggestions = append(suggestions, ionSuggestions...)

	return suggestions
}

// suggestBasedOnValence suggests formulas based on valence rules.
func suggestBasedOnValence(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)

	// Parse input to get elements
	elements := extractElementsWithCounts(input)
	if len(elements) == 0 {
		return suggestions
	}

	// Try to build valid compounds from the elements
	for _, elem1 := range elements {
		el1 := GetElement(elem1.Element)
		if el1 == nil || len(el1.OxidationStates) == 0 {
			continue
		}

		// Check if it's a single element that could form a compound with itself
		if len(elements) == 1 {
			// Single element - suggest common oxidation states
			for _, ox := range el1.OxidationStates {
				if ox != 0 {
					// Could form oxide
					oxygen := GetElement("O")
					if oxygen != nil {
						formula := buildCompoundFormula(el1, oxygen, ox, -2)
						if formula != "" && formula != input {
							if sub, err := ParseFormula(formula); err == nil {
								suggestions = append(suggestions, FormulaSuggestion{
									Formula:    formula,
									Similarity: 0.7,
									Reason:     fmt.Sprintf("Common oxide of %s (oxidation state %+d)", el1.Name, ox),
									Substance:  sub,
								})
							}
						}
					}
				}
			}
		}

		// For two-element combinations, check valence compatibility
		for _, elem2 := range elements {
			if elem1.Element == elem2.Element {
				continue
			}
			el2 := GetElement(elem2.Element)
			if el2 == nil || len(el2.OxidationStates) == 0 {
				continue
			}

			// Try common oxidation state combinations
			for _, ox1 := range el1.OxidationStates {
				for _, ox2 := range el2.OxidationStates {
					if ox1*ox2 < 0 { // Opposite signs can form compound
						formula := buildCompoundFormula(el1, el2, ox1, ox2)
						if formula != "" && formula != input {
							if sub, err := ParseFormula(formula); err == nil {
								similarity := calculateElementSimilarity(input, formula)
								if similarity > 0.5 {
									suggestions = append(suggestions, FormulaSuggestion{
										Formula:    formula,
										Similarity: similarity,
										Reason:     fmt.Sprintf("Valence-compatible compound: %s(%+d) + %s(%+d)", el1.Symbol, ox1, el2.Symbol, ox2),
										Substance:  sub,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return suggestions
}

// ElementCount represents an element with its count.
type ElementCount struct {
	Element string
	Count   int
}

// extractElementsWithCounts extracts elements and their counts from a formula.
func extractElementsWithCounts(formula string) []ElementCount {
	result := make([]ElementCount, 0)
	i := 0

	for i < len(formula) {
		if formula[i] >= 'A' && formula[i] <= 'Z' {
			elem := string(formula[i])
			i++
			if i < len(formula) && formula[i] >= 'a' && formula[i] <= 'z' {
				elem += string(formula[i])
				i++
			}

			count := 0
			numStart := i
			for i < len(formula) && formula[i] >= '0' && formula[i] <= '9' {
				count = count*10 + int(formula[i]-'0')
				i++
			}
			if count == 0 {
				count = 1
			}
			_ = numStart

			result = append(result, ElementCount{Element: elem, Count: count})
		} else {
			i++
		}
	}

	return result
}

// buildCompoundFormula builds a chemical formula from two elements and their oxidation states.
func buildCompoundFormula(el1, el2 *Element, ox1, ox2 int) string {
	if ox1 == 0 || ox2 == 0 {
		return ""
	}

	// Ensure ox1 is positive and ox2 is negative
	if ox1 > 0 && ox2 > 0 {
		return ""
	}
	if ox1 < 0 && ox2 < 0 {
		return ""
	}

	// Swap if needed
	if ox1 < 0 {
		el1, el2 = el2, el1
		ox1, ox2 = ox2, ox1
	}

	// Now ox1 > 0, ox2 < 0
	ox2 = -ox2

	// Find simplest ratio using GCD
	gcd := gcd(ox1, ox2)
	count1 := ox2 / gcd
	count2 := ox1 / gcd

	// Build formula
	formula := el1.Symbol
	if count1 > 1 {
		formula += fmt.Sprintf("%d", count1)
	}
	formula += el2.Symbol
	if count2 > 1 {
		formula += fmt.Sprintf("%d", count2)
	}

	return formula
}

// suggestBasedOnCaseErrors suggests fixes for case errors.
func suggestBasedOnCaseErrors(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)
	upper := strings.ToUpper(input)

	// Check if input is all uppercase (common typo)
	if upper != input && strings.ToUpper(upper) == upper {
		// Try to parse with proper casing
		fixed := fixElementCases(upper)
		if fixed != upper {
			if sub, err := ParseFormula(fixed); err == nil {
				suggestions = append(suggestions, FormulaSuggestion{
					Formula:    fixed,
					Similarity: 0.95,
					Reason:     "Fixed element symbol capitalization",
					Substance:  sub,
				})
			}
		}
	}

	return suggestions
}

// fixElementCases fixes capitalization of element symbols.
func fixElementCases(input string) string {
	result := make([]rune, len(input))
	i := 0

	for i < len(input) {
		r := rune(input[i])
		if r >= 'A' && r <= 'Z' {
			// Check if this could be a two-letter element
			if i+1 < len(input) && input[i+1] >= 'a' && input[i+1] <= 'z' {
				// Already proper case
				result[i] = r
				result[i+1] = rune(input[i+1])
				i += 2
				continue
			}

			// Check if uppercase followed by uppercase should be uppercase-lowercase
			if i+1 < len(input) && input[i+1] >= 'A' && input[i+1] <= 'Z' {
				// Check if uppercase-lowercase version is a valid element
				twoLetter := strings.ToUpper(string(input[i])) + strings.ToLower(string(input[i+1]))
				if IsValidElement(twoLetter) {
					result[i] = rune(input[i])
					result[i+1] = unicode.ToLower(rune(input[i+1]))
					i += 2
					continue
				}
			}

			result[i] = r
			i++
		} else if r >= 'a' && r <= 'z' {
			// Lowercase at start should be uppercase
			if i == 0 || (i > 0 && (input[i-1] < 'A' || input[i-1] > 'Z')) {
				result[i] = unicode.ToUpper(r)
			} else {
				result[i] = r
			}
			i++
		} else {
			result[i] = r
			i++
		}
	}

	return string(result)
}

// suggestBasedOnZeroConfusion suggests fixes for 0/O confusion.
func suggestBasedOnZeroConfusion(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)
	normalized := strings.ToUpper(input)

	// Check for zero instead of O
	if strings.Contains(normalized, "0") {
		fixed := strings.ReplaceAll(normalized, "0", "O")
		// Fix element cases
		fixed = fixElementCases(fixed)
		if sub, err := ParseFormula(fixed); err == nil {
			suggestions = append(suggestions, FormulaSuggestion{
				Formula:    fixed,
				Similarity: 0.95,
				Reason:     "Replaced zero (0) with oxygen (O)",
				Substance:  sub,
			})
		}
	}

	return suggestions
}

// suggestBasedOnDiatomicRules suggests diatomic forms for single elements.
func suggestBasedOnDiatomicRules(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)

	// Common diatomic elements
	diatomicElements := map[string]string{
		"H":  "H2",
		"N":  "N2",
		"O":  "O2",
		"F":  "F2",
		"Cl": "Cl2",
		"Br": "Br2",
		"I":  "I2",
	}

	elements := extractElementsWithCounts(input)
	if len(elements) == 1 {
		elem := elements[0].Element
		count := elements[0].Count

		// Single element with count 1 - suggest diatomic form
		if count == 1 {
			if diatomic, ok := diatomicElements[elem]; ok {
				if sub, err := ParseFormula(diatomic); err == nil {
					suggestions = append(suggestions, FormulaSuggestion{
						Formula:    diatomic,
						Similarity: 0.85,
						Reason:     fmt.Sprintf("%s is commonly found as diatomic %s", GetElement(elem).Name, diatomic),
						Substance:  sub,
					})
				}
			}
		}
	}

	return suggestions
}

// suggestBasedOnIonPatterns suggests common ion patterns.
func suggestBasedOnIonPatterns(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)

	// Common polyatomic ions with their charges
	polyatomicIons := map[string]struct {
		formula string
		charge  int
	}{
		"OH":   {"OH^-", -1},
		"SO4":  {"SO4^2-", -2},
		"NO3":  {"NO3^-", -1},
		"CO3":  {"CO3^2-", -2},
		"PO4":  {"PO4^3-", -3},
		"NH4":  {"NH4^+", 1},
		"ClO3": {"ClO3^-", -1},
		"MnO4": {"MnO4^-", -1},
	}

	elements := extractElementsWithCounts(input)
	if len(elements) > 1 {
		// Check if elements match a polyatomic ion pattern
		for ionName, ionData := range polyatomicIons {
			ionElements := extractElementsWithCounts(strings.TrimSuffix(ionData.formula, "^-"))
			if compositionsMatch(elements, ionElements) {
				if sub, err := ParseFormula(ionData.formula); err == nil {
					suggestions = append(suggestions, FormulaSuggestion{
						Formula:    ionData.formula,
						Similarity: 0.8,
						Reason:     fmt.Sprintf("Common polyatomic ion: %s", ionName),
						Substance:  sub,
					})
				}
			}
		}
	}

	return suggestions
}

// compositionsMatch checks if two element compositions match.
func compositionsMatch(elems1, elems2 []ElementCount) bool {
	if len(elems1) != len(elems2) {
		return false
	}

	map1 := make(map[string]int)
	map2 := make(map[string]int)

	for _, e := range elems1 {
		map1[e.Element] = e.Count
	}
	for _, e := range elems2 {
		map2[e.Element] = e.Count
	}

	for elem, count := range map1 {
		if map2[elem] != count {
			return false
		}
	}

	return true
}

// generateSimilarElementFormulas generates formulas with similar elements.
func generateSimilarElementFormulas(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)

	elements := extractElements(input)
	for _, elem := range elements {
		el := GetElement(elem)
		if el == nil {
			continue
		}

		// Find elements in the same group (similar chemical properties)
		similarElements := getElementsInSameGroup(el.Group)
		for _, similar := range similarElements {
			if similar.Symbol != elem {
				// Replace element in input with similar one
				replaced := replaceElementInFormula(input, elem, similar.Symbol)
				if sub, err := ParseFormula(replaced); err == nil {
					similarity := calculateElementSimilarity(input, replaced)
					suggestions = append(suggestions, FormulaSuggestion{
						Formula:    replaced,
						Similarity: similarity * 0.7,
						Reason:     fmt.Sprintf("Similar compound with %s (same group as %s)", similar.Symbol, elem),
						Substance:  sub,
					})
				}
			}
		}
	}

	return suggestions
}

// getElementsInSameGroup returns elements in the same periodic table group.
func getElementsInSameGroup(group int) []Element {
	if group == 0 {
		return nil
	}

	result := make([]Element, 0)
	for _, el := range GetAllElements() {
		if el.Group == group {
			result = append(result, el)
		}
	}
	return result
}

// replaceElementInFormula replaces an element symbol in a formula.
func replaceElementInFormula(formula, oldElem, newElem string) string {
	result := formula
	// Replace with word boundary consideration
	for i := 0; i < len(result); {
		if i < len(result) && result[i] >= 'A' && result[i] <= 'Z' {
			elem := string(result[i])
			j := i + 1
			if j < len(result) && result[j] >= 'a' && result[j] <= 'z' {
				elem += string(result[j])
				j++
			}
			if elem == oldElem {
				result = result[:i] + newElem + result[j:]
				i = i + len(newElem)
			} else {
				i = j
			}
		} else {
			i++
		}
	}
	return result
}

// calculateElementSimilarity calculates similarity based on shared elements.
func calculateElementSimilarity(formula1, formula2 string) float64 {
	elems1 := extractElements(formula1)
	elems2 := extractElements(formula2)

	if len(elems1) == 0 || len(elems2) == 0 {
		return 0
	}

	common := 0
	for _, e1 := range elems1 {
		for _, e2 := range elems2 {
			if e1 == e2 {
				common++
				break
			}
		}
	}

	// Jaccard similarity
	union := len(elems1) + len(elems2) - common
	return float64(common) / float64(union)
}

// detectAndFixTypos detects and fixes common typos using edit distance.
func detectAndFixTypos(input string) []FormulaSuggestion {
	suggestions := make([]FormulaSuggestion, 0)
	normalized := normalizeInput(input)

	// Generate typo fixes based on edit distance patterns
	edits := generateEditDistanceFixes(normalized)
	for _, edit := range edits {
		if sub, err := ParseFormula(edit); err == nil {
			distance := levenshteinDistance(normalized, strings.ToUpper(edit))
			maxLen := max(len(normalized), len(edit))
			similarity := 1.0 - float64(distance)/float64(maxLen)

			suggestions = append(suggestions, FormulaSuggestion{
				Formula:    edit,
				Similarity: similarity,
				Reason:     getTypoReason(normalized, edit),
				Substance:  sub,
			})
		}
	}

	return suggestions
}

// generateEditDistanceFixes generates possible fixes based on edit distance.
func generateEditDistanceFixes(input string) []string {
	fixes := make([]string, 0)
	seen := make(map[string]bool)

	// Single character edits
	for i := 0; i < len(input); i++ {
		// Deletion
		if len(input) > 1 {
			fix := input[:i] + input[i+1:]
			if !seen[fix] {
				fixes = append(fixes, fix)
				seen[fix] = true
			}
		}

		// Insertion (common characters)
		for _, ch := range []string{"2", "3", "4", "O", "H", "C", "N", "S"} {
			fix := input[:i] + ch + input[i:]
			if !seen[fix] {
				fixes = append(fixes, fix)
				seen[fix] = true
			}
		}

		// Substitution
		for _, ch := range []string{"O", "0", "2", "3"} {
			if string(input[i]) != ch {
				fix := input[:i] + ch + input[i+1:]
				if !seen[fix] {
					fixes = append(fixes, fix)
					seen[fix] = true
				}
			}
		}
	}

	// Transposition
	for i := 0; i < len(input)-1; i++ {
		if input[i] != input[i+1] {
			fix := input[:i] + string(input[i+1]) + string(input[i]) + input[i+2:]
			if !seen[fix] {
				fixes = append(fixes, fix)
				seen[fix] = true
			}
		}
	}

	return fixes
}

// getTypoReason returns a human-readable reason for a typo fix.
func getTypoReason(original, fixed string) string {
	distance := levenshteinDistance(strings.ToUpper(original), strings.ToUpper(fixed))

	if distance == 1 {
		if len(original) > len(fixed) {
			return "Possible typo: extra character"
		} else if len(original) < len(fixed) {
			return "Possible typo: missing character"
		}
		return "Possible typo: single character error"
	}

	return "Similar formula structure"
}

// normalizeInput normalizes input for comparison.
func normalizeInput(input string) string {
	return strings.TrimSpace(strings.ToUpper(input))
}

// deduplicateSuggestions removes duplicate suggestions.
func deduplicateSuggestions(suggestions []FormulaSuggestion) []FormulaSuggestion {
	seen := make(map[string]bool)
	result := make([]FormulaSuggestion, 0)

	for _, s := range suggestions {
		if !seen[s.Formula] {
			seen[s.Formula] = true
			result = append(result, s)
		}
	}

	return result
}

// SuggestEquations returns suggestions for equation corrections.
func SuggestEquations(input string) []EquationSuggestion {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	suggestions := make([]EquationSuggestion, 0)

	// Try to parse and balance as-is
	if result, err := BalanceEquation(input); err == nil && result.IsBalanced {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     input,
			Suggested:    formatBalancedEquation(result),
			Reason:       "Equation balanced successfully",
			IsBalanced:   true,
			Coefficients: result.Coefficients,
		})
	}

	// Generate suggestions based on reaction types
	typeSuggestions := generateTypeBasedEquationSuggestions(input)
	suggestions = append(suggestions, typeSuggestions...)

	// Fix common equation format issues
	formatSuggestions := fixEquationFormatIssues(input)
	suggestions = append(suggestions, formatSuggestions...)

	// Sort: balanced first, then by similarity
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].IsBalanced != suggestions[j].IsBalanced {
			return suggestions[i].IsBalanced
		}
		return suggestions[i].Reason < suggestions[j].Reason
	})

	// Return top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions
}

// generateTypeBasedEquationSuggestions generates suggestions based on reaction type analysis.
func generateTypeBasedEquationSuggestions(input string) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	reaction, err := ParseEquation(input)
	if err != nil {
		// Try to fix and parse
		fixedInput := fixArrowFormat(input)
		reaction, err = ParseEquation(fixedInput)
		if err != nil {
			return suggestions
		}
	}

	// Analyze reaction type and suggest based on chemical rules
	reactionType := analyzeReactionType(reaction)

	switch reactionType {
	case "combustion":
		// Combustion reactions: fuel + O2 -> CO2 + H2O
		suggestions = append(suggestions, suggestCombustionFixes(reaction)...)
	case "synthesis":
		// Synthesis: A + B -> AB
		suggestions = append(suggestions, suggestSynthesisFixes(reaction)...)
	case "decomposition":
		// Decomposition: AB -> A + B
		suggestions = append(suggestions, suggestDecompositionFixes(reaction)...)
	case "single_replacement":
		// Single replacement: A + BC -> AC + B
		suggestions = append(suggestions, suggestSingleReplacementFixes(reaction)...)
	case "double_replacement":
		// Double replacement: AB + CD -> AD + CB
		suggestions = append(suggestions, suggestDoubleReplacementFixes(reaction)...)
	case "acid_base":
		// Acid-base: acid + base -> salt + water
		suggestions = append(suggestions, suggestAcidBaseFixes(reaction)...)
	}

	return suggestions
}

// analyzeReactionType determines the type of reaction.
func analyzeReactionType(r *Reaction) string {
	numReactants := len(r.Reactants)
	numProducts := len(r.Products)

	// Check for combustion (O2 as reactant, CO2 and/or H2O as products)
	hasO2 := hasElement(r.Reactants, "O", 2)
	hasCO2 := hasCompound(r.Products, "CO2")
	hasH2O := hasCompound(r.Products, "H2O")

	if hasO2 && (hasCO2 || hasH2O) {
		return "combustion"
	}

	// Check for acid-base (H+ donor and OH- acceptor)
	hasAcid := hasAcidReactant(r)
	hasBase := hasBaseReactant(r)
	if hasAcid && hasBase {
		return "acid_base"
	}

	// Synthesis: multiple reactants, single product
	if numReactants > 1 && numProducts == 1 {
		return "synthesis"
	}

	// Decomposition: single reactant, multiple products
	if numReactants == 1 && numProducts > 1 {
		return "decomposition"
	}

	// Single replacement: element + compound -> element + compound
	if numReactants == 2 && numProducts == 2 {
		if isElement(r.Reactants[0]) && isCompound(r.Reactants[1]) &&
			isElement(r.Products[0]) && isCompound(r.Products[1]) {
			return "single_replacement"
		}
		// Double replacement: compound + compound -> compound + compound
		if isCompound(r.Reactants[0]) && isCompound(r.Reactants[1]) &&
			isCompound(r.Products[0]) && isCompound(r.Products[1]) {
			return "double_replacement"
		}
	}

	return "unknown"
}

// hasElement checks if a side contains a specific element with specific count.
func hasElement(components []ReactionComponent, elem string, count int) bool {
	for _, comp := range components {
		if comp.Substance.Composition[elem] == count {
			return true
		}
	}
	return false
}

// hasCompound checks if a side contains a specific compound.
func hasCompound(components []ReactionComponent, formula string) bool {
	for _, comp := range components {
		if comp.Substance.Formula == formula {
			return true
		}
	}
	return false
}

// hasAcidReactant checks if reactants contain an acid.
func hasAcidReactant(r *Reaction) bool {
	for _, comp := range r.Reactants {
		if strings.HasPrefix(comp.Substance.Formula, "H") {
			return true
		}
	}
	return false
}

// hasBaseReactant checks if reactants contain a base.
func hasBaseReactant(r *Reaction) bool {
	for _, comp := range r.Reactants {
		if strings.HasSuffix(comp.Substance.Formula, "OH") ||
			strings.Contains(comp.Substance.Formula, "OH^") {
			return true
		}
	}
	return false
}

// isElement checks if a component is a single element.
func isElement(comp ReactionComponent) bool {
	return len(comp.Substance.Composition) == 1
}

// isCompound checks if a component is a compound.
func isCompound(comp ReactionComponent) bool {
	return len(comp.Substance.Composition) > 1
}

// suggestCombustionFixes suggests fixes for combustion reactions.
func suggestCombustionFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	// Ensure O2 is present in reactants
	hasO2 := false
	for _, comp := range r.Reactants {
		if comp.Substance.Formula == "O2" {
			hasO2 = true
			break
		}
	}

	if !hasO2 {
		// Suggest adding O2
		suggestions = append(suggestions, EquationSuggestion{
			Original:   r.Raw,
			Suggested:  r.Raw + " + O2",
			Reason:     "Combustion reactions require O2",
			IsBalanced: false,
		})
	}

	// Try to balance
	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced combustion reaction",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// suggestSynthesisFixes suggests fixes for synthesis reactions.
func suggestSynthesisFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced synthesis reaction",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// suggestDecompositionFixes suggests fixes for decomposition reactions.
func suggestDecompositionFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced decomposition reaction",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// suggestSingleReplacementFixes suggests fixes for single replacement reactions.
func suggestSingleReplacementFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	// Check activity series rules
	activitySeries := []string{"Li", "K", "Ca", "Na", "Mg", "Al", "Zn", "Fe", "Ni", "Sn", "Pb", "H", "Cu", "Ag", "Au"}

	// Get the free element and the element in compound
	var freeElement, compoundElement string
	for _, comp := range r.Reactants {
		if len(comp.Substance.Composition) == 1 {
			for elem := range comp.Substance.Composition {
				freeElement = elem
				break
			}
		} else {
			for elem := range comp.Substance.Composition {
				compoundElement = elem
				break
			}
		}
	}

	if freeElement != "" && compoundElement != "" {
		// Check if reaction is favorable based on activity series
		freeIdx := getActivityIndex(freeElement, activitySeries)
		compoundIdx := getActivityIndex(compoundElement, activitySeries)

		if freeIdx > compoundIdx {
			suggestions = append(suggestions, EquationSuggestion{
				Original:  r.Raw,
				Suggested: r.Raw,
				Reason:    fmt.Sprintf("Reaction may not proceed: %s is less reactive than %s", freeElement, compoundElement),
			})
		}
	}

	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced single replacement reaction",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// suggestDoubleReplacementFixes suggests fixes for double replacement reactions.
func suggestDoubleReplacementFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	// Check for precipitate formation (simplified solubility rules)
	solubilityRules := map[string]bool{
		"NaCl": true, "KCl": true, "NaNO3": true, "KNO3": true,
		"AgCl": false, "PbCl2": false, "BaSO4": false, "CaCO3": false,
	}

	for _, comp := range r.Products {
		if !solubilityRules[comp.Substance.Formula] {
			suggestions = append(suggestions, EquationSuggestion{
				Original:  r.Raw,
				Suggested: r.Raw,
				Reason:    fmt.Sprintf("Possible precipitate: %s may form", comp.Substance.Formula),
			})
		}
	}

	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced double replacement reaction",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// suggestAcidBaseFixes suggests fixes for acid-base reactions.
func suggestAcidBaseFixes(r *Reaction) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	// Acid-base should produce salt + water
	hasWater := false
	for _, comp := range r.Products {
		if comp.Substance.Formula == "H2O" {
			hasWater = true
			break
		}
	}

	if !hasWater {
		suggestions = append(suggestions, EquationSuggestion{
			Original:  r.Raw,
			Suggested: r.Raw + " -> H2O + [salt]",
			Reason:    "Acid-base reactions typically produce water and a salt",
		})
	}

	if balanced, err := r.Balance(); err == nil {
		suggestions = append(suggestions, EquationSuggestion{
			Original:     r.Raw,
			Suggested:    formatBalancedEquation(balanced),
			Reason:       "Balanced acid-base neutralization",
			IsBalanced:   balanced.IsBalanced,
			Coefficients: balanced.Coefficients,
		})
	}

	return suggestions
}

// getActivityIndex returns the index of an element in the activity series.
func getActivityIndex(elem string, series []string) int {
	for i, e := range series {
		if e == elem {
			return i
		}
	}
	return len(series) // Not found, least reactive
}

// fixEquationFormatIssues fixes common format issues in equations.
func fixEquationFormatIssues(input string) []EquationSuggestion {
	suggestions := make([]EquationSuggestion, 0)

	// Fix arrow format
	fixed := fixArrowFormat(input)
	if fixed != input {
		if result, err := BalanceEquation(fixed); err == nil {
			suggestions = append(suggestions, EquationSuggestion{
				Original:     input,
				Suggested:    formatBalancedEquation(result),
				Reason:       "Fixed arrow format and balanced",
				IsBalanced:   result.IsBalanced,
				Coefficients: result.Coefficients,
			})
		}
	}

	// Fix spacing issues
	fixed = fixSpacing(input)
	if fixed != input {
		if result, err := BalanceEquation(fixed); err == nil {
			suggestions = append(suggestions, EquationSuggestion{
				Original:     input,
				Suggested:    formatBalancedEquation(result),
				Reason:       "Fixed spacing and balanced",
				IsBalanced:   result.IsBalanced,
				Coefficients: result.Coefficients,
			})
		}
	}

	return suggestions
}

// fixArrowFormat fixes arrow format issues.
func fixArrowFormat(input string) string {
	fixed := input
	// Handle unicode arrow first
	fixed = strings.ReplaceAll(fixed, "→", " -> ")
	// Handle various dash combinations
	fixed = strings.ReplaceAll(fixed, "==>", " -> ")
	fixed = strings.ReplaceAll(fixed, "=>", " -> ")
	fixed = strings.ReplaceAll(fixed, "-->", " -> ")
	fixed = strings.ReplaceAll(fixed, "= ", " -> ")
	fixed = strings.ReplaceAll(fixed, " = ", " -> ")
	// Normalize multiple spaces
	for strings.Contains(fixed, "  ") {
		fixed = strings.ReplaceAll(fixed, "  ", " ")
	}
	return strings.TrimSpace(fixed)
}

// fixSpacing fixes spacing issues.
func fixSpacing(input string) string {
	fixed := input
	// Fix spacing around +
	fixed = strings.ReplaceAll(fixed, " +", " + ")
	fixed = strings.ReplaceAll(fixed, "+ ", " + ")
	fixed = strings.ReplaceAll(fixed, "  ", " ")
	return strings.TrimSpace(fixed)
}

// formatBalancedEquation formats a balanced reaction as a string.
func formatBalancedEquation(result *BalancedReaction) string {
	if result == nil {
		return ""
	}

	var sb strings.Builder
	for i, comp := range result.Reaction.Reactants {
		if i > 0 {
			sb.WriteString(" + ")
		}
		coef := result.Coefficients[i]
		if coef > 1 {
			sb.WriteString(fmt.Sprintf("%d", coef))
		}
		sb.WriteString(comp.Substance.Formula)
	}

	sb.WriteString(" -> ")

	for i, comp := range result.Reaction.Products {
		if i > 0 {
			sb.WriteString(" + ")
		}
		coef := result.Coefficients[len(result.Reaction.Reactants)+i]
		if coef > 1 {
			sb.WriteString(fmt.Sprintf("%d", coef))
		}
		sb.WriteString(comp.Substance.Formula)
	}

	return sb.String()
}

// Helper functions

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// extractElements extracts element symbols from a formula string.
func extractElements(formula string) []string {
	elements := make([]string, 0)
	i := 0
	for i < len(formula) {
		if formula[i] >= 'A' && formula[i] <= 'Z' {
			elem := string(formula[i])
			i++
			if i < len(formula) && formula[i] >= 'a' && formula[i] <= 'z' {
				elem += string(formula[i])
				i++
			}
			elements = append(elements, elem)
		} else {
			i++
		}
	}
	return elements
}

// levenshteinDistance calculates the edit distance between two strings.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = minOfThree(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// minOfThree returns the minimum of three integers.
func minOfThree(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// calculateFormulaSimilarity calculates similarity between two formulas.
func calculateFormulaSimilarity(input, target string) float64 {
	input = strings.ToUpper(input)
	target = strings.ToUpper(target)

	distance := levenshteinDistance(input, target)
	maxLen := max(len(input), len(target))

	if maxLen == 0 {
		return 1.0
	}

	similarity := 1.0 - float64(distance)/float64(maxLen)

	// Bonus for same elements
	inputElems := extractElements(input)
	targetElems := extractElements(target)
	commonElems := 0
	for _, e := range inputElems {
		for _, t := range targetElems {
			if e == t {
				commonElems++
				break
			}
		}
	}

	if len(inputElems) > 0 && len(targetElems) > 0 {
		elemBonus := float64(commonElems) / float64(len(targetElems))
		similarity = similarity*0.7 + elemBonus*0.3
	}

	return similarity
}

// roundTo rounds a float to specified decimal places.
func roundTo(val float64, places int) float64 {
	mult := math.Pow(10, float64(places))
	return math.Round(val*mult) / mult
}
