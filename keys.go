package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	New    key.Binding
	Delete key.Binding
	Save   key.Binding
	Quit   key.Binding
	Switch key.Binding
	Top    key.Binding
	Enter  key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "Create new note"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Delete note"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "Save note"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "Quit Terminotes"),
		),
		Switch: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Switch focus"),
		),
		Top: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "Return to top"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select note"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Switch, k.New, k.Delete, k.Save, k.Top, k.Enter}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.New, k.Delete, k.Save},          //First column
		{k.Enter, k.Switch, k.Top, k.Quit}, //Second Column
	}
}
