# Chematrix

[![Go Reference](https://pkg.go.dev/badge/github.com/Locon213/chematrix.svg)](https://pkg.go.dev/github.com/Locon213/chematrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/Locon213/chematrix)](https://goreportcard.com/report/github.com/Locon213/chematrix)
[![License](https://img.shields.io/github/license/Locon213/chematrix?color=blue&style=flat-square)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Locon213/chematrix?color=00ADD8&logo=go&style=flat-square)](https://go.dev)
[![Build & Release](https://github.com/Locon213/chematrix/actions/workflows/release.yml/badge.svg)](https://github.com/Locon213/chematrix/actions/workflows/release.yml)
[![Tests](https://img.shields.io/github/actions/workflow/status/Locon213/chematrix/release.yml?branch=main&label=tests&style=flat-square)](https://github.com/Locon213/chematrix/actions)

**Chematrix** is a pure Go library for chemical computations, designed for educational applications (like "PhotoMath for Chemistry").

## 🚀 Features

### Formula Processing
- ✅ **Parse and validate chemical formulas** (e.g., `Fe2(SO4)3`, `OH^-`)
- ✅ **Handle ionic charges** (e.g., `SO4^2-`, `NH4^+`, `Fe^3+`)
- ✅ **Support hydrates** (e.g., `CuSO4·5H2O`)
- ✅ **State symbols** (e.g., `NaCl(s)`, `HCl(aq)`)
- ✅ **Molar mass calculation**
- ✅ **Smart formula suggestions** based on chemical rules

### Chemical Equations
- ✅ **Balance equations** using matrix methods (SVD)
- ✅ **Step-by-step explanations** for educational purposes
- ✅ **Reaction type detection** (combustion, synthesis, decomposition, acid-base, replacement)
- ✅ **Equation correction suggestions**

### Periodic Table
- ✅ **All 118 elements** with complete data
- ✅ **Oxidation states** for each element
- ✅ **Element groups** for chemical property analysis
- ✅ **Search by symbol and name**

### Advanced Suggestion System
- 🔮 **Valence-based generation** — uses oxidation states to suggest valid compounds
- 🔮 **Typo correction** — Levenshtein distance + chemical rules
- 🔮 **Case fixing** — "CL" → "Cl₂"
- 🔮 **0/O confusion detection** — "H20" → "H₂O"
- 🔮 **Diatomic elements** — "H" → "H₂"
- 🔮 **Similar elements** — suggestions from the same group

## 📦 Installation

```bash
go get github.com/Locon213/chematrix
```

## 🚀 Quick Start

### Parse a Chemical Formula

```go
package main

import (
    "fmt"
    "log"

    "github.com/Locon213/chematrix/chemistry"
)

func main() {
    // Parse a simple formula
    substance, err := chemistry.ParseFormula("H2SO4")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Formula: %s\n", substance.Formula)
    fmt.Printf("Composition: %v\n", substance.Composition)
    // Output: Composition: map[H:2 O:4 S:1]

    fmt.Printf("Molar Mass: %.2f g/mol\n", substance.MolarMass())
    // Output: Molar Mass: 98.08 g/mol
}
```

### Ionic Compounds

```go
// Parse an ion with charge
sulfate, err := chemistry.ParseFormula("SO4^2-")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Charge: %d\n", sulfate.Charge) // Output: Charge: -2

// Simple ionic format
ammonium, err := chemistry.ParseFormula("NH4^+")
fmt.Printf("Charge: %d\n", ammonium.Charge) // Output: Charge: 1
```

### Complex Formulas

```go
// Nested parentheses
complex, err := chemistry.ParseFormula("Fe2(SO4)3")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Composition: %v\n", complex.Composition)
// Output: Composition: map[Fe:2 S:3 O:12]

// Hydrates
hydrate, err := chemistry.ParseFormula("CuSO4·5H2O")
fmt.Printf("Composition: %v\n", hydrate.Composition)
// Output: Composition: map[Cu:1 S:1 O:9 H:10]

// State symbols
aqueous, err := chemistry.ParseFormula("NaCl(aq)")
fmt.Printf("State: %s\n", aqueous.State) // Output: State: (aq)
```

### Access Element Data

```go
// Get element by symbol
iron := chemistry.GetElement("Fe")
if iron != nil {
    fmt.Printf("Name: %s\n", iron.Name)           
    fmt.Printf("Atomic Number: %d\n", iron.AtomicNumber) 
    fmt.Printf("Atomic Mass: %.2f\n", iron.AtomicMass)   
    fmt.Printf("Oxidation States: %v\n", iron.OxidationStates) 
}

// Check if element exists
if chemistry.IsValidElement("Uuo") {
    fmt.Println("Element exists")
}
```

### Balance Chemical Equations

```go
package main

import (
    "fmt"
    "log"

    "github.com/Locon213/chematrix/chemistry"
)

func main() {
    // Balance a chemical equation
    result, err := chemistry.BalanceEquation("H2 + O2 -> H2O")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Balanced: %v\n", result.IsBalanced)
    // Output: Balanced: true

    fmt.Printf("Coefficients: %v\n", result.Coefficients)
    // Output: Coefficients: [2 1 2]

    // Print step-by-step explanation
    for _, step := range result.Steps {
        fmt.Printf("%s: %s\n", step.Description, step.Equation)
    }
}
```

### Smart Formula Suggestions

```go
package main

import (
    "fmt"

    "github.com/Locon213/chematrix/chemistry"
)

func main() {
    // Example 1: Typo correction
    suggestions := chemistry.SuggestFormulas("HO")
    for _, s := range suggestions {
        fmt.Printf("• %s (%.0f%%) — %s\n", s.Formula, s.Similarity*100, s.Reason)
    }
    // Output:
    // • H2O (95%) — Replaced zero (0) with oxygen (O)
    // • H2O2 (85%) — Similar formula structure

    // Example 2: Case fixing
    suggestions = chemistry.SuggestFormulas("CL")
    for _, s := range suggestions {
        fmt.Printf("• %s (%.0f%%) — %s\n", s.Formula, s.Similarity*100, s.Reason)
    }
    // Output: • Cl2 (95%) — Fixed element symbol capitalization

    // Example 3: 0/O confusion
    suggestions = chemistry.SuggestFormulas("H20")
    for _, s := range suggestions {
        fmt.Printf("• %s (%.0f%%) — %s\n", s.Formula, s.Similarity*100, s.Reason)
    }
    // Output: • H2O (95%) — Replaced zero (0) with oxygen (O)

    // Example 4: Diatomic elements
    suggestions = chemistry.SuggestFormulas("H")
    for _, s := range suggestions {
        fmt.Printf("• %s (%.0f%%) — %s\n", s.Formula, s.Similarity*100, s.Reason)
    }
    // Output: • H2 (85%) — Hydrogen is commonly found as diatomic H2

    // Example 5: Valence rules
    suggestions = chemistry.SuggestFormulas("NaO")
    for _, s := range suggestions {
        fmt.Printf("• %s (%.0f%%) — %s\n", s.Formula, s.Similarity*100, s.Reason)
    }
    // Output: • Na2O (70%) — Valence-compatible compound: Na(+1) + O(-2)
}
```

### Equation Suggestions

```go
package main

import (
    "fmt"

    "github.com/Locon213/chematrix/chemistry"
)

func main() {
    // Fix arrow format
    suggestions := chemistry.SuggestEquations("H2 + O2 = H2O")
    for _, s := range suggestions {
        fmt.Printf("• %s\n  Reason: %s\n  Balanced: %v\n", s.Suggested, s.Reason, s.IsBalanced)
    }

    // Reaction type analysis
    suggestions = chemistry.SuggestEquations("CH4 + O2 -> CO2 + H2O")
    for _, s := range suggestions {
        fmt.Printf("• %s\n  Reason: %s\n", s.Suggested, s.Reason)
    }
    // Output: • 2CH4 + 4O2 -> 2CO2 + 4H2O
    //         Reason: Balanced combustion reaction
}
```

## 📚 API Overview

### Types

#### `chemistry.Element`
Represents a chemical element from the periodic table.

```go
type Element struct {
    Symbol          string  // e.g., "Fe", "O"
    Name            string  // e.g., "Iron", "Oxygen"
    AtomicNumber    int     // e.g., 26
    AtomicMass      float64 // atomic weight in u
    OxidationStates []int   // common oxidation states
    Group           int     // periodic table group
}
```

#### `chemistry.Substance`
Represents a chemical compound or ion.

```go
type Substance struct {
    Formula     string            // original input, normalized
    Composition Composition       // map[element]count
    Charge      int               // ionic charge: -2, +1, 0, etc.
    State       string            // optional: "(s)", "(aq)", "(g)"
}
```

#### `chemistry.Composition`
A map of element symbols to atom counts.

```go
type Composition map[string]int
```

#### `chemistry.FormulaSuggestion`
A suggested formula with similarity score.

```go
type FormulaSuggestion struct {
    Formula    string     // suggested formula
    Similarity float64    // similarity score (0-1)
    Reason     string     // why this suggestion was made
    Substance  *Substance // parsed substance
}
```

#### `chemistry.EquationSuggestion`
A suggested equation correction.

```go
type EquationSuggestion struct {
    Original     string // original equation
    Suggested    string // suggested correction
    Reason       string // what was fixed
    IsBalanced   bool   // whether the suggestion is balanced
    Coefficients []int  // balancing coefficients
}
```

#### `chemistry.BalancedReaction`
A successfully balanced chemical equation.

```go
type BalancedReaction struct {
    Reaction     *Reaction
    Coefficients []int
    Steps        []BalanceStep
    IsBalanced   bool
}
```

### Functions

#### `ParseFormula(formula string) (*Substance, error)`
Converts a chemical formula string into a Substance.

**Supported formats:**
- Elements: `"H"`, `"Fe"`, `"Uuo"`
- Numbers: `"H2O"`, `"Fe2(SO4)3"`
- Parentheses: `"Ca(OH)2"`, `"Al2(SO4)3"`
- Ionic charges: `"SO4^2-"`, `"NH4^+"`, `"Fe^3+"`
- Hydrates: `"CuSO4·5H2O"`
- State symbols: `"NaCl(s)"`, `"HCl(aq)"`

#### `SuggestFormulas(input string) []FormulaSuggestion`
Generates formula suggestions based on input.

**Rules used:**
- Typo correction (Levenshtein distance)
- Element case fixing
- 0/O confusion detection
- Valence rules
- Diatomic elements
- Similar elements (same group)

#### `SuggestEquations(input string) []EquationSuggestion`
Generates suggestions for chemical equations.

**Analyzed reaction types:**
- Combustion
- Synthesis
- Decomposition
- Acid-base
- Single replacement
- Double replacement

#### `BalanceEquation(equation string) (*BalancedReaction, error)`
Balances a chemical equation.

#### `GetElement(symbol string) *Element`
Returns an Element by its symbol (case-sensitive).

#### `GetAllElements() []Element`
Returns all 118 elements in the periodic table.

#### `IsValidElement(symbol string) bool`
Checks if a symbol corresponds to a known element.

## 🖥️ Interactive CLI

Chematrix includes a beautiful interactive TUI (Text User Interface):

```bash
go run cmd/chematrix/main.go
```

**Features:**
- 🧪 Parse chemical formulas with instant feedback
- ⚛️  Look up element information from the periodic table
- 📊 Calculate molar mass with percentage breakdown
- ⚖️  Compare two formulas side by side
- ⚗️  Balance chemical equations
- 💡 **Smart formula suggestions**

## 📋 Supported Features

| Feature | Example | Status |
|---------|---------|--------|
| Simple formulas | `H2O`, `CO2` | ✅ |
| Two-letter elements | `Fe`, `Cu`, `Na` | ✅ |
| Parentheses | `Ca(OH)2` | ✅ |
| Nested parentheses | `Al2(SO4)3` | ✅ |
| Ionic charges (^) | `SO4^2-`, `Fe^3+` | ✅ |
| Simple ions | `Na+`, `Cl-` | ✅ |
| Hydrates | `CuSO4·5H2O` | ✅ |
| State symbols | `NaCl(s)`, `HCl(aq)` | ✅ |
| Molar mass | `substance.MolarMass()` | ✅ |
| All 118 elements | Periodic table data | ✅ |
| Equation balancing | `BalanceEquation("H2 + O2 -> H2O")` | ✅ |
| Smart suggestions | `SuggestFormulas("HO")` | ✅ |
| Reaction types | Automatic detection | ✅ |

## 🧪 Testing

Run tests with:

```bash
go test ./chemistry/... -v
```

The library has comprehensive test coverage including:
- Simple formulas
- Complex nested structures
- Ionic charges
- Hydrates
- Equation balancing
- Suggestion system
- Typo correction

## 📖 Usage Examples

### Example 1: User Input Validation

```go
func validateUserInput(input string) (string, error) {
    // First try to parse as-is
    if _, err := chemistry.ParseFormula(input); err == nil {
        return input, nil
    }
    
    // If failed, offer suggestions
    suggestions := chemistry.SuggestFormulas(input)
    if len(suggestions) > 0 {
        // Return the most similar formula
        return suggestions[0].Formula, nil
    }
    
    return "", fmt.Errorf("failed to recognize formula")
}
```

### Example 2: Educational Application

```go
func explainMistake(input, correct string) string {
    suggestions := chemistry.SuggestFormulas(input)
    for _, s := range suggestions {
        if s.Formula == correct {
            return fmt.Sprintf("Mistake: %s", s.Reason)
        }
    }
    return "Please check the formula spelling"
}
```

### Example 3: Generate Similar Compounds

```go
func findSimilarCompounds(formula string) []string {
    suggestions := chemistry.SuggestFormulas(formula)
    var similar []string
    for _, s := range suggestions {
        if s.Similarity > 0.5 {
            similar = append(similar, s.Formula)
        }
    }
    return similar
}
```

## 🛣️ Roadmap

### Phase 1: Core Parser & Types (Completed ✅)
- [x] Define `Element`, `Substance`, `Composition` types
- [x] Implement `ParseFormula` with regex/stack-based parser
- [x] Add periodic table data (118 elements)
- [x] Write 50+ unit tests for parser edge cases

### Phase 2: Stoichiometry Engine (Completed ✅)
- [x] Integrate `gonum.org/v1/gonum/mat`
- [x] Matrix solver for coefficient extraction
- [x] Equation balancing with explanations
- [x] Reaction type detection

### Phase 3: Smart Suggestions (Completed ✅)
- [x] Valence-based generation
- [x] Typo correction (Levenshtein)
- [x] Element case fixing
- [x] 0/O confusion detection
- [x] Diatomic elements
- [x] Similar elements (groups)

### Phase 4: Redox & Ionic Support (Planned)
- [ ] Add `HalfReaction` type and electron logic
- [ ] Embed standard reduction potentials table
- [ ] Implement ionic equation simplification
- [ ] Add solubility rules lookup

### Phase 5: Polish & Release (In Progress)
- [x] Benchmark performance
- [x] CLI demo with TUI
- [x] Complete documentation
- [ ] Tag v1.0.0 release

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 👤 Author

**Locon213**

## 🙏 Acknowledgments

- Periodic table data sourced from IUPAC standard atomic weights
- Built with ❤️ for educational applications
