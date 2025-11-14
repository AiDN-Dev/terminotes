package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

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
	currentFile string
	status      string
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
	ta.Placeholder = "Select a note to view and edit its contents"

	return model{
		list:        l,
		textarea:    ta,
		files:       files,
		listFocused: true, //Start with list focused
		currentFile: "",
		status:      "Select a file to view and edit.",
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

		//Update list and viewport sizes
		m.list.SetSize(m.width/3, m.height-statusBarHeight)
		m.textarea.SetWidth(m.width * 2 / 3)
		m.textarea.SetHeight(m.height - statusBarHeight)

	case tea.KeyMsg:
		if msg.String() == "tab" {
			m.listFocused = !m.listFocused
			if m.listFocused {
				m.textarea.Blur()
			} else {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// Handle other key presses
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+n":
			filename := "note-" + time.Now().Format("2006-01-02-15-04-05") + ".md"
			m.currentFile = filename
			m.textarea.SetValue("")
			m.status = "New note: " + filename
			m.listFocused = false
			cmd = m.textarea.Focus()
			cmds = append(cmds, cmd)
		case "ctrl+t":
			m.textarea.CursorStart()
			m.status = "Moved cursor to top"
		case "enter":
			if m.listFocused {
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					m.currentFile = selectedItem.title
					content, err := os.ReadFile("./notes/" + m.currentFile)
					if err != nil {
						m.status = "Error reading file: " + err.Error()
					} else {
						m.textarea.SetValue(string(content))
						m.textarea.CursorStart() // Sets the cursor to the start of the file
						m.status = "Editing " + m.currentFile
					}
				}
			}
		case "ctrl+s":
			if !m.listFocused && m.currentFile != "" {
				//create the notes directory if it exists
				if err := os.MkdirAll("./notes", 0755); err != nil {
					m.status = "Error creating director: " + err.Error()
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

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.textarea.View())

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar)
}

func main() {
	if _, err := tea.NewProgram(containerModel()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
