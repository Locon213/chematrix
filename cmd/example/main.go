// Example usage of the chematrix library.
// Run with: go run cmd/example/main.go
package main

import (
	"fmt"
	"log"

	"github.com/Locon213/chematrix/chemistry"
)

func main() {
	fmt.Println("=== chematrix Example ===")
	fmt.Println()

	// Example 1: Parse a simple formula
	fmt.Println("1. Parse a simple formula (H2SO4):")
	substance, err := chemistry.ParseFormula("H2SO4")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Formula: %s\n", substance.Formula)
	fmt.Printf("   Composition: %v\n", substance.Composition)
	fmt.Printf("   Molar Mass: %.2f g/mol\n\n", substance.MolarMass())

	// Example 2: Parse an ionic compound with charge
	fmt.Println("2. Parse an ionic compound (SO4^2-):")
	sulfate, err := chemistry.ParseFormula("SO4^2-")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Formula: %s\n", sulfate.Formula)
	fmt.Printf("   Composition: %v\n", sulfate.Composition)
	fmt.Printf("   Charge: %d\n\n", sulfate.Charge)

	// Example 3: Parse a complex formula with parentheses
	fmt.Println("3. Parse a complex formula (Fe2(SO4)3):")
	complex, err := chemistry.ParseFormula("Fe2(SO4)3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Formula: %s\n", complex.Formula)
	fmt.Printf("   Composition: %v\n", complex.Composition)
	fmt.Printf("   Total atoms: %d\n\n", complex.Composition.TotalAtoms())

	// Example 4: Parse a hydrate
	fmt.Println("4. Parse a hydrate (CuSO4·5H2O):")
	hydrate, err := chemistry.ParseFormula("CuSO4·5H2O")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Formula: %s\n", hydrate.Formula)
	fmt.Printf("   Composition: %v\n", hydrate.Composition)
	fmt.Printf("   Molar Mass: %.2f g/mol\n\n", hydrate.MolarMass())

	// Example 5: Parse with state symbol
	fmt.Println("5. Parse with state symbol (NaCl(aq)):")
	aqueous, err := chemistry.ParseFormula("NaCl(aq)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Formula: %s\n", aqueous.Formula)
	fmt.Printf("   State: %s\n\n", aqueous.State)

	// Example 6: Access element data
	fmt.Println("6. Access element data (Fe):")
	iron := chemistry.GetElement("Fe")
	if iron != nil {
		fmt.Printf("   Name: %s\n", iron.Name)
		fmt.Printf("   Atomic Number: %d\n", iron.AtomicNumber)
		fmt.Printf("   Atomic Mass: %.2f\n", iron.AtomicMass)
		fmt.Printf("   Oxidation States: %v\n", iron.OxidationStates)
		fmt.Printf("   Group: %d\n\n", iron.Group)
	}

	// Example 7: Check if element exists
	fmt.Println("7. Check if elements exist:")
	elements := []string{"Fe", "O", "Xx", "Og"}
	for _, sym := range elements {
		exists := chemistry.IsValidElement(sym)
		fmt.Printf("   %s: %v\n", sym, exists)
	}
	fmt.Println()

	// Example 8: Get all elements count
	fmt.Println("8. Periodic table statistics:")
	allElements := chemistry.GetAllElements()
	fmt.Printf("   Total elements: %d\n\n", len(allElements))

	// Example 9: Composition operations
	fmt.Println("9. Composition operations:")
	comp1 := chemistry.Composition{"H": 2, "O": 1}
	comp2 := comp1.Copy()
	comp2.Multiply(2)
	fmt.Printf("   Original: %v\n", comp1)
	fmt.Printf("   Multiplied by 2: %v\n", comp2)
	fmt.Printf("   Total atoms in original: %d\n\n", comp1.TotalAtoms())

	// Example 10: Balance a chemical equation
	fmt.Println("10. Balance a chemical equation (H2 + O2 -> H2O):")
	result, err := chemistry.BalanceEquation("H2 + O2 -> H2O")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("    Balanced: %v\n", result.IsBalanced)
	fmt.Printf("    Coefficients: %v\n", result.Coefficients)
	fmt.Println("    Steps:")
	for i, step := range result.Steps {
		if step.Equation != "" {
			fmt.Printf("      %d. %s: %s\n", i+1, step.Description, step.Equation)
		} else {
			fmt.Printf("      %d. %s\n", i+1, step.Description)
		}
	}
	fmt.Println()

	// Example 11: Balance another equation
	fmt.Println("11. Balance another equation (CH4 + O2 -> CO2 + H2O):")
	result2, err := chemistry.BalanceEquation("CH4 + O2 -> CO2 + H2O")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("    Balanced: %v\n", result2.IsBalanced)
	fmt.Printf("    Coefficients: %v\n\n", result2.Coefficients)

	// Example 12: Parse and balance ammonia synthesis
	fmt.Println("12. Balance ammonia synthesis (N2 + H2 -> NH3):")
	result3, err := chemistry.BalanceEquation("N2 + H2 -> NH3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("    Balanced: %v\n", result3.IsBalanced)
	fmt.Printf("    Result: ")
	for i, comp := range result3.Reaction.Reactants {
		if i > 0 {
			fmt.Print(" + ")
		}
		coef := result3.Coefficients[i]
		if coef > 1 {
			fmt.Printf("%d%s", coef, comp.Substance.Formula)
		} else {
			fmt.Print(comp.Substance.Formula)
		}
	}
	fmt.Print(" → ")
	for i, comp := range result3.Reaction.Products {
		if i > 0 {
			fmt.Print(" + ")
		}
		coef := result3.Coefficients[len(result3.Reaction.Reactants)+i]
		if coef > 1 {
			fmt.Printf("%d%s", coef, comp.Substance.Formula)
		} else {
			fmt.Print(comp.Substance.Formula)
		}
	}
	fmt.Println()
}
