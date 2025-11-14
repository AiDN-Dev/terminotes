package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list         list.Model
	textarea     textarea.Model
	files        []fs.DirEntry
	listFocused  bool //true: list is focused, false: text area is focused
	currentFile  string
	status       string
	width        int
	height       int
	quitting     bool
	keys         keyMap
	help         help.Model
	textInput    textinput.Model //for file name input
	inputting    bool            //true if currently askign for a file Name
	InputPurpose string          // "New" or "Save"
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

	// key mapping from keys.go
	keys := newKeyMap()
	help := help.New()
	help.ShowAll = true

	// Text inputting
	ti := textinput.New()
	ti.Placeholder = "Enter filename (leave blank for default)"
	ti.CharLimit = 50
	ti.Width = 40
	ti.Prompt = "Filename: "

	return model{
		list:         l,
		textarea:     ta,
		files:        files,
		listFocused:  true, //Start with list focused
		currentFile:  "",
		status:       "Select a file to view and edit.",
		keys:         keys,
		help:         help,
		textInput:    ti,
		inputting:    false, //not inputting initially
		InputPurpose: "",    //no purpose initially
	}
}

// Item is a helper struct for the list component
type item struct {
	title, desc string
}

func refreshList() ([]list.Item, []fs.DirEntry, error) {
	files, err := os.ReadDir("./notes")
	if err != nil {
		return nil, nil, err
	}

	items := make([]list.Item, len(files))
	for i, file := range files {
		items[i] = item{title: file.Name(), desc: ""}
	}
	return items, files, nil
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Handle filename input if active
	if m.inputting {
		// Pass all messages to the text input model first
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd) // Collect any commands from textinput

		// Now handle specific key presses for the text input
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.Type {
			case tea.KeyEnter:
				filename := m.textInput.Value()
				if filename == "" {
					filename = "note-" + time.Now().Format("2006-01-02-15-04-05") + ".md"
				} else {
					if !strings.HasSuffix(filename, ".md") {
						filename += ".md"
					}
				}

				if m.InputPurpose == "new" {
					m.currentFile = filename
					m.textarea.SetValue("") // Start with empty content for new note
					m.status = "New note: " + filename
					m.listFocused = false // Focus textarea for editing

					if err := os.MkdirAll("./notes", 0755); err != nil {
						m.status = "Error creating directory: " + err.Error()
					} else {
						err := os.WriteFile("./notes/"+m.currentFile, []byte(""), 0644)
						if err != nil {
							m.status = "Error creating file: " + err.Error()
						} else {
							items, files, err := refreshList()
							if err != nil {
								m.status = "Error refreshing list: " + err.Error()
							} else {
								m.list.SetItems(items)
								m.files = files
							}
						}
					}
					cmd = m.textarea.Focus()
					cmds = append(cmds, cmd) // Add focus command
				} else if m.InputPurpose == "save" {
					m.currentFile = filename
					m.status = "Saving as: " + filename

					if err := os.MkdirAll("./notes", 0755); err != nil {
						m.status = "Error creating directory: " + err.Error()
					} else {
						err := os.WriteFile("./notes/"+m.currentFile, []byte(m.textarea.Value()), 0644)
						if err != nil {
							m.status = "Error saving file: " + err.Error()
						} else {
							m.status = "Saved: " + m.currentFile
							items, files, err := refreshList()
							if err != nil {
								m.status = "Error refreshing list: " + err.Error()
							} else {
								m.list.SetItems(items)
								m.files = files
							}
						}
					}
				}

				m.inputting = false
				m.InputPurpose = ""
				m.textInput.Blur()
				return m, tea.Batch(cmds...)
			case tea.KeyEsc:
				m.inputting = false
				m.InputPurpose = ""
				m.textInput.Blur()
				m.status = "Filename input cancelled."
				return m, tea.Batch(cmds...)
			}
		}
		return m, tea.Batch(cmds...)
	}

	// Main update logic when not inputting
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		statusBarHeight := 1

		m.help.Width = m.width
		helpView := m.help.View(m.keys)
		helpHeight := lipgloss.Height(helpView)

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
				m.inputting = true
				m.InputPurpose = "new"
				m.textInput.Reset()
				cmd = m.textInput.Focus()
				return m, cmd
			case key.Matches(msg, m.keys.Delete):
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					if selectedItem.title == m.currentFile {
						m.textarea.SetValue("")
						m.currentFile = ""
					}
					err := os.Remove("./notes/" + selectedItem.title)
					if err != nil {
						m.status = "Error deleting file: " + err.Error()
					} else {
						m.status = "Deleted: " + selectedItem.title
						items, files, err := refreshList()
						if err != nil {
							m.status = "Error refreshing list: " + err.Error()
						} else {
							m.list.SetItems(items)
							m.files = files
						}
					}
				}
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
						m.status = "Editing " + m.currentFile
					}
				}
			}
		} else { // textarea focused
			switch {
			case key.Matches(msg, m.keys.Save):
				if m.currentFile != "" {
					if err := os.MkdirAll("./notes", 0755); err != nil {
						m.status = "Error creating directory: " + err.Error()
						return m, nil
					}
					err := os.WriteFile("./notes/"+m.currentFile, []byte(m.textarea.Value()), 0644)
					if err != nil {
						m.status = "Error saving file: " + err.Error()
					} else {
						m.status = "Saved: " + m.currentFile
						items, files, err := refreshList()
						if err != nil {
							m.status = "Error refreshing list: " + err.Error()
						} else {
							m.list.SetItems(items)
							m.files = files
						}
					}
				} else {
					m.inputting = true
					m.InputPurpose = "save"
					m.textInput.Reset()
					cmd = m.textInput.Focus()
					return m, cmd
				}
			case key.Matches(msg, m.keys.Top):
				m.textarea.CursorStart()
				m.status = "Moved cursor to top"
			}
		}
	}

	// Update list and textarea if not inputting
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

	//Create the status bar
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#5C5C5C")).
		Padding(0, 1).
		Render(m.status)

	helpView := m.help.View(m.keys)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.textarea.View())
	appView := lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar, helpView)

	if m.inputting {
		var (
			modalStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("62")).
					Padding(1, 2).
					Background(lipgloss.Color("#000000")).
					Foreground(lipgloss.Color("#FFFFFF"))

			promptStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("240"))
		)

		m.textInput.PromptStyle = promptStyle
		m.textInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

		inputBoxContent := lipgloss.JoinVertical(lipgloss.Left,
			"Enter filename: ",
			m.textInput.View(),
			promptStyle.Render("Press Enter to confirm, Esc to cancel"),
		)

		inputBox := modalStyle.Render(inputBoxContent)

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			inputBox,
		)
	}

	return appView
}

func main() {
	if _, err := tea.NewProgram(containerModel()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
