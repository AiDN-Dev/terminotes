package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list        list.Model
	textarea    textarea.Model
	files       []fs.DirEntry
	listFocused bool //true: list is focused, false: text area is focused
	currentFile string
	status      string
	width       int
	height      int
	quitting    bool
	keys        keyMap
	help        help.Model
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
	l.SetShowHelp(false) //Using my own help view

	//setup the text area
	ta := textarea.New()
	ta.Placeholder = "Select a note to view and edit its contents"

	keys := newKeyMap()
	help := help.New()
	help.ShowAll = true

	return model{
		list:        l,
		textarea:    ta,
		files:       files,
		listFocused: true, //Start with list focused
		currentFile: "",
		status:      "Select a file to view and edit.",
		keys:        keys,
		help:        help,
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
		statusBarHeight := 1

		//Calculate help height
		m.help.Width = m.width
		helpView := m.help.View(m.keys)
		helpHeight := lipgloss.Height(helpView)

		//Update list and viewport sizes
		m.list.SetSize(m.width/3, m.height-statusBarHeight-helpHeight)
		m.textarea.SetWidth(m.width * 2 / 3)
		m.textarea.SetHeight(m.height - statusBarHeight - helpHeight)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Switch):
			m.listFocused = !m.listFocused
			if m.listFocused {
				m.textarea.Blur()
			} else {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		if m.listFocused {
			switch {
			case key.Matches(msg, m.keys.New):
				filename := "note-" + time.Now().Format("2006-01-02-15-04-05") + ".md"
				m.currentFile = filename
				m.textarea.SetValue("")
				m.status = "New note: " + filename
				m.listFocused = false
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, m.keys.Enter):
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					m.currentFile = selectedItem.title
					content, err := os.ReadFile("./notes/" + m.currentFile)
					if err != nil {
						m.status = "Error reading file: " + err.Error()
					} else {
						m.textarea.SetValue(string(content))
						m.textarea.CursorStart()
						m.status = "Editing: " + m.currentFile
					}
				}
			}
		} else {
			switch {
			case key.Matches(msg, m.keys.Save):
				if m.currentFile != "" {
					//Create notes directory if it exists
					if err := os.MkdirAll("./notes", 0755); err != nil {
						m.status = "Error creating directory: " + err.Error()
						return m, nil
					}

					err := os.WriteFile("./notes/"+m.currentFile, []byte(m.textarea.Value()), 0644)
					if err != nil {
						m.status = "Error saving file: " + err.Error()
					} else {
						m.status = "Saved: " + m.currentFile
					}
				} else {
					m.status = "Cannot Save: focus the textarea and try to open a file first."
				}
			case key.Matches(msg, m.keys.Top):
				m.textarea.CursorStart()
				m.status = "Movd cursor to Top"
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

	// Create a status bar
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#5C5C5C")).
		Padding(0, 1).
		Render(m.status)

	helpView := m.help.View(m.keys)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.textarea.View())

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar, helpView)
}

func main() {
	if _, err := tea.NewProgram(containerModel()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
