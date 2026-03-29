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

// Language represents the current UI language
type Language int

const (
	English Language = iota
	Russian
)

// Translations contains all UI strings for both languages
type Translations struct {
	Title               string
	SelectOption        string
	PressQToQuit        string
	ParseFormula        string
	LookUpElement       string
	CalculateMolarMass  string
	CompareFormulas     string
	BalanceEquation     string
	SuggestFormula      string
	Quit                string
	EnterFormula        string
	EnterElementSymbol  string
	EnterEquation       string
	EnterFirstFormula   string
	EnterSecondFormula  string
	FirstFormulaParsed  string
	PressEscToBack      string
	Formula             string
	Charge              string
	State               string
	Composition         string
	TotalAtoms          string
	Element             string
	Name                string
	AtomicNumber        string
	AtomicMass          string
	Group               string
	OxidationStates     string
	MolarMass           string
	Breakdown           string
	ComparisonResults   string
	Difference          string
	Balanced            string
	Coefficients        string
	Steps               string
	UnbalancedEquation  string
	AtomCount           string
	BalancedEquation    string
	Left                string
	Right               string
	ElementNotFound     string
	ThanksForUsing      string
	SelectLanguage      string
	EnglishLabel        string
	RussianLabel        string
	EnterForSuggestions string
	Suggestions         string
	Similarity          string
	Reason              string
	NoSuggestions       string
}

var translations = map[Language]Translations{
	English: {
		Title:               "🧬 CHEMATRIX Calculator 🧬",
		SelectOption:        "Select an option (use ↑↓ or j/k):",
		PressQToQuit:        "Press 'q' to quit",
		ParseFormula:        "🧪 Parse Chemical Formula",
		LookUpElement:       "⚛️  Look Up Element",
		CalculateMolarMass:  "📊 Calculate Molar Mass",
		CompareFormulas:     "⚖️  Compare Formulas",
		BalanceEquation:     "⚗️  Balance Equation",
		SuggestFormula:      "💡 Get Formula Suggestions",
		Quit:                "🚪 Quit",
		EnterFormula:        "Enter a chemical formula (e.g., H2SO4, Fe2(SO4)3, CuSO4·5H2O):",
		EnterElementSymbol:  "Enter an element symbol (e.g., Fe, O, Na):",
		EnterEquation:       "Enter a chemical equation (e.g., H2 + O2 -> H2O):",
		EnterFirstFormula:   "Enter the first formula:",
		EnterSecondFormula:  "Enter the second formula:",
		FirstFormulaParsed:  "✓ First formula parsed!",
		PressEscToBack:      "Press Esc to go back",
		Formula:             "Formula",
		Charge:              "Charge",
		State:               "State",
		Composition:         "Composition",
		TotalAtoms:          "Total atoms",
		Element:             "Element",
		Name:                "Name",
		AtomicNumber:        "Atomic Number",
		AtomicMass:          "Atomic Mass",
		Group:               "Group",
		OxidationStates:     "Oxidation States",
		MolarMass:           "Molar Mass",
		Breakdown:           "Breakdown",
		ComparisonResults:   "Comparison Results",
		Difference:          "Difference",
		Balanced:            "Balanced",
		Coefficients:        "Coefficients",
		Steps:               "Steps",
		UnbalancedEquation:  "Unbalanced equation",
		AtomCount:           "Atom count",
		BalancedEquation:    "Balanced equation",
		Left:                "left",
		Right:               "right",
		ElementNotFound:     "Element '%s' not found",
		ThanksForUsing:      "Thanks for using Chematrix! 🧬",
		SelectLanguage:      "Select language / Выберите язык:",
		EnglishLabel:        "English",
		RussianLabel:        "Русский",
		EnterForSuggestions: "Enter a formula for suggestions (e.g., HO, H20, CL):",
		Suggestions:         "Suggestions",
		Similarity:          "Similarity",
		Reason:              "Reason",
		NoSuggestions:       "No suggestions available",
	},
	Russian: {
		Title:               "🧬 CHEMATRIX Калькулятор 🧬",
		SelectOption:        "Выберите опцию (используйте ↑↓ или j/k):",
		PressQToQuit:        "Нажмите 'q' для выхода",
		ParseFormula:        "🧪 Разобрать формулу",
		LookUpElement:       "⚛️  Найти элемент",
		CalculateMolarMass:  "📊 Молярная масса",
		CompareFormulas:     "⚖️  Сравнить формулы",
		BalanceEquation:     "⚗️  Балансировка уравнения",
		SuggestFormula:      "💡 Подсказки формул",
		Quit:                "🚪 Выход",
		EnterFormula:        "Введите химическую формулу (например, H2SO4, Fe2(SO4)3, CuSO4·5H2O):",
		EnterElementSymbol:  "Введите символ элемента (например, Fe, O, Na):",
		EnterEquation:       "Введите химическое уравнение (например, H2 + O2 -> H2O):",
		EnterFirstFormula:   "Введите первую формулу:",
		EnterSecondFormula:  "Введите вторую формулу:",
		FirstFormulaParsed:  "✓ Первая формула разобрана!",
		PressEscToBack:      "Нажмите Esc для возврата",
		Formula:             "Формула",
		Charge:              "Заряд",
		State:               "Состояние",
		Composition:         "Состав",
		TotalAtoms:          "Всего атомов",
		Element:             "Элемент",
		Name:                "Название",
		AtomicNumber:        "Атомный номер",
		AtomicMass:          "Атомная масса",
		Group:               "Группа",
		OxidationStates:     "Степени окисления",
		MolarMass:           "Молярная масса",
		Breakdown:           "Разбор",
		ComparisonResults:   "Результаты сравнения",
		Difference:          "Разница",
		Balanced:            "Сбалансировано",
		Coefficients:        "Коэффициенты",
		Steps:               "Шаги",
		UnbalancedEquation:  "Исходное уравнение",
		AtomCount:           "Подсчёт атомов",
		BalancedEquation:    "Сбалансированное уравнение",
		Left:                "слева",
		Right:               "справа",
		ElementNotFound:     "Элемент '%s' не найден",
		ThanksForUsing:      "Спасибо за использование Chematrix! 🧬",
		SelectLanguage:      "Select language / Выберите язык:",
		EnglishLabel:        "English",
		RussianLabel:        "Русский",
		EnterForSuggestions: "Введите формулу для подсказок (например, HO, H20, CL):",
		Suggestions:         "Подсказки",
		Similarity:          "Схожесть",
		Reason:              "Причина",
		NoSuggestions:       "Нет подсказок",
	},
}

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
	languageScreen screen = iota
	menuScreen
	parseScreen
	elementScreen
	calculateScreen
	compareScreen
	balanceScreen
	suggestScreen
	quitScreen
)

