package logutils

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestLevelFilter_impl(t *testing.T) {
	var _ io.Writer = new(LevelFilter)
}

func TestLevelFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   buf,
	}

	logger := log.New(filter, "", 0)
	logger.Print("[WARN] foo")
	logger.Println("[ERROR] bar")
	logger.Println("[DEBUG] baz")
	logger.Println("[WARN] buzz")

	result := buf.String()
	expected := "\x1b[33m[WARN] foo\n\x1b[0m\x1b[31m[ERROR] bar\n\x1b[0m\x1b[33m[WARN] buzz\n\x1b[0m"
	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
}

func TestLevelFilterCheck(t *testing.T) {
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   nil,
	}

	testCases := []struct {
		line  string
		check bool
	}{
		{"[WARN] foo\n", true},
		{"[ERROR] bar\n", true},
		{"[DEBUG] baz\n", false},
		{"[WARN] buzz\n", true},
	}

	for _, testCase := range testCases {
		_, result := filter.Check([]byte(testCase.line))
		if result != testCase.check {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}

func TestLevelFilter_SetMinLevel(t *testing.T) {
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "ERROR",
		Writer:   nil,
	}

	testCases := []struct {
		line        string
		checkBefore bool
		checkAfter  bool
	}{
		{"[WARN] foo\n", false, true},
		{"[ERROR] bar\n", true, true},
		{"[DEBUG] baz\n", false, false},
		{"[WARN] buzz\n", false, true},
	}

	for _, testCase := range testCases {
		_, result := filter.Check([]byte(testCase.line))
		if result != testCase.checkBefore {
			t.Errorf("Fail: %s", testCase.line)
		}
	}

	// Update the minimum level to WARN
	filter.SetMinLevel("WARN")

	for _, testCase := range testCases {
		_, result := filter.Check([]byte(testCase.line))
		if result != testCase.checkAfter {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}
