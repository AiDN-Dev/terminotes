package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"

	//"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list        list.Model
	textarea    textarea.Model
	files       []fs.DirEntry
	listFocused bool //true: list is focused, false: text area is focused
	width       int
	height      int
	quitting    bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func containerModel() model {
	//Read files from the current directory
	files, err := os.ReadDir("./notes")
	if err != nil {
		//Handle error ....
		fmt.Println("Error Could not read directory")
	}

	items := make([]list.Item, len(files))
	for i, file := range files {
		items[i] = item{title: file.Name(), desc: ""}
	}

	//Setup the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Notes"

	//setup the text area
	ta := textarea.New()
	ta.Placeholder = "Select a file to view and edit its contents"

	return model{
		list:        l,
		textarea:    ta,
		files:       files,
		listFocused: true, //Start with list focused
	}
}

// Item is a helper struct for the list component
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		//Update list and viewport sizes
		m.list.SetSize(m.width/3, m.height)
		m.textarea.SetWidth(m.width * 2 / 3)
		m.textarea.SetHeight(m.height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.listFocused = !m.listFocused
			if m.listFocused {
				m.textarea.Blur()
			} else {
				m.textarea.Focus()
			}
		case "enter":
			if m.listFocused {
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					content, err := os.ReadFile("./notes/" + selectedItem.title)
					if err != nil {
						m.textarea.SetValue("Error reading file: " + err.Error())
					} else {
						m.textarea.SetValue(string(content))
					}
				}
			}
		}
	}
	if m.listFocused {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.textarea.View())
}

func main() {
	if _, err := tea.NewProgram(containerModel()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
