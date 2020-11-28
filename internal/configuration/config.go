package configuration

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	ErrConfigurationFileDoesNotExist  = errors.New("configuration file does not exist")
	ErrConfigurationFilePathUndefined = errors.New("configuration file path wasn't provided")
)

// Config is an application's configuration.
type Config struct {
	configPath string
	// Applications describes configuration for remote application's endpoints.
	// Key is an application's name.
	Applications map[string]struct {
		// Endpoint is a remote application endpoint which should give us metrics
		// in Prometheus format.
		Endpoint string
		// TimeBetweenRequests is a minimal amount of time which should pass
		// between requests.
		TimeBetweenRequests time.Duration
	}
	// Datastore describes data storage configuration.
	Datastore struct {
		// ValidTimeout is a timeout for which every data entry will be considered
		// as valid. After that timeout if value wasn't updated it will be considered
		// as invalid and purged from memory.
		ValidTimeout time.Duration `yaml:"valid_timeout"`
	} `yaml:"datastore"`
}

// NewConfig returns new configuration.
func NewConfig() *Config {
	c := &Config{}
	c.initialize()

	return c
}

// Initializes configuration.
func (c *Config) initialize() {
	flag.StringVar(&c.configPath, "config", "", "Configuration file path.")
}

// Parse parses configuration.
func (c *Config) Parse() error {
	if c.configPath == "" {
		return ErrConfigurationFilePathUndefined
	}

	if strings.HasPrefix(c.configPath, "~") {
		userDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("%s: %w", "file path normalization: getting user's home directory", err)
		}

		c.configPath = strings.Replace(c.configPath, "~", userDir, 1)
	}

	cfgPath, err := filepath.Abs(c.configPath)
	if err != nil {
		return fmt.Errorf("%s: %w", "file path normalization: getting absolute path", err)
	}

	c.configPath = cfgPath

	if c.configPath == "" {
		return fmt.Errorf("%s: %w", "file path normalization", ErrConfigurationFilePathUndefined)
	}

	fileData, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return fmt.Errorf("%s: %w", "configuration file read", err)
	}

	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		return fmt.Errorf("%s: %w", "configuration file parsing", err)
	}

	return nil
}
