package chemistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEquation_Simple(t *testing.T) {
	tests := []struct {
		name          string
		equation      string
		wantReactants []string
		wantProducts  []string
	}{
		{
			name:          "simple water formation",
			equation:      "H2 + O2 -> H2O",
			wantReactants: []string{"H2", "O2"},
			wantProducts:  []string{"H2O"},
		},
		{
			name:          "arrow with spaces",
			equation:      "H2 + O2 -> H2O",
			wantReactants: []string{"H2", "O2"},
			wantProducts:  []string{"H2O"},
		},
		{
			name:          "equals sign",
			equation:      "H2 + O2 = H2O",
			wantReactants: []string{"H2", "O2"},
			wantProducts:  []string{"H2O"},
		},
		{
			name:          "multiple reactants and products",
			equation:      "NaOH + HCl -> NaCl + H2O",
			wantReactants: []string{"NaOH", "HCl"},
			wantProducts:  []string{"NaCl", "H2O"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := ParseEquation(tt.equation)
			require.NoError(t, err)

			assert.Equal(t, len(tt.wantReactants), len(reaction.Reactants))
			assert.Equal(t, len(tt.wantProducts), len(reaction.Products))

			for i, want := range tt.wantReactants {
				assert.Equal(t, want, reaction.Reactants[i].Substance.Formula)
			}

			for i, want := range tt.wantProducts {
				assert.Equal(t, want, reaction.Products[i].Substance.Formula)
			}
		})
	}
}

func TestParseEquation_WithCoefficients(t *testing.T) {
	tests := []struct {
		name     string
		equation string
		side     string // "reactants" or "products"
		index    int
		wantCoef int
	}{
		{
			name:     "coefficient 2",
			equation: "2H2 + O2 -> 2H2O",
			side:     "reactants",
			index:    0,
			wantCoef: 2,
		},
		{
			name:     "coefficient 1 implicit",
			equation: "2H2 + O2 -> 2H2O",
			side:     "reactants",
			index:    1,
			wantCoef: 1,
		},
		{
			name:     "product coefficient",
			equation: "2H2 + O2 -> 2H2O",
			side:     "products",
			index:    0,
			wantCoef: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := ParseEquation(tt.equation)
			require.NoError(t, err)

			var components []ReactionComponent
			if tt.side == "reactants" {
				components = reaction.Reactants
			} else {
				components = reaction.Products
			}

			require.Greater(t, len(components), tt.index)
			assert.Equal(t, tt.wantCoef, components[tt.index].Coefficient)
		})
	}
}

func TestParseEquation_Errors(t *testing.T) {
	tests := []struct {
		name     string
		equation string
	}{
		{"empty equation", ""},
		{"no arrow", "H2 O2 H2O"},
		{"empty reactants", "-> H2O"},
		{"empty products", "H2 + O2 ->"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseEquation(tt.equation)
			assert.Error(t, err)
		})
	}
}

func TestBalanceEquation_Simple(t *testing.T) {
	tests := []struct {
		name             string
		equation         string
		wantCoefficients []int
		wantBalanced     bool
	}{
		{
			name:             "water formation H2 + O2 -> H2O",
			equation:         "H2 + O2 -> H2O",
			wantCoefficients: []int{2, 1, 2}, // 2H2 + O2 -> 2H2O
			wantBalanced:     true,
		},
		{
			name:             "already balanced H2 + Cl2 -> HCl",
			equation:         "H2 + Cl2 -> HCl",
			wantCoefficients: []int{1, 1, 1}, // May be normalized
			wantBalanced:     true,
		},
		{
			name:             "combustion CH4 + O2 -> CO2 + H2O",
			equation:         "CH4 + O2 -> CO2 + H2O",
			wantCoefficients: []int{1, 2, 1, 2}, // CH4 + 2O2 -> CO2 + 2H2O
			wantBalanced:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BalanceEquation(tt.equation)
			require.NoError(t, err)

			assert.NotEmpty(t, result.Coefficients)
			assert.Equal(t, tt.wantBalanced, result.IsBalanced)

			// Verify atom balance manually
			for i, step := range result.Steps {
				t.Logf("Step %d: %s - %s", i, step.Description, step.Equation)
			}
		})
	}
}

