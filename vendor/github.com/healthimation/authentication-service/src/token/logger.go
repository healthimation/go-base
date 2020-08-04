package token

import "log"

// Logger implements Printf
type Logger interface {
	Printf(format string, args ...interface{})
}

// SimpleLogger will print using the log package
type SimpleLogger struct{}

// Printf takes a format and args just like fmt.Printf
func (s SimpleLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
