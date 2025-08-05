package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type ANSICodes struct {
	HomeCursor []byte
	HideCursor []byte
	ShowCursor []byte
}

func ctrlCodes() ANSICodes {
	return ANSICodes{
		HomeCursor: []byte("\033[H\033[2J"),
		HideCursor: []byte("\033[?25l"),
		ShowCursor: []byte("\033[?25h"),
	}
}

type Cursor struct {
	X int
	Y int
}

func newCursor(YPos int) Cursor {
	return Cursor{
		X: 0,
		Y: YPos,
	}
}

const message string = "pending "

func (c *Cursor) write(t byte) {
	n := rand.Intn(5)
	fmt.Print(message)
	for range n {
		time.Sleep(time.Duration(1) * time.Second)
		c.pos()
		fmt.Print(string(t))
		c.X++
	}
	fmt.Println()
}

func (c Cursor) pos() {
	row := fmt.Appendf([]byte{}, "\033[%v;1H\033[%vC", c.Y, len(message)+c.X)
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
	os.Stdout.Write(c.ShowCursor)
}

func Progress(id int, c byte) {
	control(hideCursor)

	pos := newCursor(id)
	pos.write(c)

	defer control(returnCursor)
}
