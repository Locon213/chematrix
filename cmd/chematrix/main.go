// Chematrix CLI - Interactive Chemistry Calculator
// Run with: go run cmd/chematrix/main.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Locon213/chematrix/chemistry"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6")).
			MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			MarginLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50FA7B")).
			MarginLeft(2)

	resultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8BE9FD")).
			Padding(1, 2).
			Margin(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Bold(true)
)

// Screen types
type screen int

const (
	menuScreen screen = iota
	parseScreen
	elementScreen
	calculateScreen
	compareScreen
	quitScreen
)

// Main model
type model struct {
	screen      screen
	menuIndex   int
	input       string
	result      string
	errorMsg    string
	substance   *chemistry.Substance
	elements    []chemistry.Element
	comp1       chemistry.Composition
	comp2       chemistry.Composition
	compareStep int
}

// Messages
type inputMsg string
type resultMsg string
type errorMsg string

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		screen:   menuScreen,
		elements: chemistry.GetAllElements(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case inputMsg:
		m.input = string(msg)
		return m, nil
	case resultMsg:
		m.result = string(msg)
		return m, nil
	case errorMsg:
		m.errorMsg = string(msg)
		return m, nil
	}
	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case menuScreen:
		return m.handleMenuInput(msg)
	case parseScreen:
		return m.handleParseInput(msg)
	case elementScreen:
		return m.handleElementInput(msg)
	case calculateScreen:
		return m.handleCalculateInput(msg)
	case compareScreen:
		return m.handleCompareInput(msg)
	}

	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < 5 {
			m.menuIndex++
		}
	case "enter":
		switch m.menuIndex {
		case 0:
			m.screen = parseScreen
			m.input = ""
			m.result = ""
			m.errorMsg = ""
		case 1:
			m.screen = elementScreen
			m.input = ""
			m.result = ""
		case 2:
			m.screen = calculateScreen
			m.input = ""
			m.result = ""
		case 3:
			m.screen = compareScreen
			m.comp1 = nil
			m.comp2 = nil
			m.compareStep = 0
			m.result = ""
		case 4:
			m.screen = quitScreen
			return m, tea.Quit
		}
	case "q":
		m.screen = quitScreen
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleParseInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.input != "" {
			sub, err := chemistry.ParseFormula(m.input)
			if err != nil {
				m.errorMsg = err.Error()
				m.result = ""
			} else {
				m.substance = sub
				m.errorMsg = ""
				m.result = formatSubstanceResult(sub)
			}
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.input = ""
		m.result = ""
		m.errorMsg = ""
	default:
		m.input = updateInput(m.input, msg)
	}
	return m, nil
}

func (m model) handleElementInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.input != "" {
			el := chemistry.GetElement(m.input)
			if el == nil {
				m.errorMsg = fmt.Sprintf("Element '%s' not found", m.input)
				m.result = ""
			} else {
				m.errorMsg = ""
				m.result = formatElementResult(el)
			}
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.input = ""
		m.result = ""
		m.errorMsg = ""
	default:
		m.input = updateInput(m.input, msg)
	}
	return m, nil
}

func (m model) handleCalculateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.input != "" {
			sub, err := chemistry.ParseFormula(m.input)
			if err != nil {
				m.errorMsg = err.Error()
				m.result = ""
			} else {
				m.errorMsg = ""
				m.result = formatCalculationResult(sub)
			}
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.input = ""
		m.result = ""
		m.errorMsg = ""
	default:
		m.input = updateInput(m.input, msg)
	}
	return m, nil
}

func (m model) handleCompareInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.compareStep == 0 {
			sub1, err := chemistry.ParseFormula(m.input)
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.comp1 = sub1.Composition
			m.compareStep = 1
			m.input = ""
			m.errorMsg = ""
		} else {
			sub2, err := chemistry.ParseFormula(m.input)
			if err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.comp2 = sub2.Composition
			m.errorMsg = ""
			m.result = formatCompareResult(m.comp1, m.comp2)
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.comp1 = nil
		m.comp2 = nil
		m.compareStep = 0
		m.input = ""
		m.result = ""
		m.errorMsg = ""
	default:
		m.input = updateInput(m.input, msg)
	}
	return m, nil
}

