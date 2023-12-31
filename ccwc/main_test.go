package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestCCWC(t *testing.T) {
	t.Run("no flags", func(t *testing.T) {
		actual, err := runCCWC(t, "testdata/test.txt")
		expected := "7145 58164 342190 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("lines flag", func(t *testing.T) {
		actual, err := runCCWC(t, "-l", "testdata/test.txt")
		expected := "7145 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("words flag", func(t *testing.T) {
		actual, err := runCCWC(t, "-w", "testdata/test.txt")
		expected := "58164 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("characters flag", func(t *testing.T) {
		actual, err := runCCWC(t, "-m", "testdata/test.txt")
		expected := "339292 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("bytes flag", func(t *testing.T) {
		actual, err := runCCWC(t, "-c", "testdata/test.txt")
		expected := "342190 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("all flags", func(t *testing.T) {
		actual, err := runCCWC(t, "-l", "-w", "-m", "-c", "testdata/test.txt")
		expected := "7145 58164 339292 342190 testdata/test.txt\n"

		assertNil(t, err)
		assertEqual(t, actual, expected)
	})

	t.Run("no such file or directory", func(t *testing.T) {
		actual, err := runCCWC(t, "testdata/test1.txt")
		expectedErr := "exit status 1"
		expected := "ccwc: open testdata/test1.txt: no such file or directory\n"

		assertEqual(t, err.Error(), expectedErr)
		assertEqual(t, actual, expected)
	})

	t.Run("is a directory", func(t *testing.T) {
		actual, err := runCCWC(t, "testdata")
		expectedErr := "exit status 1"
		expected := "ccwc: testdata: Is a directory\n0 0 0 testdata\n"

		assertEqual(t, err.Error(), expectedErr)
		assertEqual(t, actual, expected)
	})
}

func runCCWC(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}

	cmd := exec.Command("bin/ccwc", args...)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()

	return buf.String(), err
}

func assertNil(t *testing.T, actual any) {
	t.Helper()
	if actual != nil {
		t.Errorf("got %v, want nil", actual)
	}
}

func assertEqual(t *testing.T, actual, expected any) {
	t.Helper()
	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}
