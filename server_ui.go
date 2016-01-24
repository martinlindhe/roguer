package rogue

import (
	"fmt"

	"github.com/gobuild/log"
	"github.com/nsf/termbox-go"
)

var (
	generalLog messageList
)

type messageList struct {
	messages []string
}

func (m *messageList) Error(a ...interface{}) {

	// XXX mark msg type
	m.Info(a...)
}

func (m *messageList) Debug(a ...interface{}) {

	// XXX mark msg type
	m.Info(a...)
}

func (m *messageList) Info(a ...interface{}) {

	s := fmt.Sprint(a...)
	log.Info(s) // XXX
	m.messages = append(m.messages, s)
}

func (m *messageList) Infof(format string, a ...interface{}) {

	//m.messages = append(m.messages, fmt.Sprintf(format, a...))
	m.Info(fmt.Sprintf(format, a...))
}

func (m *messageList) repaintMostRecent() {

	y := 10

	maxMessages := 30 // XXX depending on console window height

	// draw last x lines of log
	msgs := m.messages[len(m.messages)-maxMessages : len(m.messages)]

	for _, m := range msgs {
		tbPrint(0, y, termbox.ColorWhite, termbox.ColorDefault, m)
		y++
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
			//log.SetLevel(log.DebugLevel)
		}
	}

	return true
}
