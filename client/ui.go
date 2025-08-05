package main

import (
	"fmt"
	"os"
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
	X     int
	Y     int
	Msg   string
	Style byte
}

func newCursor(YPos int, msg string, style byte) Cursor {
	return Cursor{
		X:     0,
		Y:     YPos,
		Msg:   msg,
		Style: style,
	}
}

func (c *Cursor) write() {
	c.positionWriter()
	fmt.Print(c.Msg)
	//c.pos(len(c.Msg)) // print after the message
	//c.pos() // print after the message
	//fmt.Print(string(c.Style))
	//c.X++
	//fmt.Println()
}

func (c Cursor) positionWriter() {
	//row := fmt.Appendf([]byte{}, "\033[%v;1H\033[%vC", c.Y, offset+c.X)
	row := fmt.Appendf([]byte{}, "\033[%v;1H\033[C", c.Y)
	os.Stdout.Write(row)
}

func control(fn func(ANSICodes)) {
	codes := ctrlCodes()
	fn(codes)
}

func hideCursor(c ANSICodes) {
	os.Stdout.Write(c.HomeCursor)
	//os.Stdout.Write(c.HideCursor)
}

func returnCursor(c ANSICodes) {
	os.Stdout.Write(c.ShowCursor)
}

func Clear() {
	control(hideCursor)
}

func Return() {
	control(returnCursor)
}

func Progress(id int, message string) {
	pos := newCursor(id, message, '.')
	pos.write()
}
