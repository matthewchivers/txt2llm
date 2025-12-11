// Package cli handles command-line argument parsing and configuration for txt2llm.
package cli

import (
	"github.com/spf13/pflag"
)

// Config holds parsed CLI flags.
type Config struct {
	Recursive    bool
	Relative     bool
	MarkerPrefix string
	MarkerSuffix string
}

// Parse parses command-line flags and returns configuration.
func Parse() Config {
	var cfg Config
	pflag.BoolVar(&cfg.Recursive, "recursive", false, "Process directories recursively")
	pflag.BoolVar(&cfg.Relative, "relative", false, "Use paths relative to current directory in output")
	pflag.StringVar(&cfg.MarkerPrefix, "marker-prefix", "<<<", "Prefix for start/end marker lines")
	pflag.StringVar(&cfg.MarkerSuffix, "marker-suffix", ">>>", "Suffix for start/end marker lines")
	pflag.CommandLine.SetInterspersed(true)
	pflag.Parse()
	return cfg
}

// Patterns returns positional arguments treated as patterns.
func Patterns() []string {
	return append([]string{}, pflag.Args()...)
}
