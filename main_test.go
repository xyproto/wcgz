package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	stats, err := Examine("testdata/test.txt.gz")
	if err != nil {
		t.Fail()
	}
	wcOutput, err := exec.Command("wc", "--bytes", "testdata/test.txt").CombinedOutput()
	if err != nil {
		t.Fail()
	}
	fields := strings.Fields(string(wcOutput))
	if len(fields) != 2 {
		t.Fail()
	}
	wcByteCount, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fail()
	}
	if uint64(wcByteCount) != stats.byteCounter {
		t.Errorf("expected a %d bytes, got %d", wcByteCount, stats.byteCounter)
	}
}

func TestLines(t *testing.T) {
	stats, err := Examine("testdata/test.txt.gz")
	if err != nil {
		t.Fail()
	}
	wcOutput, err := exec.Command("wc", "--lines", "testdata/test.txt").CombinedOutput()
	if err != nil {
		t.Fail()
	}
	fields := strings.Fields(string(wcOutput))
	if len(fields) != 2 {
		t.Fail()
	}
	wcLineCount, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fail()
	}
	if uint64(wcLineCount) != stats.lineCounter {
		t.Errorf("expected a %d lines, got %d", wcLineCount, stats.lineCounter)
	}
}

func TestChars(t *testing.T) {
	stats, err := Examine("testdata/test.txt.gz")
	if err != nil {
		t.Fail()
	}
	wcOutput, err := exec.Command("wc", "--chars", "testdata/test.txt").CombinedOutput()
	if err != nil {
		t.Fail()
	}
	fields := strings.Fields(string(wcOutput))
	if len(fields) != 2 {
		t.Fail()
	}
	wcCharCount, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fail()
	}
	if uint64(wcCharCount) != stats.runeCounter {
		t.Errorf("expected a %d chars/runes, got %d", wcCharCount, stats.runeCounter)
	}
}

func TestWords(t *testing.T) {
	stats, err := Examine("testdata/test.txt.gz")
	if err != nil {
		t.Fail()
	}
	wcOutput, err := exec.Command("wc", "--words", "testdata/test.txt").CombinedOutput()
	if err != nil {
		t.Fail()
	}
	fields := strings.Fields(string(wcOutput))
	if len(fields) != 2 {
		t.Fail()
	}
	wcWordCount, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fail()
	}
	if uint64(wcWordCount) != stats.wordCounter {
		t.Errorf("expected a %d words, got %d", wcWordCount, stats.wordCounter)
	}
}
