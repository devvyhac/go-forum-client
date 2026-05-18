package main

import (
	"net"

	tea "charm.land/bubbletea/v2"
)

func Receive(ch chan string) tea.Cmd {
	// Receiver Goroutine (Server -> Client)
	return func() tea.Msg {
		return ServerMsg(<-ch)
	}
}

func ChatSend(conn net.Conn, text string) tea.Cmd {
	return func() tea.Msg {
		conn.Write([]byte(text + "\n"))
		return nil
	}
}
