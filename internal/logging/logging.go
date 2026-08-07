// Package logging normalizes the subset of mritd/logger configuration that is
// exposed through CLI flags (level, output format). Keeping the parsing here
// makes it unit-testable without touching the process-wide logger singleton.
package logging

import (
	"fmt"
	"strings"

	"github.com/mritd/logger"
)

// Level is a supported log level.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Format is a supported log output format.
type Format string

const (
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
)

var levelValues = map[Level]bool{
	LevelDebug: true,
	LevelInfo:  true,
	LevelWarn:  true,
	LevelError: true,
}

var formatValues = map[Format]bool{
	FormatConsole: true,
	FormatJSON:    true,
}

// ParseLevel validates and normalizes a --log-level value. Empty means the
// default (info).
func ParseLevel(s string) (Level, error) {
	if s == "" {
		return LevelInfo, nil
	}
	l := Level(strings.ToLower(strings.TrimSpace(s)))
	if !levelValues[l] {
		return "", fmt.Errorf("invalid log level %q (want debug|info|warn|error)", s)
	}
	return l, nil
}

// ParseFormat validates and normalizes a --log-format value. Empty means the
// default (console).
func ParseFormat(s string) (Format, error) {
	if s == "" {
		return FormatConsole, nil
	}
	f := Format(strings.ToLower(strings.TrimSpace(s)))
	if !formatValues[f] {
		return "", fmt.Errorf("invalid log format %q (want console|json)", s)
	}
	return f, nil
}

// ApplyLevel maps a parsed Level onto the mritd logger.
func ApplyLevel(l Level) {
	switch l {
	case LevelDebug:
		logger.SetLevel(logger.LevelDebug)
	case LevelInfo:
		logger.SetLevel(logger.LevelInfo)
	case LevelWarn:
		logger.SetLevel(logger.LevelWarn)
	case LevelError:
		logger.SetLevel(logger.LevelError)
	}
}

// ApplyFormat maps a parsed Format onto the mritd logger.
func ApplyFormat(f Format) {
	switch f {
	case FormatConsole:
		logger.SetEncoder(logger.EncoderConsole)
	case FormatJSON:
		logger.SetEncoder(logger.EncoderJSON)
	}
}
