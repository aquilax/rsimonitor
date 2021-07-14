package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MarinX/keylogger"
)

type EventType int
type Verbosity int

const (
	verbosityNone Verbosity = iota
	verbosityPrivacy
	verbosityNames
	verbosityCodes
)

const privacy = true

type Logger struct {
	writer    io.Writer
	verbosity Verbosity
}

func (l Logger) log(t time.Time, e keylogger.InputEvent) error {
	var err error
	switch l.verbosity {
	case verbosityPrivacy:
		_, err = fmt.Fprintf(l.writer, "%d\t%d\t%s\n", t.Unix(), e.Value, keyCodeMap[e.Code])
	case verbosityNames:
		_, err = fmt.Fprintf(l.writer, "%d\t%d\t%s\n", t.Unix(), e.Value, e.KeyString())
	case verbosityCodes:
		_, err = fmt.Fprintf(l.writer, "%d\t%d\t%s\t%d\n", t.Unix(), e.Value, e.KeyString(), e.Code)
	default:
		_, err = fmt.Fprintf(l.writer, "%d\t%d\n", t.Unix(), e.Value)
	}
	return err
}

func main() {
	verbosity := verbosityNone
	v1 := flag.Bool("v", false, "Log keys but anonymize letters and digits")
	v2 := flag.Bool("vv", false, "Log raw characters")
	v3 := flag.Bool("vvv", false, "Log raw characters and key codes")
	flag.Parse()
	if *v1 {
		verbosity = verbosityPrivacy
	}
	if *v2 {
		verbosity = verbosityNames
	}
	if *v3 {
		verbosity = verbosityCodes
	}

	logger := Logger{
		os.Stdout,
		verbosity,
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

	var e keylogger.InputEvent

	// range of events
	for e = range k.Read() {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:
			logger.log(time.Now(), e)
		}
	}
}
