package chemistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFormula_Simple(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		wantComp    Composition
		wantCharge  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "single element H",
			formula:    "H",
			wantComp:   Composition{"H": 1},
			wantCharge: 0,
		},
		{
			name:       "single element O",
			formula:    "O",
			wantComp:   Composition{"O": 1},
			wantCharge: 0,
		},
		{
			name:       "water H2O",
			formula:    "H2O",
			wantComp:   Composition{"H": 2, "O": 1},
			wantCharge: 0,
		},
		{
			name:       "carbon dioxide CO2",
			formula:    "CO2",
			wantComp:   Composition{"C": 1, "O": 2},
			wantCharge: 0,
		},
		{
			name:       "ammonia NH3",
			formula:    "NH3",
			wantComp:   Composition{"N": 1, "H": 3},
			wantCharge: 0,
		},
		{
			name:       "methane CH4",
			formula:    "CH4",
			wantComp:   Composition{"C": 1, "H": 4},
			wantCharge: 0,
		},
		{
			name:       "sodium chloride NaCl",
			formula:    "NaCl",
			wantComp:   Composition{"Na": 1, "Cl": 1},
			wantCharge: 0,
		},
		{
			name:       "sulfuric acid H2SO4",
			formula:    "H2SO4",
			wantComp:   Composition{"H": 2, "S": 1, "O": 4},
			wantCharge: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
			assert.Equal(t, tt.wantCharge, sub.Charge)
		})
	}
}

