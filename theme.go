package main

import "github.com/charmbracelet/lipgloss"

type Theme struct {
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
}

type Styles struct {
	StatusBar      lipgloss.Style
	Modal          lipgloss.Style
	Prompt         lipgloss.Style
	TextInput      lipgloss.Style
	ListTitle      lipgloss.Style
	ListItem       lipgloss.Style
	ListItemActive lipgloss.Style
	ActivePane     lipgloss.Style
	InactivePane   lipgloss.Style
}

func NewStyles(theme Theme) Styles {
	s := Styles{}

	s.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text)).
		Background(lipgloss.Color(theme.Surface0))

	s.Modal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Mauve)).
		Padding(1, 2).
		Background(lipgloss.Color(theme.Base)).
		Foreground(lipgloss.Color(theme.Text))

	s.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Subtext0))

	s.TextInput = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Peach))

	s.ListTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Sapphire)).
		Padding(0, 1)

	s.ListItem = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Center)

	s.ListItemActive = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color(theme.Mauve)).
		Foreground(lipgloss.Color(theme.Base)).
		Align(lipgloss.Center)

	s.ActivePane = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Mauve))

	s.InactivePane = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Overlay0))

	return s
}
