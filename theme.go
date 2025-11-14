package main

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name                string
	FileName            string
	Accent              string
	Title               string
	Input               string
	Text                string
	SubtleText          string
	StatusBarBackground string
	ModalBackground     string
	Border              string
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
		Background(lipgloss.Color(theme.StatusBarBackground))

	s.Modal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Accent)).
		Padding(1, 2).
		Background(lipgloss.Color(theme.ModalBackground)).
		Foreground(lipgloss.Color(theme.Text))

	s.Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SubtleText))

	s.TextInput = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Input))

	s.ListTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Title)).
		Padding(0, 1)

	s.ListItem = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Center)

	s.ListItemActive = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color(theme.Accent)).
		Foreground(lipgloss.Color(theme.ModalBackground)).
		Align(lipgloss.Center)

	s.ActivePane = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Accent))

	s.InactivePane = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border))

	return s
}