func updateInput(input string, msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyBackspace:
		if len(input) > 0 {
			return input[:len(input)-1]
		}
	case tea.KeyCtrlU:
		return ""
	default:
		if len(msg.String()) == 1 {
			return input + msg.String()
		}
	}
	return input
}

func (m model) View() string {
	var s string

	switch m.screen {
	case menuScreen:
		s = m.viewMenu()
	case parseScreen:
		s = m.viewParse()
	case elementScreen:
		s = m.viewElement()
	case calculateScreen:
		s = m.viewCalculate()
	case compareScreen:
		s = m.viewCompare()
	case quitScreen:
		s = viewQuit()
	}

	return s
}

func (m model) viewMenu() string {
	menuItems := []string{
		"🧪 Parse Chemical Formula",
		"⚛️  Look Up Element",
		"📊 Calculate Molar Mass",
		"⚖️  Compare Formulas",
		"🚪 Quit",
	}

	s := titleStyle.Render("╔════════════════════════════════════╗")
	s += "\n"
	s += titleStyle.Render("║   🧬 CHEMATRIX Calculator 🧬       ║")
	s += "\n"
	s += titleStyle.Render("╚════════════════════════════════════╝")
	s += "\n\n"
	s += infoStyle.Render("Select an option (use ↑↓ or j/k):")
	s += "\n\n"

	for i, item := range menuItems {
		if i == m.menuIndex {
			s += selectedStyle.Render("❯ " + item)
		} else {
			s += menuStyle.Render("  " + item)
		}
		s += "\n"
	}

	s += "\n" + menuStyle.Render("Press 'q' to quit")
	return s
}

func (m model) viewParse() string {
	s := titleStyle.Render("🧪 Parse Chemical Formula")
	s += "\n\n"
	s += infoStyle.Render("Enter a chemical formula (e.g., H2SO4, Fe2(SO4)3, CuSO4·5H2O):")
	s += "\n\n"
	s += inputStyle.Render("> " + m.input + "_")
	s += "\n\n"

	if m.errorMsg != "" {
		s += errorStyle.Render("❌ " + m.errorMsg)
		s += "\n"
	}

	if m.result != "" {
		s += resultStyle.Render(m.result)
	}

	s += "\n" + menuStyle.Render("Press Esc to go back")
	return s
}

func (m model) viewElement() string {
	s := titleStyle.Render("⚛️  Look Up Element")
	s += "\n\n"
	s += infoStyle.Render("Enter an element symbol (e.g., Fe, O, Na):")
	s += "\n\n"
	s += inputStyle.Render("> " + m.input + "_")
	s += "\n\n"

	if m.errorMsg != "" {
		s += errorStyle.Render("❌ " + m.errorMsg)
		s += "\n"
	}

	if m.result != "" {
		s += resultStyle.Render(m.result)
	}

	s += "\n" + menuStyle.Render("Press Esc to go back")
	return s
}

func (m model) viewCalculate() string {
	s := titleStyle.Render("📊 Calculate Molar Mass")
	s += "\n\n"
	s += infoStyle.Render("Enter a chemical formula:")
	s += "\n\n"
	s += inputStyle.Render("> " + m.input + "_")
	s += "\n\n"

	if m.errorMsg != "" {
		s += errorStyle.Render("❌ " + m.errorMsg)
		s += "\n"
	}

	if m.result != "" {
		s += resultStyle.Render(m.result)
	}

	s += "\n" + menuStyle.Render("Press Esc to go back")
	return s
}

func (m model) viewCompare() string {
	s := titleStyle.Render("⚖️  Compare Formulas")
	s += "\n\n"

	if m.compareStep == 0 {
		s += infoStyle.Render("Enter the first formula:")
		s += "\n\n"
		s += inputStyle.Render("> " + m.input + "_")
	} else if m.compareStep == 1 {
		s += successStyle.Render("✓ First formula parsed!")
		s += "\n\n"
		s += infoStyle.Render("Enter the second formula:")
		s += "\n\n"
		s += inputStyle.Render("> " + m.input + "_")
	}

	if m.errorMsg != "" {
		s += "\n" + errorStyle.Render("❌ "+m.errorMsg)
	}

	if m.result != "" {
		s += "\n" + resultStyle.Render(m.result)
	}

	s += "\n\n" + menuStyle.Render("Press Esc to go back")
	return s
}

