package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MarinX/keylogger"
)

type EventType int

const (
	eventKeyUp   = 0
	eventKeyDown = 1
)

const privacy = true

type Logger struct {
	writer io.Writer
}

func (l Logger) log(t time.Time, e EventType, key string) error {
	fmt.Fprintf(l.writer, "%d\t%d\t%s\n", t.Unix(), e, key)
	return nil
}

func main() {

	logger := Logger{
		os.Stdout,
	}

	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		panic("No keyboard found...you will need to provide manual input path")
	}

	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	var keyCode string
	var e keylogger.InputEvent

	// range of events
	for e = range k.Read() {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:
			if privacy {
				keyCode = keyCodeMap[e.Code]
			} else {
				keyCode = e.KeyString()
			}
			logger.log(time.Now(), EventType(e.Value), keyCode)
		}
	}
}
