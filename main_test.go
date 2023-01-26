package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	stats, err := examine("testdata/test.txt.gz")
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
		t.Fail()
	}
}
