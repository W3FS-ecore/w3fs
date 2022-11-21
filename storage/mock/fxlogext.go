package mock

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var _exit = func() { os.Exit(1) }

// Printer is a formatting printer.
type Printer interface {
	Printf(string, ...interface{})
}

// New returns a new Logger backed by the standard library's log package.
func NewLogger() *Logger {
	return &Logger{log.New(os.Stderr, "", log.LstdFlags)}
}

// A Logger writes output to standard error.
type Logger struct {
	Printer
}

// Printf logs a formatted Fx line.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Printer.Printf(prepend(format), v...)
}

// PrintProvide logs a type provided into the dig.Container.
func (l *Logger) PrintProvide(t interface{}) {
	//for _, rtype := range fxreflect.ReturnTypes(t) {
	//	l.Printf("PROVIDE\t%s <= %s", rtype, fxreflect.FuncName(t))
	//}
	l.Printf("PROVIDE\t%s ", t)
}

// PrintSignal logs an os.Signal.
func (l *Logger) PrintSignal(signal os.Signal) {
	l.Printf(strings.ToUpper(signal.String()))
}

// Panic logs an Fx line then panics.
func (l *Logger) Panic(err error) {
	l.Printer.Printf(prepend(err.Error()))
	panic(err)
}

// Fatalf logs an Fx line then fatals.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Printer.Printf(prepend(format), v...)
	_exit()
}

func prepend(str string) string {
	return fmt.Sprintf("[Fx] %s", str)
}
