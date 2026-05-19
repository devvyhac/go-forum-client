package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"net"
	"os"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func Chat(conn net.Conn, name string) {
	p := tea.NewProgram(initialModel(conn, name))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	conn        net.Conn
	name        string
}

func initialModel(conn net.Conn, name string) *model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		conn:        conn,
		name:        name,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		Receive(NetworkChannel),
	)
}

func (m *model) RebuildViewport() {
	content := lipgloss.NewStyle().
		Width(m.viewport.Width()).
		Render(strings.Join(m.messages, "\n"))

	m.viewport.SetContent(content)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// 1. Handle app logic FIRST
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

	case ServerMsg:
		atBottom := m.viewport.AtBottom()

		msgSplit := strings.Split(string(msg), ":")
		var message string

		if strings.HasPrefix(strings.TrimSpace(string(msg)), m.name) {
			message = m.senderStyle.Render("You: ") + strings.TrimSpace(msgSplit[1])
		} else {
			randomStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
			message = randomStyle.Render(fmt.Sprintf("%s: ", msgSplit[0])) + strings.TrimSpace(msgSplit[1])
		}

		m.messages = append(m.messages, message)
		m.RebuildViewport()

		if atBottom {
			m.viewport.GotoBottom()
		}

		return m, Receive(NetworkChannel)

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			chatMsg := m.textarea.Value()
			if strings.TrimSpace(chatMsg) == "" {
				return m, nil
			}

			m.textarea.Reset()
			return m, ChatSend(m.conn, chatMsg)
		}
	}

	// 2. ALWAYS forward msg to components (IMPORTANT)
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.textarea.View())
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