func viewQuit() string {
	return "\n" + titleStyle.Render("Thanks for using Chematrix! 🧬") + "\n\n"
}

// Formatting helpers
func formatSubstanceResult(sub *chemistry.Substance) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Formula: %s\n", successStyle.Render(sub.Formula)))
	sb.WriteString(fmt.Sprintf("Charge: %d\n", sub.Charge))
	if sub.State != "" {
		sb.WriteString(fmt.Sprintf("State: %s\n", sub.State))
	}
	sb.WriteString("\nComposition:\n")

	for elem, count := range sub.Composition {
		el := chemistry.GetElement(elem)
		name := elem
		if el != nil {
			name = fmt.Sprintf("%s (%s)", elem, el.Name)
		}
		sb.WriteString(fmt.Sprintf("  • %s: %d atom(s)\n", name, count))
	}

	sb.WriteString(fmt.Sprintf("\nTotal atoms: %d", sub.Composition.TotalAtoms()))

	return sb.String()
}

func formatElementResult(el *chemistry.Element) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Element: %s\n", successStyle.Render(el.Symbol)))
	sb.WriteString(fmt.Sprintf("Name: %s\n", el.Name))
	sb.WriteString(fmt.Sprintf("Atomic Number: %d\n", el.AtomicNumber))
	sb.WriteString(fmt.Sprintf("Atomic Mass: %.3f u\n", el.AtomicMass))
	sb.WriteString(fmt.Sprintf("Group: %d\n", el.Group))

	if len(el.OxidationStates) > 0 {
		states := make([]string, len(el.OxidationStates))
		for i, s := range el.OxidationStates {
			if s > 0 {
				states[i] = fmt.Sprintf("+%d", s)
			} else {
				states[i] = fmt.Sprintf("%d", s)
			}
		}
		sb.WriteString(fmt.Sprintf("Oxidation States: %s\n", strings.Join(states, ", ")))
	}

	return sb.String()
}

func formatCalculationResult(sub *chemistry.Substance) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Formula: %s\n\n", successStyle.Render(sub.Formula)))

	molarMass := sub.MolarMass()
	sb.WriteString(fmt.Sprintf("Molar Mass: %s\n", successStyle.Render(fmt.Sprintf("%.3f g/mol", molarMass))))
	sb.WriteString(fmt.Sprintf("Total atoms: %d\n", sub.Composition.TotalAtoms()))

	sb.WriteString("\nBreakdown:\n")
	for elem, count := range sub.Composition {
		el := chemistry.GetElement(elem)
		if el != nil {
			mass := el.AtomicMass * float64(count)
			percent := (mass / molarMass) * 100
			sb.WriteString(fmt.Sprintf("  • %s: %.3f g/mol (%.1f%%)\n", elem, mass, percent))
		}
	}

	return sb.String()
}

func formatCompareResult(comp1, comp2 chemistry.Composition) string {
	var sb strings.Builder

	sb.WriteString("Comparison Results:\n\n")

	// Get all unique elements
	allElements := make(map[string]bool)
	for elem := range comp1 {
		allElements[elem] = true
	}
	for elem := range comp2 {
		allElements[elem] = true
	}

	sb.WriteString("Element | Formula 1 | Formula 2 | Difference\n")
	sb.WriteString("--------|-----------|-----------|------------\n")

	for elem := range allElements {
		c1 := comp1[elem]
		c2 := comp2[elem]
		diff := c2 - c1
		diffStr := strconv.Itoa(diff)
		if diff > 0 {
			diffStr = "+" + diffStr
		}
		sb.WriteString(fmt.Sprintf("  %s     |    %d      |    %d      |   %s\n", elem, c1, c2, diffStr))
	}

	total1 := comp1.TotalAtoms()
	total2 := comp2.TotalAtoms()

	sb.WriteString(fmt.Sprintf("\nTotal atoms: %d → %d (%+d)\n", total1, total2, total2-total1))

	return sb.String()
}
