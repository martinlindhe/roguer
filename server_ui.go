package rogue

import (
	"fmt"

	"github.com/gobuild/log"
)

var (
	generalLog messageList
)

type messageList struct {
	messages []string
}

func (m *messageList) Error(a ...interface{}) {

	// XXX mark msg type
	//m.Info(a...)
	log.Error(a...)
}

func (m *messageList) Debug(a ...interface{}) {

	// XXX mark msg type
	//m.Info(a...)
	log.Debug(a...)
}

func (m *messageList) Info(a ...interface{}) {

	s := fmt.Sprint(a...)
	log.Info(s)
	//m.messages = append(m.messages, s)
}

func (m *messageList) Infof(format string, a ...interface{}) {

	log.Info(fmt.Sprintf(format, a...))
	//m.Info(fmt.Sprintf(format, a...))
}

/*
func (m *messageList) repaintMostRecent() {

	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	y := 10

	maxMessages := 30 // XXX depending on console window height
	if maxMessages > len(m.messages) {
		maxMessages = len(m.messages)
	}

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
*/