func TestParseFormula_TwoLetterElements(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
	}{
		{
			name:     "iron Fe",
			formula:  "Fe",
			wantComp: Composition{"Fe": 1},
		},
		{
			name:     "iron oxide Fe2O3",
			formula:  "Fe2O3",
			wantComp: Composition{"Fe": 2, "O": 3},
		},
		{
			name:     "copper sulfate CuSO4",
			formula:  "CuSO4",
			wantComp: Composition{"Cu": 1, "S": 1, "O": 4},
		},
		{
			name:     "sodium hydroxide NaOH",
			formula:  "NaOH",
			wantComp: Composition{"Na": 1, "O": 1, "H": 1},
		},
		{
			name:     "calcium carbonate CaCO3",
			formula:  "CaCO3",
			wantComp: Composition{"Ca": 1, "C": 1, "O": 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_Parentheses(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
	}{
		{
			name:     "calcium hydroxide Ca(OH)2",
			formula:  "Ca(OH)2",
			wantComp: Composition{"Ca": 1, "O": 2, "H": 2},
		},
		{
			name:     "iron(III) sulfate Fe2(SO4)3",
			formula:  "Fe2(SO4)3",
			wantComp: Composition{"Fe": 2, "S": 3, "O": 12},
		},
		{
			name:     "aluminum sulfate Al2(SO4)3",
			formula:  "Al2(SO4)3",
			wantComp: Composition{"Al": 2, "S": 3, "O": 12},
		},
		{
			name:     "magnesium nitrate Mg(NO3)2",
			formula:  "Mg(NO3)2",
			wantComp: Composition{"Mg": 1, "N": 2, "O": 6},
		},
		{
			name:     "ammonium sulfate (NH4)2SO4",
			formula:  "(NH4)2SO4",
			wantComp: Composition{"N": 2, "H": 8, "S": 1, "O": 4},
		},
		{
			name:     "barium phosphate Ba3(PO4)2",
			formula:  "Ba3(PO4)2",
			wantComp: Composition{"Ba": 3, "P": 2, "O": 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_NestedParentheses(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
	}{
		{
			name:     "nested parens 1",
			formula:  "Ca3(PO4)2",
			wantComp: Composition{"Ca": 3, "P": 2, "O": 8},
		},
		{
			name:     "complex nested",
			formula:  "Al2(SO4)3",
			wantComp: Composition{"Al": 2, "S": 3, "O": 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_IonicCharges(t *testing.T) {
	tests := []struct {
		name       string
		formula    string
		wantComp   Composition
		wantCharge int
	}{
		{
			name:       "hydroxide ion OH^-",
			formula:    "OH^-",
			wantComp:   Composition{"O": 1, "H": 1},
			wantCharge: -1,
		},
		{
			name:       "sulfate ion SO4^2-",
			formula:    "SO4^2-",
			wantComp:   Composition{"S": 1, "O": 4},
			wantCharge: -2,
		},
		{
			name:       "ammonium ion NH4^+",
			formula:    "NH4^+",
			wantComp:   Composition{"N": 1, "H": 4},
			wantCharge: 1,
		},
		{
			name:       "iron(III) ion Fe^3+",
			formula:    "Fe^3+",
			wantComp:   Composition{"Fe": 1},
			wantCharge: 3,
		},
		{
			name:       "copper(II) ion Cu^2+",
			formula:    "Cu^2+",
			wantComp:   Composition{"Cu": 1},
			wantCharge: 2,
		},
		{
			name:       "chloride ion Cl^-",
			formula:    "Cl^-",
			wantComp:   Composition{"Cl": 1},
			wantCharge: -1,
		},
		{
			name:       "calcium ion Ca^2+",
			formula:    "Ca^2+",
			wantComp:   Composition{"Ca": 1},
			wantCharge: 2,
		},
		{
			name:       "phosphate ion PO4^3-",
			formula:    "PO4^3-",
			wantComp:   Composition{"P": 1, "O": 4},
			wantCharge: -3,
		},
		{
			name:       "simple sodium ion Na+",
			formula:    "Na+",
			wantComp:   Composition{"Na": 1},
			wantCharge: 1,
		},
		{
			name:       "simple chloride ion Cl-",
			formula:    "Cl-",
			wantComp:   Composition{"Cl": 1},
			wantCharge: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
			assert.Equal(t, tt.wantCharge, sub.Charge)
		})
	}
}

func TestParseFormula_Hydrates(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
	}{
		{
			name:     "copper(II) sulfate pentahydrate CuSO4·5H2O",
			formula:  "CuSO4·5H2O",
			wantComp: Composition{"Cu": 1, "S": 1, "O": 9, "H": 10},
		},
		{
			name:     "calcium sulfate dihydrate CaSO4·2H2O",
			formula:  "CaSO4·2H2O",
			wantComp: Composition{"Ca": 1, "S": 1, "O": 6, "H": 4},
		},
		{
			name:     "magnesium sulfate heptahydrate MgSO4·7H2O",
			formula:  "MgSO4·7H2O",
			wantComp: Composition{"Mg": 1, "S": 1, "O": 11, "H": 14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_StateSymbols(t *testing.T) {
	tests := []struct {
		name      string
		formula   string
		wantState string
	}{
		{
			name:      "solid NaCl(s)",
			formula:   "NaCl(s)",
			wantState: "(s)",
		},
		{
			name:      "liquid H2O(l)",
			formula:   "H2O(l)",
			wantState: "(l)",
		},
		{
			name:      "gas O2(g)",
			formula:   "O2(g)",
			wantState: "(g)",
		},
		{
			name:      "aqueous HCl(aq)",
			formula:   "HCl(aq)",
			wantState: "(aq)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantState, sub.State)
		})
	}
}

func TestParseFormula_Errors(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		errContains string
	}{
		{
			name:        "empty formula",
			formula:     "",
			errContains: "empty formula",
		},
		{
			name:        "unknown element Xx",
			formula:     "Xx",
			errContains: "unknown element",
		},
		{
			name:        "unknown element in compound",
			formula:     "H2XxO4",
			errContains: "unknown element",
		},
		{
			name:        "unmatched opening parenthesis",
			formula:     "Ca(OH",
			errContains: "unmatched parenthesis",
		},
		{
			name:        "unmatched closing parenthesis",
			formula:     "CaOH)",
			errContains: "unmatched closing parenthesis",
		},
		{
			name:        "invalid character",
			formula:     "H2@O",
			errContains: "unexpected character",
		},
		{
			name:        "lowercase start",
			formula:     "h2o",
			errContains: "unexpected character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFormula(tt.formula)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}

func TestParseFormula_Complex(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
		wantErr  bool
	}{
		{
			name:     "iron(III) sulfate",
			formula:  "Fe2(SO4)3",
			wantComp: Composition{"Fe": 2, "S": 3, "O": 12},
		},
		{
			name:     "potassium permanganate",
			formula:  "KMnO4",
			wantComp: Composition{"K": 1, "Mn": 1, "O": 4},
		},
		{
			name:     "potassium dichromate",
			formula:  "K2Cr2O7",
			wantComp: Composition{"K": 2, "Cr": 2, "O": 7},
		},
		{
			name:     "calcium acetate",
			formula:  "Ca(C2H3O2)2",
			wantComp: Composition{"Ca": 1, "C": 4, "H": 6, "O": 4},
		},
		{
			name:     "ammonium phosphate",
			formula:  "(NH4)3PO4",
			wantComp: Composition{"N": 3, "H": 12, "P": 1, "O": 4},
		},
		{
			name:     "magnesium ammonium phosphate",
			formula:  "MgNH4PO4",
			wantComp: Composition{"Mg": 1, "N": 1, "H": 4, "P": 1, "O": 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_HeavyElements(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
	}{
		{
			name:     "uranium hexafluoride",
			formula:  "UF6",
			wantComp: Composition{"U": 1, "F": 6},
		},
		{
			name:     "tungsten carbide",
			formula:  "WC",
			wantComp: Composition{"W": 1, "C": 1},
		},
		{
			name:     "platinum chloride",
			formula:  "PtCl2",
			wantComp: Composition{"Pt": 1, "Cl": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}

func TestParseFormula_MolarMass(t *testing.T) {
	tests := []struct {
		name       string
		formula    string
		wantMassMin float64
		wantMassMax float64
	}{
		{
			name:        "water",
			formula:     "H2O",
			wantMassMin: 18.0,
			wantMassMax: 18.1,
		},
		{
			name:        "carbon dioxide",
			formula:     "CO2",
			wantMassMin: 44.0,
			wantMassMax: 44.1,
		},
		{
			name:        "sulfuric acid",
			formula:     "H2SO4",
			wantMassMin: 98.0,
			wantMassMax: 98.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			mass := sub.MolarMass()
			assert.GreaterOrEqual(t, mass, tt.wantMassMin)
			assert.LessOrEqual(t, mass, tt.wantMassMax)
		})
	}
}

func TestParseFormula_String(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantStr  string
	}{
		{
			name:    "water",
			formula: "H2O",
			wantStr: "H2O",
		},
		{
			name:    "hydroxide ion",
			formula: "OH^-",
			wantStr: "OH^-",
		},
		{
			name:    "sulfate ion",
			formula: "SO4^2-",
			wantStr: "SO4^2-",
		},
		{
			name:    "ammonium ion",
			formula: "NH4^+",
			wantStr: "NH4^+",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			require.NoError(t, err)
			assert.NotEmpty(t, sub.String())
		})
	}
}

func TestComposition_Copy(t *testing.T) {
	orig := Composition{"H": 2, "O": 1}
	copy := orig.Copy()

	assert.Equal(t, orig, copy)

	// Modify copy
	copy["H"] = 10
	assert.NotEqual(t, orig, copy)
	assert.Equal(t, Composition{"H": 2, "O": 1}, orig)
}

func TestComposition_Add(t *testing.T) {
	c1 := Composition{"H": 2, "O": 1}
	c2 := Composition{"H": 1, "O": 2}

	c1.Add(c2)

	assert.Equal(t, Composition{"H": 3, "O": 3}, c1)
}

func TestComposition_Multiply(t *testing.T) {
	c := Composition{"H": 2, "O": 1}
	c.Multiply(3)

	assert.Equal(t, Composition{"H": 6, "O": 3}, c)
}

func TestComposition_TotalAtoms(t *testing.T) {
	c := Composition{"H": 2, "O": 1}
	assert.Equal(t, 3, c.TotalAtoms())

	c2 := Composition{"C": 6, "H": 12, "O": 6}
	assert.Equal(t, 24, c2.TotalAtoms())
}

func TestComposition_Equals(t *testing.T) {
	c1 := Composition{"H": 2, "O": 1}
	c2 := Composition{"H": 2, "O": 1}
	c3 := Composition{"H": 1, "O": 2}

	assert.True(t, c1.Equals(c2))
	assert.False(t, c1.Equals(c3))

	// Different keys
	c4 := Composition{"H": 2, "O": 1, "C": 0}
	assert.False(t, c1.Equals(c4))
}

// TestParseFormula_EdgeCases tests various edge cases
func TestParseFormula_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		formula  string
		wantComp Composition
		wantErr  bool
	}{
		{
			name:    "single atom with large count",
			formula: "C100",
			wantComp: Composition{"C": 100},
		},
		{
			name:    "multiple same elements",
			formula: "CH3CH2OH",
			wantComp: Composition{"C": 2, "H": 6, "O": 1},
		},
		{
			name:    "element appearing multiple times",
			formula: "H2O2",
			wantComp: Composition{"H": 2, "O": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := ParseFormula(tt.formula)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantComp, sub.Composition)
		})
	}
}