// Main model
type model struct {
	screen        screen
	menuIndex     int
	input         string
	result        string
	errorMsg      string
	substance     *chemistry.Substance
	elements      []chemistry.Element
	comp1         chemistry.Composition
	comp2         chemistry.Composition
	compareStep   int
	balanceResult *chemistry.BalancedReaction
	suggestions   []chemistry.FormulaSuggestion
	lang          Language
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
		screen:    languageScreen,
		menuIndex: 0,
		elements:  chemistry.GetAllElements(),
		lang:      English,
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
	case languageScreen:
		return m.handleLanguageInput(msg)
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
	case balanceScreen:
		return m.handleBalanceInput(msg)
	case suggestScreen:
		return m.handleSuggestInput(msg)
	}

	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleLanguageInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < 1 {
			m.menuIndex++
		}
	case "enter":
		if m.menuIndex == 0 {
			m.lang = English
		} else {
			m.lang = Russian
		}
		m.screen = menuScreen
		m.menuIndex = 0
	case "q":
		m.screen = quitScreen
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
		if m.menuIndex < 6 {
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
			m.screen = balanceScreen
			m.input = ""
			m.result = ""
			m.errorMsg = ""
			m.balanceResult = nil
		case 5:
			m.screen = suggestScreen
			m.input = ""
			m.result = ""
			m.errorMsg = ""
			m.suggestions = nil
		case 6:
			m.screen = quitScreen
			return m, tea.Quit
		}
	case "l":
		// Switch language
		m.screen = languageScreen
		m.menuIndex = int(m.lang)
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
				m.result = formatSubstanceResult(sub, m.lang)
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
				t := translations[m.lang]
				m.errorMsg = fmt.Sprintf(t.ElementNotFound, m.input)
				m.result = ""
			} else {
				m.errorMsg = ""
				m.result = formatElementResult(el, m.lang)
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
				m.result = formatCalculationResult(sub, m.lang)
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
			m.result = formatCompareResult(m.comp1, m.comp2, m.lang)
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

func (m model) handleBalanceInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.input != "" {
			result, err := chemistry.BalanceEquation(m.input)
			if err != nil {
				m.errorMsg = err.Error()
				m.result = ""
			} else {
				m.errorMsg = ""
				m.balanceResult = result
				m.result = formatBalanceResult(result, m.lang)
			}
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.input = ""
		m.result = ""
		m.errorMsg = ""
		m.balanceResult = nil
	default:
		m.input = updateInput(m.input, msg)
	}
	return m, nil
}

