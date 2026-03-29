// Package chemistry provides core types for representing chemical elements and compounds.
package chemistry

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// ReactionComponent represents a single substance in a chemical equation with its coefficient.
type ReactionComponent struct {
	Substance   *Substance
	Coefficient int // stoichiometric coefficient (1 if not specified)
}

// Reaction represents a chemical equation with reactants and products.
type Reaction struct {
	Reactants []ReactionComponent
	Products  []ReactionComponent
	Raw       string // original equation string
}

// BalanceStep represents a step in the balancing explanation.
type BalanceStep struct {
	Description string
	Equation    string
}

// BalancedReaction represents a successfully balanced chemical equation.
type BalancedReaction struct {
	Reaction     *Reaction
	Coefficients []int // coefficients for all components (reactants + products)
	Steps        []BalanceStep
	IsBalanced   bool
}

// ParseEquation parses a chemical equation string into a Reaction.
// Example: "H2 + O2 -> H2O" or "H2 + O2 = H2O"
func ParseEquation(equation string) (*Reaction, error) {
	equation = strings.TrimSpace(equation)
	if equation == "" {
		return nil, fmt.Errorf("parse error: empty equation")
	}

	// Split by arrow or equals sign
	splitPattern := regexp.MustCompile(`\s*(?:->|→|=)\s*`)
	parts := splitPattern.Split(equation, -1)

	if len(parts) != 2 {
		return nil, fmt.Errorf("parse error: invalid equation format, expected 'reactants -> products'")
	}

	reactants, err := parseEquationSide(parts[0])
	if err != nil {
		return nil, fmt.Errorf("parse error in reactants: %w", err)
	}

	products, err := parseEquationSide(parts[1])
	if err != nil {
		return nil, fmt.Errorf("parse error in products: %w", err)
	}

	return &Reaction{
		Reactants: reactants,
		Products:  products,
		Raw:       equation,
	}, nil
}

// parseEquationSide parses one side of a chemical equation.
func parseEquationSide(side string) ([]ReactionComponent, error) {
	side = strings.TrimSpace(side)
	if side == "" {
		return nil, fmt.Errorf("empty equation side")
	}

	// Split by + sign
	parts := strings.Split(side, "+")
	components := make([]ReactionComponent, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check for coefficient at the beginning
		coefficient := 1
		numEnd := 0
		for numEnd < len(part) && part[numEnd] >= '0' && part[numEnd] <= '9' {
			numEnd++
		}
		if numEnd > 0 {
			var err error
			coefficient, err = strconv.Atoi(part[:numEnd])
			if err != nil {
				return nil, fmt.Errorf("invalid coefficient: %w", err)
			}
			part = part[numEnd:]
		}

		substance, err := ParseFormula(part)
		if err != nil {
			return nil, err
		}

		components = append(components, ReactionComponent{
			Substance:   substance,
			Coefficient: coefficient,
		})
	}

	if len(components) == 0 {
		return nil, fmt.Errorf("no valid components found")
	}

	return components, nil
}

// Balance attempts to balance a chemical equation using matrix methods (SVD).
// Returns a BalancedReaction with coefficients and explanation steps.
func (r *Reaction) Balance() (*BalancedReaction, error) {
	// Collect all unique elements
	elementSet := make(map[string]bool)
	allComponents := append(r.Reactants, r.Products...)

	for _, comp := range allComponents {
		for elem := range comp.Substance.Composition {
			elementSet[elem] = true
		}
	}

	elements := make([]string, 0, len(elementSet))
	for elem := range elementSet {
		elements = append(elements, elem)
	}

	numElements := len(elements)
	numReactants := len(r.Reactants)
	numProducts := len(r.Products)
	numComponents := numReactants + numProducts

	if numComponents == 0 {
		return nil, fmt.Errorf("no components to balance")
	}

	// Build matrix: rows = elements, columns = components
	// Reactants have positive coefficients, products have negative
	matrixData := make([]float64, numElements*numComponents)

	for i, elem := range elements {
		for j, comp := range r.Reactants {
			count := comp.Substance.Composition[elem]
			matrixData[i*numComponents+j] = float64(count * comp.Coefficient)
		}
		for j, comp := range r.Products {
			count := comp.Substance.Composition[elem]
			// Products are negative in the matrix equation
			matrixData[i*numComponents+numReactants+j] = -float64(count * comp.Coefficient)
		}
	}

	matrix := mat.NewDense(numElements, numComponents, matrixData)

	// Solve using SVD to find null space
	coefficients, err := solveNullSpace(matrix, numComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to balance equation: %w", err)
	}

	// Convert to smallest whole numbers
	coefficients = normalizeCoefficients(coefficients)

	// Verify the balance
	isBalanced := verifyBalance(r, coefficients, elements)

	// Generate explanation steps
	steps := generateBalanceSteps(r, coefficients, elements)

	return &BalancedReaction{
		Reaction:     r,
		Coefficients: coefficients,
		Steps:        steps,
		IsBalanced:   isBalanced,
	}, nil
}

