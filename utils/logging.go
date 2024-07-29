package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var Mylog log.Logger

// Initialise MyLog.
// Override standard log defaults.
// Write to LogFile instead of terminal.
func LogInit() {
	Mylog = *log.Default()

	f, err := os.OpenFile("./logging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	Mylog.SetOutput(f)
}

// Wrapper for the standard library log.
// Adds a colour that can be used to identify the caller easily.
//
//	startColour: selects a colour from utils.colour rendered using ANSI codes.
//	message: the message to display
//	returns: the formatted message for further reporting such as browser output
func TraceInfo(startColour string, message string) string {
	formattedMessage := "[TRACE] " + startColour + message + Reset

	// log the message in the log file
	Mylog.Output(1, message) //TODO find out how to print ANSI escapes to a log file

	// (Development only) send the message to the terminal
	fmt.Println(formattedMessage)

	// pass the message back to the caller
	return formattedMessage
}

func TraceError(message string) error {
	formattedMessage := "[ERROR] " + BrightRed + message + Reset

	// log the message in the log file
	Mylog.Output(2, message)

	// (Development only) send the message to the terminal
	fmt.Println(formattedMessage)

	// pass the message back to the caller
	return errors.New(formattedMessage)
}

// convenience function to include a variable in an TraceError message
//
//	startColour: the colour of the message
//	message: the message, which must contain one formatting symbol
//	details: the variable to be printed
func TraceErrorf(message string, details ...any) error {
	return TraceError(fmt.Sprintf(message, details...))
}

// convenience function to include a variable in a TraceInfo message
//
//	startColour: the colour of the message
//	message: the message, which must contain one formatting symbol
//	details: the variable to be printed
func TraceInfof(startColour string, message string, details ...any) string {
	return TraceInfo(startColour, fmt.Sprintf(message, details...))
}

// convenience function to pretty print structs, maps, etc
//
//	m: a description of the object which must contain one formatting symbol
//	o: the object
func TracePretty(m string, o any) string {
	b, _ := json.MarshalIndent(o, " ", " ")
	return TraceInfof(BrightWhite, m, string(b))
}
