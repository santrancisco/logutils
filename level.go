// Package logutils augments the standard log package with levels.
package logutils

import (
	"bytes"
	"io"
	"sync"

	color "github.com/fatih/color"
)

type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger that
// will filter out log messages that aren't at least a certain level.
//
// Once the filter is in use somewhere, it is not safe to modify
// the structure.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"DEBUG", "WARN", "ERROR"}.
	Levels []LogLevel

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	badLevels map[LogLevel]struct{}
	once      sync.Once
}

var colorlist = []color.Attribute{color.FgRed, color.FgYellow, color.FgGreen, color.FgBlue}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) (color.Attribute, bool) {
	f.once.Do(f.init)

	// Check for a log level
	var level LogLevel
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = LogLevel(line[x+1 : x+y])
		}
	}

	_, isbad := f.badLevels[level]
	// if it is in the list of bad levels, we return nothing
	if isbad {
		return color.FgBlack, false
	}
	cl := color.Reset
	// This is where we determine the color of the log.
	// Assuming the highest severity log is always RED, we pre-determine the color for the last 4
	// highest severity from Blue, Green, Yellow and Red. Anything fall outside of this will have
	// default terminal color.
	for i, v := range f.Levels {
		if v == level {
			gap := len(f.Levels) - i - 1
			if gap > len(colorlist) {
				break
			}
			cl = colorlist[gap]
			break
		}
	}
	return cl, true
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line. We use that as a slight optimization within
	// this method, assuming we're dealing with a single, complete line
	// of log data.

	cl, ok := f.Check(p)
	if ok == false {
		return len(p), nil
	}
	return color.New(cl).Fprint(f.Writer, string(p))
	// return f.Writer.Write(p)
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	badLevels := make(map[LogLevel]struct{})
	for _, level := range f.Levels {
		if level == f.MinLevel {
			break
		}
		badLevels[level] = struct{}{}
	}
	f.badLevels = badLevels
}
