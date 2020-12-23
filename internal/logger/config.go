package logger

// Config represents logging configuration.
type Config struct {
	// Debug is a flag that indicates that we should print out debug output.
	Debug bool `yaml:"debug"`
}