func (m model) handleSuggestInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.input != "" {
			suggestions := chemistry.SuggestFormulas(m.input)
			m.suggestions = suggestions
			m.errorMsg = ""
			m.result = formatSuggestionsResult(suggestions, m.lang)
		}
	case tea.KeyEsc:
		m.screen = menuScreen
		m.input = ""
		m.result = ""
		m.errorMsg = ""
		m.suggestions = nil
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
	case languageScreen:
		s = m.viewLanguage()
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
	case balanceScreen:
		s = m.viewBalance()
	case suggestScreen:
		s = m.viewSuggest()
	case quitScreen:
		s = viewQuit(m.lang)
	}

	return s
}

func (m model) viewLanguage() string {
	t := translations[m.lang]

	languages := []string{
		t.EnglishLabel,
		t.RussianLabel,
	}

	s := titleStyle.Render("╔════════════════════════════════════╗")
	s += "\n"
	s += titleStyle.Render("║   🌐 Language Selection / Язык     ║")
	s += "\n"
	s += titleStyle.Render("╚════════════════════════════════════╝")
	s += "\n\n"
	s += infoStyle.Render(t.SelectLanguage)
	s += "\n\n"

	for i, lang := range languages {
		if i == m.menuIndex {
			s += selectedStyle.Render("❯ " + lang)
		} else {
			s += menuStyle.Render("  " + lang)
		}
		s += "\n"
	}

	s += "\n" + menuStyle.Render(t.PressQToQuit)
	return s
}

func (m model) viewMenu() string {
	t := translations[m.lang]

	menuItems := []string{
		t.ParseFormula,
		t.LookUpElement,
		t.CalculateMolarMass,
		t.CompareFormulas,
		t.BalanceEquation,
		t.SuggestFormula,
		t.Quit,
	}

	s := titleStyle.Render("╔════════════════════════════════════╗")
	s += "\n"
	s += titleStyle.Render("║   " + t.Title + "       ║")
	s += "\n"
	s += titleStyle.Render("╚════════════════════════════════════╝")
	s += "\n\n"
	s += infoStyle.Render(t.SelectOption)
	s += "\n\n"

	for i, item := range menuItems {
		if i == m.menuIndex {
			s += selectedStyle.Render("❯ " + item)
		} else {
			s += menuStyle.Render("  " + item)
		}
		s += "\n"
	}

	s += "\n" + menuStyle.Render(t.PressQToQuit)
	s += "\n" + menuStyle.Render("Press 'l' to change language / Сменить язык")
	return s
}

func (m model) viewParse() string {
	t := translations[m.lang]

	s := titleStyle.Render("🧪 " + t.ParseFormula)
	s += "\n\n"
	s += infoStyle.Render(t.EnterFormula)
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

	s += "\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func (m model) viewElement() string {
	t := translations[m.lang]

	s := titleStyle.Render("⚛️  " + t.LookUpElement)
	s += "\n\n"
	s += infoStyle.Render(t.EnterElementSymbol)
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

	s += "\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func (m model) viewCalculate() string {
	t := translations[m.lang]

	s := titleStyle.Render("📊 " + t.CalculateMolarMass)
	s += "\n\n"
	s += infoStyle.Render(t.EnterFormula)
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

	s += "\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func (m model) viewCompare() string {
	t := translations[m.lang]

	s := titleStyle.Render("⚖️  " + t.CompareFormulas)
	s += "\n\n"

	if m.compareStep == 0 {
		s += infoStyle.Render(t.EnterFirstFormula)
		s += "\n\n"
		s += inputStyle.Render("> " + m.input + "_")
	} else if m.compareStep == 1 {
		s += successStyle.Render(t.FirstFormulaParsed)
		s += "\n\n"
		s += infoStyle.Render(t.EnterSecondFormula)
		s += "\n\n"
		s += inputStyle.Render("> " + m.input + "_")
	}

	if m.errorMsg != "" {
		s += "\n" + errorStyle.Render("❌ "+m.errorMsg)
	}

	if m.result != "" {
		s += "\n" + resultStyle.Render(m.result)
	}

	s += "\n\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func (m model) viewBalance() string {
	t := translations[m.lang]

	s := titleStyle.Render("⚗️  " + t.BalanceEquation)
	s += "\n\n"
	s += infoStyle.Render(t.EnterEquation)
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

	s += "\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func (m model) viewSuggest() string {
	t := translations[m.lang]

	s := titleStyle.Render("💡 " + t.SuggestFormula)
	s += "\n\n"
	s += infoStyle.Render(t.EnterForSuggestions)
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

	s += "\n" + menuStyle.Render(t.PressEscToBack)
	return s
}

func viewQuit(lang Language) string {
	t := translations[lang]
	return "\n" + titleStyle.Render(t.ThanksForUsing) + "\n\n"
}

