package main

import "time"

type Message struct {
	Type    string
	Sender  string
	Content string
	Time    time.Time
}

type Handshake struct {
	Type    string
	Status  string
	Message string
}
