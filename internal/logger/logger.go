package logger

import "log"

// Logger responsible for all logging actions.
type Logger struct {
	config *Config
}

// NewLogger creates new logging wrapper and returns it to caller.
func NewLogger(config *Config) *Logger {
	l := &Logger{config: config}

	return l
}

// Debugf is a wrapper around log.Printf which will work only when debug mode
// is enabled.
func (l *Logger) Debugf(template string, params ...interface{}) {
	if l.config.Debug {
		log.Printf(template, params...)
	}
}

// Debugln is a wrapper around log.Println which will work only when debug mode
// is enabled.
func (l *Logger) Debugln(params ...interface{}) {
	if l.config.Debug {
		log.Println(params...)
	}
}

// Infof is a wrapper around log.Printf.
func (l *Logger) Infof(template string, params ...interface{}) {
	log.Printf(template, params...)
}

// Infoln is a wrapper around log.Println.
func (l *Logger) Infoln(params ...interface{}) {
	log.Println(params...)
}
