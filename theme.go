package main

import "github.com/charmbracelet/lipgloss"

var CatppucinMocha = struct {
	Rosewater string
	Flamingo  string
	Pink      string
	Mauve     string
	Red       string
	Maroon    string
	Peach     string
	Yellow    string
	Green     string
	Teal      string
	Sky       string
	Sapphire  string
	Blue      string
	Lavender  string
	Text      string
	Subtext1  string
	Subtext0  string
	Overlay2  string
	Overlay1  string
	Overlay0  string
	Surface2  string
	Surface1  string
	Surface0  string
	Base      string
	Mantle    string
	Crust     string
}{
	Rosewater: "#F5E0DC",
	Flamingo:  "#F2CDCD",
	Pink:      "#F5C2E7",
	Mauve:     "#CBA6F7",
	Red:       "#F38BA8",
	Maroon:    "#EBA0AC",
	Peach:     "#FAB387",
	Yellow:    "#F9E2AF",
	Green:     "#A6E3A1",
	Teal:      "#94E2D5",
	Sky:       "#89DCEB",
	Sapphire:  "#74C7EC",
	Blue:      "#89B4FA",
	Lavender:  "#B4BEFE",
	Text:      "#CDD6F4",
	Subtext1:  "#BAC2DE",
	Subtext0:  "#A6ADC8",
	Overlay2:  "#9399B2",
	Overlay1:  "#7F849C",
	Overlay0:  "#6C7086",
	Surface2:  "#585B70",
	Surface1:  "#454574A",
	Surface0:  "#313244",
	Base:      "#1E1E2E",
	Mantle:    "#181825",
	Crust:     "#11111B",
}

type Styles struct {
	StatusBar      lipgloss.Style
	Modal          lipgloss.Style
	Prompt         lipgloss.Style
	TextInput      lipgloss.Style
	ListTitle      lipgloss.Style
	ListItem       lipgloss.Style
	ListItemActive lipgloss.Style
	Overlay        lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(CatppucinMocha.Text)).
			Background(lipgloss.Color(CatppucinMocha.Surface0)),
		Modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Background(lipgloss.Color(CatppucinMocha.Mauve)).
			Padding(1, 2).
			Background(lipgloss.Color(CatppucinMocha.Base)).
			Foreground(lipgloss.Color(CatppucinMocha.Text)),
		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color(CatppucinMocha.Subtext0)),
		TextInput: lipgloss.NewStyle().
			Foreground(lipgloss.Color(CatppucinMocha.Peach)),
		ListTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(CatppucinMocha.Sapphire)).
			Padding(0, 1),
		ListItem: lipgloss.NewStyle().
			Padding(0, 1),
		ListItemActive: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color(CatppucinMocha.Mauve)).
			Foreground(lipgloss.Color(CatppucinMocha.Base)),
		Overlay: lipgloss.NewStyle().
			Background(lipgloss.Color(CatppucinMocha.Surface0)),
	}
}
