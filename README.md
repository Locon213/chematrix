# chematrix

[![Go Reference](https://pkg.go.dev/badge/github.com/Locon213/chematrix.svg)](https://pkg.go.dev/github.com/Locon213/chematrix)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.26.1+-blue.svg)](https://go.dev)

A pure-Go library for chemical computations, designed for educational applications (like "PhotoMath for Chemistry").

## Features

- **Parse and validate chemical formulas** (e.g., `Fe2(SO4)3`, `OH^-`)
- **Handle ionic charges** (e.g., `SO4^2-`, `NH4^+`, `Fe^3+`)
- **Support hydrates** (e.g., `CuSO4·5H2O`)
- **State symbols** (e.g., `NaCl(s)`, `HCl(aq)`)
- **Molar mass calculation**
- **Balance chemical equations** (coming soon)
- **Redox reactions** (coming soon)
- **Step-by-step explanations** for educational purposes (coming soon)

## Installation

```bash
go get github.com/Locon213/chematrix
```

## Quick Start

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

### Parse Ionic Compounds

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

### Parse Complex Formulas

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
    fmt.Printf("Name: %s\n", iron.Name)           // Output: Name: Iron
    fmt.Printf("Atomic Number: %d\n", iron.AtomicNumber) // Output: Atomic Number: 26
    fmt.Printf("Atomic Mass: %.2f\n", iron.AtomicMass)   // Output: Atomic Mass: 55.85
    fmt.Printf("Oxidation States: %v\n", iron.OxidationStates) // Output: Oxidation States: [2 3]
}

// Check if element exists
if chemistry.IsValidElement("Uuo") {
    fmt.Println("Element exists")
}
```

## API Overview

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

#### `GetElement(symbol string) *Element`
Returns an Element by its symbol (case-sensitive).

#### `GetAllElements() []Element`
Returns a slice of all 118 elements in the periodic table.

## Supported Features

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
| Equation balancing | Coming in v0.2 | 🔄 |
| Redox reactions | Coming in v0.3 | 🔄 |

## Examples

### Interactive CLI

Chematrix includes a beautiful interactive TUI (Text User Interface) for exploring chemical formulas:

```bash
go run cmd/chematrix/main.go
```

Features:
- 🧪 Parse chemical formulas with instant feedback
- ⚛️  Look up element information from the periodic table
- 📊 Calculate molar mass with percentage breakdown
- ⚖️  Compare two formulas side by side

### Simple Formula Parsing

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

## Testing

Run tests with:

```bash
go test ./chemistry/... -v
```

The parser has comprehensive test coverage including:
- Simple formulas
- Complex nested structures
- Ionic charges
- Hydrates
- Error cases

## Roadmap

### Phase 1: Core Parser & Types (Current)
- [x] Define `Element`, `Substance`, `Composition` types
- [x] Implement `ParseFormula` with regex/stack-based parser
- [x] Add periodic table data (118 elements)
- [x] Write 50+ unit tests for parser edge cases

### Phase 2: Stoichiometry Engine (Next)
- [ ] Integrate `gonum.org/v1/gonum/mat`
- [ ] Implement matrix builder from reaction components
- [ ] Add null-space solver for coefficient extraction
- [ ] Implement equation balancing with explanations

### Phase 3: Redox & Ionic Support
- [ ] Add `HalfReaction` type and electron logic
- [ ] Embed standard reduction potentials table
- [ ] Implement ionic equation simplification
- [ ] Add solubility rules lookup

### Phase 4: Polish & Release
- [ ] Benchmark performance
- [ ] Add CLI demo
- [ ] Complete documentation
- [ ] Tag v1.0.0 release

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

**Locon213**

## Acknowledgments

- Periodic table data sourced from IUPAC standard atomic weights
- Built with ❤️ for educational applications
