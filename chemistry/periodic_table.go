// Package chemistry provides core types for representing chemical elements and compounds.
package chemistry

import (
	"strings"
)

// periodicTable contains all 118 known chemical elements.
// Data sourced from IUPAC standard atomic weights.
var periodicTable = map[string]Element{
	"H":  {Symbol: "H", Name: "Hydrogen", AtomicNumber: 1, AtomicMass: 1.008, OxidationStates: []int{-1, 1}, Group: 1},
	"He": {Symbol: "He", Name: "Helium", AtomicNumber: 2, AtomicMass: 4.0026, OxidationStates: []int{0}, Group: 18},
	"Li": {Symbol: "Li", Name: "Lithium", AtomicNumber: 3, AtomicMass: 6.94, OxidationStates: []int{1}, Group: 1},
	"Be": {Symbol: "Be", Name: "Beryllium", AtomicNumber: 4, AtomicMass: 9.0122, OxidationStates: []int{2}, Group: 2},
	"B":  {Symbol: "B", Name: "Boron", AtomicNumber: 5, AtomicMass: 10.81, OxidationStates: []int{3}, Group: 13},
	"C":  {Symbol: "C", Name: "Carbon", AtomicNumber: 6, AtomicMass: 12.011, OxidationStates: []int{-4, -3, -2, -1, 1, 2, 3, 4}, Group: 14},
	"N":  {Symbol: "N", Name: "Nitrogen", AtomicNumber: 7, AtomicMass: 14.007, OxidationStates: []int{-3, -2, -1, 1, 2, 3, 4, 5}, Group: 15},
	"O":  {Symbol: "O", Name: "Oxygen", AtomicNumber: 8, AtomicMass: 15.999, OxidationStates: []int{-2, -1, 1, 2}, Group: 16},
	"F":  {Symbol: "F", Name: "Fluorine", AtomicNumber: 9, AtomicMass: 18.998, OxidationStates: []int{-1}, Group: 17},
	"Ne": {Symbol: "Ne", Name: "Neon", AtomicNumber: 10, AtomicMass: 20.180, OxidationStates: []int{0}, Group: 18},
	"Na": {Symbol: "Na", Name: "Sodium", AtomicNumber: 11, AtomicMass: 22.990, OxidationStates: []int{1}, Group: 1},
	"Mg": {Symbol: "Mg", Name: "Magnesium", AtomicNumber: 12, AtomicMass: 24.305, OxidationStates: []int{2}, Group: 2},
	"Al": {Symbol: "Al", Name: "Aluminium", AtomicNumber: 13, AtomicMass: 26.982, OxidationStates: []int{3}, Group: 13},
	"Si": {Symbol: "Si", Name: "Silicon", AtomicNumber: 14, AtomicMass: 28.085, OxidationStates: []int{-4, -3, -2, -1, 1, 2, 3, 4}, Group: 14},
	"P":  {Symbol: "P", Name: "Phosphorus", AtomicNumber: 15, AtomicMass: 30.974, OxidationStates: []int{-3, -2, -1, 1, 2, 3, 4, 5}, Group: 15},
	"S":  {Symbol: "S", Name: "Sulfur", AtomicNumber: 16, AtomicMass: 32.06, OxidationStates: []int{-2, -1, 1, 2, 3, 4, 5, 6}, Group: 16},
	"Cl": {Symbol: "Cl", Name: "Chlorine", AtomicNumber: 17, AtomicMass: 35.45, OxidationStates: []int{-1, 1, 2, 3, 4, 5, 6, 7}, Group: 17},
	"Ar": {Symbol: "Ar", Name: "Argon", AtomicNumber: 18, AtomicMass: 39.948, OxidationStates: []int{0}, Group: 18},
	"K":  {Symbol: "K", Name: "Potassium", AtomicNumber: 19, AtomicMass: 39.098, OxidationStates: []int{1}, Group: 1},
	"Ca": {Symbol: "Ca", Name: "Calcium", AtomicNumber: 20, AtomicMass: 40.078, OxidationStates: []int{2}, Group: 2},
	"Sc": {Symbol: "Sc", Name: "Scandium", AtomicNumber: 21, AtomicMass: 44.956, OxidationStates: []int{3}, Group: 3},
	"Ti": {Symbol: "Ti", Name: "Titanium", AtomicNumber: 22, AtomicMass: 47.867, OxidationStates: []int{2, 3, 4}, Group: 4},
	"V":  {Symbol: "V", Name: "Vanadium", AtomicNumber: 23, AtomicMass: 50.942, OxidationStates: []int{2, 3, 4, 5}, Group: 5},
	"Cr": {Symbol: "Cr", Name: "Chromium", AtomicNumber: 24, AtomicMass: 51.996, OxidationStates: []int{2, 3, 6}, Group: 6},
	"Mn": {Symbol: "Mn", Name: "Manganese", AtomicNumber: 25, AtomicMass: 54.938, OxidationStates: []int{2, 3, 4, 6, 7}, Group: 7},
	"Fe": {Symbol: "Fe", Name: "Iron", AtomicNumber: 26, AtomicMass: 55.845, OxidationStates: []int{2, 3}, Group: 8},
	"Co": {Symbol: "Co", Name: "Cobalt", AtomicNumber: 27, AtomicMass: 58.933, OxidationStates: []int{2, 3}, Group: 9},
	"Ni": {Symbol: "Ni", Name: "Nickel", AtomicNumber: 28, AtomicMass: 58.693, OxidationStates: []int{2, 3}, Group: 10},
	"Cu": {Symbol: "Cu", Name: "Copper", AtomicNumber: 29, AtomicMass: 63.546, OxidationStates: []int{1, 2}, Group: 11},
	"Zn": {Symbol: "Zn", Name: "Zinc", AtomicNumber: 30, AtomicMass: 65.38, OxidationStates: []int{2}, Group: 12},
	"Ga": {Symbol: "Ga", Name: "Gallium", AtomicNumber: 31, AtomicMass: 69.723, OxidationStates: []int{3}, Group: 13},
	"Ge": {Symbol: "Ge", Name: "Germanium", AtomicNumber: 32, AtomicMass: 72.63, OxidationStates: []int{-4, 2, 4}, Group: 14},
	"As": {Symbol: "As", Name: "Arsenic", AtomicNumber: 33, AtomicMass: 74.922, OxidationStates: []int{-3, 2, 3, 5}, Group: 15},
	"Se": {Symbol: "Se", Name: "Selenium", AtomicNumber: 34, AtomicMass: 78.971, OxidationStates: []int{-2, 2, 4, 6}, Group: 16},
	"Br": {Symbol: "Br", Name: "Bromine", AtomicNumber: 35, AtomicMass: 79.904, OxidationStates: []int{-1, 1, 3, 4, 5}, Group: 17},
	"Kr": {Symbol: "Kr", Name: "Krypton", AtomicNumber: 36, AtomicMass: 83.798, OxidationStates: []int{0, 2}, Group: 18},
	"Rb": {Symbol: "Rb", Name: "Rubidium", AtomicNumber: 37, AtomicMass: 85.468, OxidationStates: []int{1}, Group: 1},
	"Sr": {Symbol: "Sr", Name: "Strontium", AtomicNumber: 38, AtomicMass: 87.62, OxidationStates: []int{2}, Group: 2},
	"Y":  {Symbol: "Y", Name: "Yttrium", AtomicNumber: 39, AtomicMass: 88.906, OxidationStates: []int{3}, Group: 3},
	"Zr": {Symbol: "Zr", Name: "Zirconium", AtomicNumber: 40, AtomicMass: 91.224, OxidationStates: []int{4}, Group: 4},
	"Nb": {Symbol: "Nb", Name: "Niobium", AtomicNumber: 41, AtomicMass: 92.906, OxidationStates: []int{3, 5}, Group: 5},
	"Mo": {Symbol: "Mo", Name: "Molybdenum", AtomicNumber: 42, AtomicMass: 95.95, OxidationStates: []int{3, 4, 6}, Group: 6},
	"Tc": {Symbol: "Tc", Name: "Technetium", AtomicNumber: 43, AtomicMass: 98, OxidationStates: []int{4, 6, 7}, Group: 7},
	"Ru": {Symbol: "Ru", Name: "Ruthenium", AtomicNumber: 44, AtomicMass: 101.07, OxidationStates: []int{3, 4}, Group: 8},
	"Rh": {Symbol: "Rh", Name: "Rhodium", AtomicNumber: 45, AtomicMass: 102.91, OxidationStates: []int{3}, Group: 9},
	"Pd": {Symbol: "Pd", Name: "Palladium", AtomicNumber: 46, AtomicMass: 106.42, OxidationStates: []int{2, 4}, Group: 10},
	"Ag": {Symbol: "Ag", Name: "Silver", AtomicNumber: 47, AtomicMass: 107.87, OxidationStates: []int{1}, Group: 11},
	"Cd": {Symbol: "Cd", Name: "Cadmium", AtomicNumber: 48, AtomicMass: 112.41, OxidationStates: []int{2}, Group: 12},
	"In": {Symbol: "In", Name: "Indium", AtomicNumber: 49, AtomicMass: 114.82, OxidationStates: []int{3}, Group: 13},
	"Sn": {Symbol: "Sn", Name: "Tin", AtomicNumber: 50, AtomicMass: 118.71, OxidationStates: []int{2, 4}, Group: 14},
	"Sb": {Symbol: "Sb", Name: "Antimony", AtomicNumber: 51, AtomicMass: 121.76, OxidationStates: []int{-3, 3, 5}, Group: 15},
	"Te": {Symbol: "Te", Name: "Tellurium", AtomicNumber: 52, AtomicMass: 127.60, OxidationStates: []int{-2, 2, 4, 6}, Group: 16},
	"I":  {Symbol: "I", Name: "Iodine", AtomicNumber: 53, AtomicMass: 126.90, OxidationStates: []int{-1, 1, 3, 5, 7}, Group: 17},
	"Xe": {Symbol: "Xe", Name: "Xenon", AtomicNumber: 54, AtomicMass: 131.29, OxidationStates: []int{2, 4, 6}, Group: 18},
	"Cs": {Symbol: "Cs", Name: "Caesium", AtomicNumber: 55, AtomicMass: 132.91, OxidationStates: []int{1}, Group: 1},
	"Ba": {Symbol: "Ba", Name: "Barium", AtomicNumber: 56, AtomicMass: 137.33, OxidationStates: []int{2}, Group: 2},
	"La": {Symbol: "La", Name: "Lanthanum", AtomicNumber: 57, AtomicMass: 138.91, OxidationStates: []int{3}, Group: 3},
	"Ce": {Symbol: "Ce", Name: "Cerium", AtomicNumber: 58, AtomicMass: 140.12, OxidationStates: []int{3, 4}, Group: 0},
	"Pr": {Symbol: "Pr", Name: "Praseodymium", AtomicNumber: 59, AtomicMass: 140.91, OxidationStates: []int{3}, Group: 0},
	"Nd": {Symbol: "Nd", Name: "Neodymium", AtomicNumber: 60, AtomicMass: 144.24, OxidationStates: []int{3}, Group: 0},
	"Pm": {Symbol: "Pm", Name: "Promethium", AtomicNumber: 61, AtomicMass: 145, OxidationStates: []int{3}, Group: 0},
	"Sm": {Symbol: "Sm", Name: "Samarium", AtomicNumber: 62, AtomicMass: 150.36, OxidationStates: []int{2, 3}, Group: 0},
	"Eu": {Symbol: "Eu", Name: "Europium", AtomicNumber: 63, AtomicMass: 151.96, OxidationStates: []int{2, 3}, Group: 0},
	"Gd": {Symbol: "Gd", Name: "Gadolinium", AtomicNumber: 64, AtomicMass: 157.25, OxidationStates: []int{3}, Group: 0},
	"Tb": {Symbol: "Tb", Name: "Terbium", AtomicNumber: 65, AtomicMass: 158.93, OxidationStates: []int{3}, Group: 0},
	"Dy": {Symbol: "Dy", Name: "Dysprosium", AtomicNumber: 66, AtomicMass: 162.50, OxidationStates: []int{3}, Group: 0},
	"Ho": {Symbol: "Ho", Name: "Holmium", AtomicNumber: 67, AtomicMass: 164.93, OxidationStates: []int{3}, Group: 0},
	"Er": {Symbol: "Er", Name: "Erbium", AtomicNumber: 68, AtomicMass: 167.26, OxidationStates: []int{3}, Group: 0},
	"Tm": {Symbol: "Tm", Name: "Thulium", AtomicNumber: 69, AtomicMass: 168.93, OxidationStates: []int{3}, Group: 0},
	"Yb": {Symbol: "Yb", Name: "Ytterbium", AtomicNumber: 70, AtomicMass: 173.05, OxidationStates: []int{2, 3}, Group: 0},
	"Lu": {Symbol: "Lu", Name: "Lutetium", AtomicNumber: 71, AtomicMass: 174.97, OxidationStates: []int{3}, Group: 3},
	"Hf": {Symbol: "Hf", Name: "Hafnium", AtomicNumber: 72, AtomicMass: 178.49, OxidationStates: []int{4}, Group: 4},
	"Ta": {Symbol: "Ta", Name: "Tantalum", AtomicNumber: 73, AtomicMass: 180.95, OxidationStates: []int{5}, Group: 5},
	"W":  {Symbol: "W", Name: "Tungsten", AtomicNumber: 74, AtomicMass: 183.84, OxidationStates: []int{4, 6}, Group: 6},
	"Re": {Symbol: "Re", Name: "Rhenium", AtomicNumber: 75, AtomicMass: 186.21, OxidationStates: []int{4, 6, 7}, Group: 7},
	"Os": {Symbol: "Os", Name: "Osmium", AtomicNumber: 76, AtomicMass: 190.23, OxidationStates: []int{3, 4, 6}, Group: 8},
	"Ir": {Symbol: "Ir", Name: "Iridium", AtomicNumber: 77, AtomicMass: 192.22, OxidationStates: []int{3, 4}, Group: 9},
	"Pt": {Symbol: "Pt", Name: "Platinum", AtomicNumber: 78, AtomicMass: 195.08, OxidationStates: []int{2, 4}, Group: 10},
	"Au": {Symbol: "Au", Name: "Gold", AtomicNumber: 79, AtomicMass: 196.97, OxidationStates: []int{1, 3}, Group: 11},
	"Hg": {Symbol: "Hg", Name: "Mercury", AtomicNumber: 80, AtomicMass: 200.59, OxidationStates: []int{1, 2}, Group: 12},
	"Tl": {Symbol: "Tl", Name: "Thallium", AtomicNumber: 81, AtomicMass: 204.38, OxidationStates: []int{1, 3}, Group: 13},
	"Pb": {Symbol: "Pb", Name: "Lead", AtomicNumber: 82, AtomicMass: 207.2, OxidationStates: []int{2, 4}, Group: 14},
	"Bi": {Symbol: "Bi", Name: "Bismuth", AtomicNumber: 83, AtomicMass: 208.98, OxidationStates: []int{3, 5}, Group: 15},
	"Po": {Symbol: "Po", Name: "Polonium", AtomicNumber: 84, AtomicMass: 209, OxidationStates: []int{2, 4}, Group: 16},
	"At": {Symbol: "At", Name: "Astatine", AtomicNumber: 85, AtomicMass: 210, OxidationStates: []int{-1, 1}, Group: 17},
	"Rn": {Symbol: "Rn", Name: "Radon", AtomicNumber: 86, AtomicMass: 222, OxidationStates: []int{2}, Group: 18},
	"Fr": {Symbol: "Fr", Name: "Francium", AtomicNumber: 87, AtomicMass: 223, OxidationStates: []int{1}, Group: 1},
	"Ra": {Symbol: "Ra", Name: "Radium", AtomicNumber: 88, AtomicMass: 226, OxidationStates: []int{2}, Group: 2},
	"Ac": {Symbol: "Ac", Name: "Actinium", AtomicNumber: 89, AtomicMass: 227, OxidationStates: []int{3}, Group: 3},
	"Th": {Symbol: "Th", Name: "Thorium", AtomicNumber: 90, AtomicMass: 232.04, OxidationStates: []int{4}, Group: 0},
	"Pa": {Symbol: "Pa", Name: "Protactinium", AtomicNumber: 91, AtomicMass: 231.04, OxidationStates: []int{5}, Group: 0},
	"U":  {Symbol: "U", Name: "Uranium", AtomicNumber: 92, AtomicMass: 238.03, OxidationStates: []int{3, 4, 6}, Group: 0},
	"Np": {Symbol: "Np", Name: "Neptunium", AtomicNumber: 93, AtomicMass: 237, OxidationStates: []int{4, 6}, Group: 0},
	"Pu": {Symbol: "Pu", Name: "Plutonium", AtomicNumber: 94, AtomicMass: 244, OxidationStates: []int{4, 6}, Group: 0},
	"Am": {Symbol: "Am", Name: "Americium", AtomicNumber: 95, AtomicMass: 243, OxidationStates: []int{3, 4}, Group: 0},
	"Cm": {Symbol: "Cm", Name: "Curium", AtomicNumber: 96, AtomicMass: 247, OxidationStates: []int{3}, Group: 0},
	"Bk": {Symbol: "Bk", Name: "Berkelium", AtomicNumber: 97, AtomicMass: 247, OxidationStates: []int{3, 4}, Group: 0},
	"Cf": {Symbol: "Cf", Name: "Californium", AtomicNumber: 98, AtomicMass: 251, OxidationStates: []int{3}, Group: 0},
	"Es": {Symbol: "Es", Name: "Einsteinium", AtomicNumber: 99, AtomicMass: 252, OxidationStates: []int{3}, Group: 0},
	"Fm": {Symbol: "Fm", Name: "Fermium", AtomicNumber: 100, AtomicMass: 257, OxidationStates: []int{3}, Group: 0},
	"Md": {Symbol: "Md", Name: "Mendelevium", AtomicNumber: 101, AtomicMass: 258, OxidationStates: []int{3}, Group: 0},
	"No": {Symbol: "No", Name: "Nobelium", AtomicNumber: 102, AtomicMass: 259, OxidationStates: []int{2, 3}, Group: 0},
	"Lr": {Symbol: "Lr", Name: "Lawrencium", AtomicNumber: 103, AtomicMass: 266, OxidationStates: []int{3}, Group: 3},
	"Rf": {Symbol: "Rf", Name: "Rutherfordium", AtomicNumber: 104, AtomicMass: 267, OxidationStates: []int{4}, Group: 4},
	"Db": {Symbol: "Db", Name: "Dubnium", AtomicNumber: 105, AtomicMass: 268, OxidationStates: []int{5}, Group: 5},
	"Sg": {Symbol: "Sg", Name: "Seaborgium", AtomicNumber: 106, AtomicMass: 269, OxidationStates: []int{6}, Group: 6},
	"Bh": {Symbol: "Bh", Name: "Bohrium", AtomicNumber: 107, AtomicMass: 270, OxidationStates: []int{7}, Group: 7},
	"Hs": {Symbol: "Hs", Name: "Hassium", AtomicNumber: 108, AtomicMass: 277, OxidationStates: []int{8}, Group: 8},
	"Mt": {Symbol: "Mt", Name: "Meitnerium", AtomicNumber: 109, AtomicMass: 278, OxidationStates: []int{}, Group: 9},
	"Ds": {Symbol: "Ds", Name: "Darmstadtium", AtomicNumber: 110, AtomicMass: 281, OxidationStates: []int{}, Group: 10},
	"Rg": {Symbol: "Rg", Name: "Roentgenium", AtomicNumber: 111, AtomicMass: 282, OxidationStates: []int{}, Group: 11},
	"Cn": {Symbol: "Cn", Name: "Copernicium", AtomicNumber: 112, AtomicMass: 285, OxidationStates: []int{}, Group: 12},
	"Nh": {Symbol: "Nh", Name: "Nihonium", AtomicNumber: 113, AtomicMass: 286, OxidationStates: []int{}, Group: 13},
	"Fl": {Symbol: "Fl", Name: "Flerovium", AtomicNumber: 114, AtomicMass: 289, OxidationStates: []int{}, Group: 14},
	"Mc": {Symbol: "Mc", Name: "Moscovium", AtomicNumber: 115, AtomicMass: 290, OxidationStates: []int{}, Group: 15},
	"Lv": {Symbol: "Lv", Name: "Livermorium", AtomicNumber: 116, AtomicMass: 293, OxidationStates: []int{}, Group: 16},
	"Ts": {Symbol: "Ts", Name: "Tennessine", AtomicNumber: 117, AtomicMass: 294, OxidationStates: []int{}, Group: 17},
	"Og": {Symbol: "Og", Name: "Oganesson", AtomicNumber: 118, AtomicMass: 294, OxidationStates: []int{}, Group: 18},
}

// GetElement returns an Element by its symbol (case-sensitive).
// Returns nil if the element is not found.
func GetElement(symbol string) *Element {
	if el, ok := periodicTable[symbol]; ok {
		return &el
	}
	return nil
}

// GetElementByName returns an Element by its name (case-insensitive).
// Returns nil if the element is not found.
func GetElementByName(name string) *Element {
	nameLower := strings.ToLower(name)
	for _, el := range periodicTable {
		if strings.ToLower(el.Name) == nameLower {
			return &el
		}
	}
	return nil
}

// GetAllElements returns a slice of all elements in the periodic table.
func GetAllElements() []Element {
	elements := make([]Element, 0, len(periodicTable))
	for _, el := range periodicTable {
		elements = append(elements, el)
	}
	return elements
}

// IsValidElement checks if a symbol corresponds to a known element.
func IsValidElement(symbol string) bool {
	_, ok := periodicTable[symbol]
	return ok
}
