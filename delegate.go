package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newItemDelegate(styles Styles) list.ItemDelegate {
	return &itemDelagte{styles: styles}
}

type itemDelagte struct {
	styles Styles
}

func (d *itemDelagte) Height() int                               { return 1 }
func (d *itemDelagte) Spacing() int                              { return 0 }
func (d *itemDelagte) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d *itemDelagte) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	s, ok := listItem.(item)
	if !ok {
		return
	}

	title := s.Title()
	maxWidth := m.Width() - 4

	if len(title) > maxWidth && maxWidth > 0 {
		title = title[:maxWidth-3] + "..."
	}

	var style lipgloss.Style
	if index == m.Index() {
		style = d.styles.ListItemActive
	} else {
		style = d.styles.ListItem
	}

	fmt.Fprint(w, style.Render(title))
}