// Formatting helpers
func formatSubstanceResult(sub *chemistry.Substance, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s: %s\n", t.Formula, successStyle.Render(sub.Formula)))
	sb.WriteString(fmt.Sprintf("%s: %d\n", t.Charge, sub.Charge))
	if sub.State != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\n", t.State, sub.State))
	}
	sb.WriteString("\n" + t.Composition + ":\n")

	for elem, count := range sub.Composition {
		el := chemistry.GetElement(elem)
		name := elem
		if el != nil {
			if lang == Russian {
				name = fmt.Sprintf("%s (%s)", elem, el.Name)
			} else {
				name = fmt.Sprintf("%s (%s)", elem, el.Name)
			}
		}
		sb.WriteString(fmt.Sprintf("  • %s: %d atom(s)\n", name, count))
	}

	sb.WriteString(fmt.Sprintf("\n%s: %d", t.TotalAtoms, sub.Composition.TotalAtoms()))

	return sb.String()
}

func formatElementResult(el *chemistry.Element, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s: %s\n", t.Element, successStyle.Render(el.Symbol)))
	sb.WriteString(fmt.Sprintf("%s: %s\n", t.Name, el.Name))
	sb.WriteString(fmt.Sprintf("%s: %d\n", t.AtomicNumber, el.AtomicNumber))
	sb.WriteString(fmt.Sprintf("%s: %.3f u\n", t.AtomicMass, el.AtomicMass))
	sb.WriteString(fmt.Sprintf("%s: %d\n", t.Group, el.Group))

	if len(el.OxidationStates) > 0 {
		states := make([]string, len(el.OxidationStates))
		for i, s := range el.OxidationStates {
			if s > 0 {
				states[i] = fmt.Sprintf("+%d", s)
			} else {
				states[i] = fmt.Sprintf("%d", s)
			}
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", t.OxidationStates, strings.Join(states, ", ")))
	}

	return sb.String()
}

func formatCalculationResult(sub *chemistry.Substance, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s: %s\n\n", t.Formula, successStyle.Render(sub.Formula)))

	molarMass := sub.MolarMass()
	sb.WriteString(fmt.Sprintf("%s: %s\n", t.MolarMass, successStyle.Render(fmt.Sprintf("%.3f g/mol", molarMass))))
	sb.WriteString(fmt.Sprintf("%s: %d\n", t.TotalAtoms, sub.Composition.TotalAtoms()))

	sb.WriteString("\n" + t.Breakdown + ":\n")
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

func formatCompareResult(comp1, comp2 chemistry.Composition, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	sb.WriteString(t.ComparisonResults + ":\n\n")

	allElements := make(map[string]bool)
	for elem := range comp1 {
		allElements[elem] = true
	}
	for elem := range comp2 {
		allElements[elem] = true
	}

	sb.WriteString("Element | Formula 1 | Formula 2 | " + t.Difference + "\n")
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

func formatBalanceResult(result *chemistry.BalancedReaction, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s: %v\n", t.Balanced, successStyle.Render(fmt.Sprintf("%v", result.IsBalanced))))
	sb.WriteString(fmt.Sprintf("%s: %v\n\n", t.Coefficients, result.Coefficients))

	sb.WriteString(t.Steps + ":\n")
	for i, step := range result.Steps {
		if step.Equation != "" {
			desc := step.Description
			// Translate common descriptions
			switch step.Description {
			case "Unbalanced equation":
				desc = t.UnbalancedEquation
			case "Atom count:":
				desc = t.AtomCount
			case "Balanced equation":
				desc = t.BalancedEquation
			}
			// Translate left/right in atom count
			equation := step.Equation
			if lang == Russian {
				equation = strings.ReplaceAll(equation, "left=", t.Left+"=")
				equation = strings.ReplaceAll(equation, "right=", t.Right+"=")
			}
			sb.WriteString(fmt.Sprintf("  %d. %s: %s\n", i+1, desc, equation))
		} else {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step.Description))
		}
	}

	return sb.String()
}

func formatSuggestionsResult(suggestions []chemistry.FormulaSuggestion, lang Language) string {
	t := translations[lang]
	var sb strings.Builder

	if len(suggestions) == 0 {
		return infoStyle.Render(t.NoSuggestions)
	}

	sb.WriteString(fmt.Sprintf("%s:\n\n", t.Suggestions))

	for i, s := range suggestions {
		sb.WriteString(fmt.Sprintf("%d. ", i+1))
		sb.WriteString(successStyle.Render(s.Formula))
		sb.WriteString(fmt.Sprintf(" (%s: %.0f%%)\n", t.Similarity, s.Similarity*100))

		if s.Reason != "" {
			sb.WriteString(fmt.Sprintf("   %s: %s\n", t.Reason, infoStyle.Render(s.Reason)))
		}

		if s.Substance != nil {
			molarMass := s.Substance.MolarMass()
			sb.WriteString(fmt.Sprintf("   %s: %.2f g/mol\n", t.MolarMass, molarMass))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
