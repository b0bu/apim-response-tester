package main

import (
	"fmt"
	"os"
	"sync"
)

const (
	Reset = "\033[0m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Home  = "\033[H\033[2J"
	Hide  = "\033[?25l"
	Show  = "\033[?25h"
)

type ANSICodes struct {
	HomeCursor []byte
	HideCursor []byte
	ShowCursor []byte
}

func ctrlCodes() ANSICodes {
	return ANSICodes{
		HomeCursor: []byte(Home),
		HideCursor: []byte(Hide),
		ShowCursor: []byte(Show),
	}
}

type Cursor struct {
	Payload
	X     int
	Y     int
	Style byte
}

func newCursor(YPos int, p Payload, style byte) Cursor {
	return Cursor{
		X:       0,
		Y:       YPos,
		Style:   style,
		Payload: p,
	}
}

func (c *Cursor) write() {
	c.positionWriter()
	switch c.Status {
	case "pending":
		fmt.Print("id: " + c.ID + ", backend: " + c.Backend + ", status: " + Blue + c.Status + Reset + "\r")
	case "complete":
		fmt.Print("id: " + c.ID + ", backend: " + c.Backend + ", status: " + Green + c.Status + Reset + "\r")
	}
}

func (c Cursor) positionWriter() {
	row := fmt.Appendf([]byte{}, "\033[%v;1H", c.Y)
	os.Stdout.Write(row)
}

func control(fn func(ANSICodes)) {
	codes := ctrlCodes()
	fn(codes)
}

func hideCursor(c ANSICodes) {
	os.Stdout.Write(c.HomeCursor)
	os.Stdout.Write(c.HideCursor)
}

func returnCursor(c ANSICodes) {
	os.Stdout.Write(c.HomeCursor)
	os.Stdout.Write(c.ShowCursor)
}

func Clear() {
	control(hideCursor)
}

func Return() {
	control(returnCursor)
}

var termMu sync.Mutex

func Progress(id int, p Payload) {
	termMu.Lock()
	defer termMu.Unlock()

	pos := newCursor(id, p, '.')
	pos.write()

}

type UIState struct {
	MaxLines int
}

var state UIState
