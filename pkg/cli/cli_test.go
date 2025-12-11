package cli

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "default values",
			args: []string{},
			expected: Config{
				Recursive:    false,
				Relative:     false,
				MarkerPrefix: "<<<",
				MarkerSuffix: ">>>",
			},
		},
		{
			name: "recursive flag",
			args: []string{"--recursive"},
			expected: Config{
				Recursive:    true,
				Relative:     false,
				MarkerPrefix: "<<<",
				MarkerSuffix: ">>>",
			},
		},
		{
			name: "relative flag",
			args: []string{"--relative"},
			expected: Config{
				Recursive:    false,
				Relative:     true,
				MarkerPrefix: "<<<",
				MarkerSuffix: ">>>",
			},
		},
		{
			name: "custom markers",
			args: []string{"--marker-prefix", "[[[", "--marker-suffix", "]]]"},
			expected: Config{
				Recursive:    false,
				Relative:     false,
				MarkerPrefix: "[[[",
				MarkerSuffix: "]]]",
			},
		},
		{
			name: "all flags combined",
			args: []string{"--recursive", "--relative", "--marker-prefix", "***", "--marker-suffix", "---"},
			expected: Config{
				Recursive:    true,
				Relative:     true,
				MarkerPrefix: "***",
				MarkerSuffix: "---",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset pflag state
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			// Set up args
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			// Parse and test
			cfg := Parse()
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestPatterns(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "no patterns",
			args:     []string{},
			expected: []string{},
		},
		{
			name:     "single pattern",
			args:     []string{"*.go"},
			expected: []string{"*.go"},
		},
		{
			name:     "multiple patterns",
			args:     []string{"*.go", "*.txt", "docs/"},
			expected: []string{"*.go", "*.txt", "docs/"},
		},
		{
			name:     "patterns with flags",
			args:     []string{"--recursive", "*.go", "--relative", "*.txt"},
			expected: []string{"*.go", "*.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset pflag state
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			// Set up args
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			// Parse first to consume flags, then get patterns
			Parse()
			patterns := Patterns()
			assert.Equal(t, tt.expected, patterns)
		})
	}
}

func TestConfig(t *testing.T) {
	t.Run("Config struct fields", func(t *testing.T) {
		cfg := Config{
			Recursive:    true,
			Relative:     false,
			MarkerPrefix: "START",
			MarkerSuffix: "END",
		}

		assert.True(t, cfg.Recursive)
		assert.False(t, cfg.Relative)
		assert.Equal(t, "START", cfg.MarkerPrefix)
		assert.Equal(t, "END", cfg.MarkerSuffix)
	})
}

// TestParseWithInvalidFlags tests edge cases with flag parsing
func TestParseEdgeCases(t *testing.T) {
	t.Run("empty marker values", func(t *testing.T) {
		// Reset pflag state
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		oldArgs := os.Args
		os.Args = []string{"test", "--marker-prefix", "", "--marker-suffix", ""}
		defer func() { os.Args = oldArgs }()

		cfg := Parse()
		assert.Equal(t, "", cfg.MarkerPrefix)
		assert.Equal(t, "", cfg.MarkerSuffix)
	})
}