// solveNullSpace finds the null space of the matrix using SVD.
// For equation balancing, we need to find coefficient vector x such that A*x = 0.
func solveNullSpace(matrix *mat.Dense, numVars int) ([]int, error) {
	rows, cols := matrix.Dims()

	// For small matrices or when rows >= cols, use algebraic method
	if cols <= rows+1 {
		return balanceAlgebraicFromMatrix(matrix, numVars), nil
	}

	// Perform SVD decomposition
	svd := new(mat.SVD)
	ok := svd.Factorize(matrix, mat.SVDThin)
	if !ok {
		return balanceAlgebraicFromMatrix(matrix, numVars), nil
	}

	// Get singular values
	values := svd.Values(nil)

	// Number of non-zero singular values = matrix rank
	rank := 0
	for i := 0; i < len(values); i++ {
		if values[i] > 1e-10 {
			rank++
		}
	}

	nullity := cols - rank
	if nullity <= 0 {
		return balanceAlgebraicFromMatrix(matrix, numVars), nil
	}

	// For SVDThin: V has dimensions cols x min(rows,cols)
	vCols := min(rows, cols)
	V := mat.NewDense(cols, vCols, nil)
	svd.VTo(V)

	// Take the last column of V (corresponds to smallest singular value)
	nullVec := make([]float64, cols)
	for i := 0; i < cols; i++ {
		nullVec[i] = V.At(i, vCols-1)
	}

	if allZero(nullVec) {
		return balanceAlgebraicFromMatrix(matrix, numVars), nil
	}

	return floatToInt(normalizeVector(nullVec)), nil
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// balanceAlgebraicFromMatrix implements algebraic balancing.
func balanceAlgebraicFromMatrix(matrix *mat.Dense, numVars int) []int {
	rows, cols := matrix.Dims()

	// Determine number of reactants
	numReactants := 0
	for col := 0; col < cols; col++ {
		isReactant := false
		isProduct := false
		for row := 0; row < rows; row++ {
			val := matrix.At(row, col)
			if val > 0 {
				isReactant = true
			}
			if val < 0 {
				isProduct = true
			}
		}
		if isReactant && !isProduct {
			numReactants = col + 1
		}
	}

	if numReactants == 0 || numReactants >= cols {
		numReactants = cols / 2
	}

	coefficients := make([]int, cols)
	for i := range coefficients {
		coefficients[i] = 1
	}

	// Bidirectional iterative balancing
	for iteration := 0; iteration < 100; iteration++ {
		allBalanced := true

		for row := 0; row < rows; row++ {
			leftSum := 0.0
			rightSum := 0.0

			for col := 0; col < numReactants; col++ {
				val := matrix.At(row, col)
				if val > 0 {
					leftSum += val * float64(coefficients[col])
				}
			}

			for col := numReactants; col < cols; col++ {
				val := matrix.At(row, col)
				if val < 0 {
					rightSum += -val * float64(coefficients[col])
				}
			}

			if math.Abs(leftSum-rightSum) > 1e-6 {
				allBalanced = false

				if leftSum > 1e-6 && rightSum > 1e-6 {
					// If left > right, scale the right side
					if leftSum > rightSum {
						mult := math.Round(leftSum / rightSum)
						if mult < 1 {
							mult = 1
						}
						for col := numReactants; col < cols; col++ {
							val := matrix.At(row, col)
							if val != 0 {
								coefficients[col] = int(math.Round(float64(coefficients[col]) * mult))
								if coefficients[col] < 1 {
									coefficients[col] = 1
								}
								break
							}
						}
					} else {
						// If right > left, scale the left side
						mult := math.Round(rightSum / leftSum)
						if mult < 1 {
							mult = 1
						}
						for col := 0; col < numReactants; col++ {
							val := matrix.At(row, col)
							if val != 0 {
								coefficients[col] = int(math.Round(float64(coefficients[col]) * mult))
								if coefficients[col] < 1 {
									coefficients[col] = 1
								}
								break
							}
						}
					}
				} else if leftSum > 1e-6 && rightSum <= 1e-6 {
					// Element only on left - add on right
					for col := numReactants; col < cols; col++ {
						val := matrix.At(row, col)
						if val != 0 {
							coefficients[col] = int(math.Round(leftSum / -val))
							if coefficients[col] < 1 {
								coefficients[col] = 1
							}
							break
						}
					}
				} else if rightSum > 1e-6 && leftSum <= 1e-6 {
					// Element only on right - add on left
					for col := 0; col < numReactants; col++ {
						val := matrix.At(row, col)
						if val != 0 {
							coefficients[col] = int(math.Round(rightSum / val))
							if coefficients[col] < 1 {
								coefficients[col] = 1
							}
							break
						}
					}
				}
			}
		}

		if allBalanced {
			break
		}
	}

	// Normalize coefficients (divide by GCD)
	gcdAll := 0
	for _, c := range coefficients {
		if c > 0 {
			if gcdAll == 0 {
				gcdAll = c
			} else {
				gcdAll = gcd(gcdAll, c)
			}
		}
	}

	if gcdAll > 1 {
		for i := range coefficients {
			coefficients[i] /= gcdAll
		}
	}

	return coefficients
}

// lcm returns the least common multiple.
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// normalizeCoefficients converts coefficients to smallest whole numbers.
func normalizeCoefficients(coeffs []int) []int {
	if len(coeffs) == 0 {
		return coeffs
	}

	// Find GCD of all coefficients
	gcdVal := 0
	for _, c := range coeffs {
		if c <= 0 {
			c = -c
		}
		if c > 0 {
			if gcdVal == 0 {
				gcdVal = c
			} else {
				gcdVal = gcd(gcdVal, c)
			}
		}
	}

	if gcdVal == 0 {
		gcdVal = 1
	}

	// Divide by GCD
	result := make([]int, len(coeffs))
	for i, c := range coeffs {
		result[i] = c / gcdVal
		if result[i] <= 0 {
			result[i] = 1
		}
	}

	return result
}

// normalizeVector normalizes a float64 vector.
func normalizeVector(vec []float64) []float64 {
	if len(vec) == 0 {
		return vec
	}

	// Find minimum absolute non-zero value
	minVal := math.MaxFloat64
	for _, v := range vec {
		absV := math.Abs(v)
		if absV > 1e-10 && absV < minVal {
			minVal = absV
		}
	}

	if minVal == math.MaxFloat64 {
		return vec
	}

	// Normalize
	result := make([]float64, len(vec))
	for i, v := range vec {
		result[i] = v / minVal
	}

	return result
}

// floatToInt converts float64 slice to int slice.
func floatToInt(vec []float64) []int {
	result := make([]int, len(vec))
	for i, v := range vec {
		result[i] = int(math.Round(math.Abs(v)))
		if result[i] <= 0 {
			result[i] = 1
		}
	}
	return result
}

// allZero checks if all values in a slice are zero or near-zero.
func allZero(vec []float64) bool {
	for _, v := range vec {
		if math.Abs(v) > 1e-10 {
			return false
		}
	}
	return true
}

// gcd computes the greatest common divisor.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// verifyBalance checks if the balanced equation is actually balanced.
func verifyBalance(r *Reaction, coefficients []int, elements []string) bool {
	if len(coefficients) != len(r.Reactants)+len(r.Products) {
		return false
	}

	for _, elem := range elements {
		leftCount := 0
		rightCount := 0

		for i, comp := range r.Reactants {
			leftCount += comp.Substance.Composition[elem] * coefficients[i]
		}

		for i, comp := range r.Products {
			rightCount += comp.Substance.Composition[elem] * coefficients[len(r.Reactants)+i]
		}

		if leftCount != rightCount {
			return false
		}
	}

	return true
}

// generateBalanceSteps creates explanation steps for the balancing process.
func generateBalanceSteps(r *Reaction, coefficients []int, elements []string) []BalanceStep {
	steps := make([]BalanceStep, 0)

	// Step 1: Show unbalanced equation
	steps = append(steps, BalanceStep{
		Description: "Unbalanced equation",
		Equation:    formatEquation(r, make([]int, len(r.Reactants)+len(r.Products))),
	})

	// Step 2: Count atoms on each side
	stepDesc := "Atom count:"
	for _, elem := range elements {
		leftCount := 0
		rightCount := 0
		for _, comp := range r.Reactants {
			leftCount += comp.Substance.Composition[elem]
		}
		for _, comp := range r.Products {
			rightCount += comp.Substance.Composition[elem]
		}
		if leftCount != rightCount {
			stepDesc += fmt.Sprintf(" %s: left=%d, right=%d;", elem, leftCount, rightCount)
		}
	}
	steps = append(steps, BalanceStep{
		Description: stepDesc,
		Equation:    "",
	})

	// Step 3: Show balanced equation
	steps = append(steps, BalanceStep{
		Description: "Balanced equation",
		Equation:    formatEquation(r, coefficients),
	})

	return steps
}

// formatEquation formats a reaction with given coefficients.
func formatEquation(r *Reaction, coefficients []int) string {
	var sb strings.Builder

	for i, comp := range r.Reactants {
		if i > 0 {
			sb.WriteString(" + ")
		}
		coef := 1
		if i < len(coefficients) && coefficients[i] > 1 {
			coef = coefficients[i]
			sb.WriteString(fmt.Sprintf("%d", coef))
		}
		sb.WriteString(comp.Substance.Formula)
	}

	sb.WriteString(" → ")

	for i, comp := range r.Products {
		if i > 0 {
			sb.WriteString(" + ")
		}
		coef := 1
		idx := len(r.Reactants) + i
		if idx < len(coefficients) && coefficients[idx] > 1 {
			coef = coefficients[idx]
			sb.WriteString(fmt.Sprintf("%d", coef))
		}
		sb.WriteString(comp.Substance.Formula)
	}

	return sb.String()
}

// BalanceEquation is a convenience function to parse and balance an equation in one step.
func BalanceEquation(equation string) (*BalancedReaction, error) {
	reaction, err := ParseEquation(equation)
	if err != nil {
		return nil, err
	}
	return reaction.Balance()
}
