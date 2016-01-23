package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type MessageList struct {
	messages []string
}

func (m MessageList) Info(a ...interface{}) {

	m.messages = append(m.messages, fmt.Sprint(a))
}

func (m MessageList) Infof(format string, a ...interface{}) {

	m.messages = append(m.messages, fmt.Sprintf(format, a))
}

func (m MessageList) repaintMostRecent() {

	y := 0

	// XXX draw last x lines, depending on console window height
	for _, m := range logMessages.messages {
		tbPrint(0, y, termbox.ColorWhite, termbox.ColorDefault, m)
		fmt.Println(m)
		y++ // XXX
	}

	termbox.Flush()
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// returns false to signal shutdown
func handleEvents() bool {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		if ev.Key == termbox.KeyCtrlQ {
			fmt.Println("ctrl-q pressed, exiting")
			return false
		}
		if ev.Key == termbox.KeyCtrlX {
			fmt.Println("ctrl-x pressed, exiting")
			return false
		}
		if ev.Key == 'v' {
			fmt.Println("xxx verbosity")
		}
	}

	return true
}