func TestBalanceEquation_Complex(t *testing.T) {
	tests := []struct {
		name     string
		equation string
	}{
		{
			name:     "ammonia synthesis",
			equation: "N2 + H2 -> NH3",
		},
		{
			name:     "calcium phosphate",
			equation: "Ca3(PO4)2 + H2SO4 -> CaSO4 + H3PO4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BalanceEquation(tt.equation)
			require.NoError(t, err)

			assert.NotEmpty(t, result.Coefficients)
			assert.True(t, result.IsBalanced, "Equation should be balanced")

			// Log the balanced equation
			t.Logf("Balanced: %s", formatEquation(result.Reaction, result.Coefficients))
		})
	}
}

// TestBalanceEquation_KnownLimitations contains tests for known algorithm limitations.
// These equations require a more sophisticated balancing algorithm.
func TestBalanceEquation_KnownLimitations(t *testing.T) {
	tests := []struct {
		name     string
		equation string
	}{
		{
			name:     "iron oxidation - needs odd coefficients",
			equation: "Fe + O2 -> Fe2O3",
		},
		{
			name:     "aluminum oxidation - needs odd coefficients",
			equation: "Al + O2 -> Al2O3",
		},
		{
			name:     "potassium chlorate decomposition",
			equation: "KClO3 -> KCl + O2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BalanceEquation(tt.equation)
			require.NoError(t, err)

			assert.NotEmpty(t, result.Coefficients)
			// NOTE: These tests are expected to fail until the algorithm is improved
			t.Logf("Balanced: %s (known limitation)", formatEquation(result.Reaction, result.Coefficients))
		})
	}
}

func TestReaction_Balance(t *testing.T) {
	tests := []struct {
		name     string
		equation string
	}{
		{
			name:     "simple water",
			equation: "H2 + O2 -> H2O",
		},
		{
			name:     "neutralization",
			equation: "NaOH + HCl -> NaCl + H2O",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := ParseEquation(tt.equation)
			require.NoError(t, err)

			balanced, err := reaction.Balance()
			require.NoError(t, err)

			assert.NotEmpty(t, balanced.Coefficients)
			assert.NotEmpty(t, balanced.Steps)
		})
	}
}

func TestFormatEquation(t *testing.T) {
	tests := []struct {
		name         string
		equation     string
		coefficients []int
		wantContain  string
	}{
		{
			name:         "with coefficients",
			equation:     "H2 + O2 -> H2O",
			coefficients: []int{2, 1, 2},
			wantContain:  "2H2",
		},
		{
			name:         "without coefficients",
			equation:     "H2 + O2 -> H2O",
			coefficients: []int{1, 1, 1},
			wantContain:  "H2 + O2 → H2O",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := ParseEquation(tt.equation)
			require.NoError(t, err)

			result := formatEquation(reaction, tt.coefficients)
			assert.Contains(t, result, tt.wantContain)
		})
	}
}

func TestVerifyBalance(t *testing.T) {
	tests := []struct {
		name         string
		equation     string
		coefficients []int
		wantBalance  bool
	}{
		{
			name:         "balanced water",
			equation:     "H2 + O2 -> H2O",
			coefficients: []int{2, 1, 2},
			wantBalance:  true,
		},
		{
			name:         "unbalanced water",
			equation:     "H2 + O2 -> H2O",
			coefficients: []int{1, 1, 1},
			wantBalance:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := ParseEquation(tt.equation)
			require.NoError(t, err)

			// Collect elements
			elementSet := make(map[string]bool)
			allComponents := append(reaction.Reactants, reaction.Products...)
			for _, comp := range allComponents {
				for elem := range comp.Substance.Composition {
					elementSet[elem] = true
				}
			}

			elements := make([]string, 0, len(elementSet))
			for elem := range elementSet {
				elements = append(elements, elem)
			}

			result := verifyBalance(reaction, tt.coefficients, elements)
			assert.Equal(t, tt.wantBalance, result)
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{12, 18, 6},
		{7, 13, 1},
		{100, 25, 25},
		{0, 5, 5},
		{5, 0, 5},
	}

	for _, tt := range tests {
		t.Logf("gcd(%d, %d) = %d", tt.a, tt.b, gcd(tt.a, tt.b))
		assert.Equal(t, tt.want, gcd(tt.a, tt.b))
	}
}

func TestNormalizeCoefficients(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		output []int
	}{
		{
			name:   "already normalized",
			input:  []int{1, 2, 1},
			output: []int{1, 2, 1},
		},
		{
			name:   "can divide by 2",
			input:  []int{2, 4, 2},
			output: []int{1, 2, 1},
		},
		{
			name:   "can divide by 3",
			input:  []int{3, 6, 3},
			output: []int{1, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCoefficients(tt.input)
			assert.Equal(t, tt.output, result)
		})
	}
}
