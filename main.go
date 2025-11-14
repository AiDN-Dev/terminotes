package main

import (
	"encoding/json"
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

const (
	viewNotes = iota
	viewSettings
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
	styles       Styles
	themes       map[string]Theme
	appState     int
	settingsList list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func containerModel() model {
	loadThemes, err := loadThemes()
	if err != nil {
		fmt.Println("Could not load themes: ", err)
		os.Exit(1)
	}

	defaultTheme := loadThemes["mocha.json"]
	styles := NewStyles(defaultTheme)

	settingsItems := make([]list.Item, 0, len(loadThemes))
	for _, theme := range loadThemes {
		settingsItems = append(settingsItems, item{title: theme.Name, desc: theme.FileName})
	}

	settingsDelegate := list.NewDefaultDelegate()
	settingsDelegate.ShowDescription = false
	settingsDelegate.SetSpacing(0)
	settingsDelegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		return nil
	}

	sl := list.New(settingsItems, settingsDelegate, 0, 0)
	sl.Title = "Themes"
	sl.SetShowHelp(false)

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
	delegate := newItemDelegate(styles)
	l := list.New(items, delegate, 0, 0)
	l.Title = "Notes"
	l.Styles.Title = styles.ListTitle
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
	ti.Prompt = ""

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
		styles:       styles,
		themes:       loadThemes,
		appState:     viewNotes,
		settingsList: sl,
	}
}

// Item is a helper struct for the list component
type item struct {
	title, desc string
}

func loadThemes() (map[string]Theme, error) {
	themes := make(map[string]Theme)
	files, err := os.ReadDir("./themes")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := "./themes/" + file.Name()
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading theme file %s: %v\n", filePath, err)
				continue
			}

			var theme Theme
			err = json.Unmarshal(content, &theme)
			if err != nil {
				fmt.Printf("Error parsing theme file %s: %v\n", filePath, err)
				continue
			}
			theme.FileName = file.Name()
			if theme.Name == "" {
				theme.Name = strings.TrimSuffix(file.Name(), ".json")
			}
			themes[file.Name()] = theme
		}
	}
	return themes, nil
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

	// Global keybindings that work in any view
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Settings):
			if m.appState == viewNotes { // Only enter settings from notes view
				m.appState = viewSettings
				m.status = "Settings - Press Enter to select a theme, Esc to return"
				return m, nil
			}
		}
	}

	// Handle updates based on the application state
	switch m.appState {
	case viewSettings:
		// Handle key presses specific to the settings view
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(msg, m.keys.Enter):
				selectedItem, ok := m.settingsList.SelectedItem().(item)
				if ok {
					themeFileName := selectedItem.desc
					if selectedTheme, found := m.themes[themeFileName]; found {
						m.styles = NewStyles(selectedTheme)
						m.status = fmt.Sprintf("Theme changed to %s", selectedTheme.Name)
					} else {
						m.status = fmt.Sprintf("Error: Theme %s not found", themeFileName)
					}
				}
				m.appState = viewNotes
				return m, nil
			case msg.Type == tea.KeyEsc:
				// Exit settings view
				m.appState = viewNotes
				m.status = "Select a file to view and edit."
				return m, nil
			}
		}

		// Pass updates to the settings list
		m.settingsList, cmd = m.settingsList.Update(msg)
		cmds = append(cmds, cmd)

	default: // This is viewNotes
		// Handle updates specific to the notes view
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

			const listWidth = 30

			paneStyle := m.styles.ActivePane
			borderWidth := paneStyle.GetHorizontalFrameSize()
			borderHeight := paneStyle.GetVerticalFrameSize()

			m.list.SetSize(listWidth-borderWidth, m.height-statusBarHeight-helpHeight-borderHeight)

			textareaWidth := m.width - m.list.Width() - borderWidth*2
			m.textarea.SetWidth(textareaWidth)
			m.textarea.SetHeight(m.height - statusBarHeight - helpHeight - borderHeight)

			// Also size the settings list
			m.settingsList.SetSize(m.width/2, m.height/2)

		case tea.KeyMsg:
			switch {
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
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Handle the view based on the application state
	switch m.appState {
	case viewSettings:
		// Settings View
		// We'll place the settings list in the center of the screen.
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.settingsList.View(),
		)

	default: // This is viewNotes
		// If inputting, show the input modal over the notes view
		if m.inputting {
			m.textInput.PromptStyle = m.styles.Prompt
			m.textInput.TextStyle = m.styles.TextInput

			inputBoxContent := lipgloss.JoinVertical(lipgloss.Left,
				"Enter filename: ",
				m.textInput.View(),
				m.styles.Prompt.Render("Press Enter to confirm, Esc to cancel."),
			)

			inputBox := m.styles.Modal.Render(inputBoxContent)

			return lipgloss.Place(
				m.width,
				m.height,
				lipgloss.Center,
				lipgloss.Center,
				inputBox,
			)
		}

		// Default Notes View (when not inputting)
		statusBar := m.styles.StatusBar.Padding(0, 1).Width(m.width).Render(m.status)
		helpView := m.help.View(m.keys)

		var listStyle, textareaStyle lipgloss.Style
		if m.listFocused {
			listStyle = m.styles.ActivePane
			textareaStyle = m.styles.InactivePane
		} else {
			listStyle = m.styles.InactivePane
			textareaStyle = m.styles.ActivePane
		}
		listStyle = listStyle.Width(m.list.Width())

		mainContent := lipgloss.JoinHorizontal(lipgloss.Top,
			listStyle.Render(m.list.View()),
			textareaStyle.Render(m.textarea.View()),
		)
		return lipgloss.JoinVertical(lipgloss.Left, mainContent, statusBar, helpView)
	}
}

func main() {
	if _, err := tea.NewProgram(containerModel()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
