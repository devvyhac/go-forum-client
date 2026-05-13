package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	conn := Connect()
	defer conn.Close()

	fmt.Println("Connected to Server!")

	netReader := bufio.NewReader(conn)
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

			conn.Write([]byte(username)) // Ensure server expects raw string or JSON here
		} else if data.Type == "auth_success" {
			fmt.Print(data.Message)
			break
		}
	}

	// Phase 2: Concurrent Operations
	done := make(chan struct{})

	// Receiver Goroutine (Server -> Client)
	go func() {
		defer close(done)
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

			if msg.Type == "chat" {
				fmt.Printf("\r%s: %s\n", msg.Sender, msg.Content)
				fmt.Printf("You~(%s): ", name) // Re-print prompt for UX
			}
		}
	}()

	// Sender Loop (Client -> Server)
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("You~(%s): ", name)
			text, err := inputReader.ReadString('\n')
			if err != nil {
				return
			}

			// Optional: Wrap in JSON if server expects Message struct
			conn.Write([]byte(text))
		}
	}
}

// func ReadFromConnection() {

// }
