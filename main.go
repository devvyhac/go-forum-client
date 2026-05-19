package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ServerMsg string

var NetworkChannel = make(chan string, 100)

func main() {
	conn := Connect()
	defer conn.Close()

	fmt.Println("Connected to Server!")

	netReader := bufio.NewReader(conn)

	// textarea below here.
	inputReader := bufio.NewReader(os.Stdin)
	var name string

	// Phase 1: Synchronous Handshake/Authentication
	// We handle this sequentially before launching concurrent loops
	for {
		line, err := netReader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection lost during handshake.")
			return
		}

		var data Handshake
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		if data.Type == "prompt" {
			fmt.Print(data.Message)
			username, _ := inputReader.ReadString('\n')
			name = strings.TrimSpace(username)

			conn.Write([]byte(name + "\n")) // Ensure server expects raw string or JSON here
		} else if data.Type == "auth_success" {
			inputReader = nil
			break
		}
	}

	// Phase 2: Concurrent Operations

	// Receiver Goroutine (Server -> Client)
	go func() {
		defer conn.Close()
		for {
			line, err := netReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server.")
				return
			}

			var msg Message
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				continue
			}

			if msg.Type == "chat" || msg.Type == "event" {
				NetworkChannel <- fmt.Sprintf("\r%s: %s\n", msg.Sender, msg.Content)
			} else if msg.Type == "command" {
				NetworkChannel <- fmt.Sprintf("\r%s: %s\n", msg.Sender, msg.Content)
			}
		}
	}()

	os.Stdout.Sync()
	os.Stdin.Sync()

	Chat(conn, name)
}
