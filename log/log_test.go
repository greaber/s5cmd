package log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureOutputs(t *testing.T, f func()) (stdout, stderr string) {
	t.Helper()

	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	f()
	if err := stdoutWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := stderrWriter.Close(); err != nil {
		t.Fatal(err)
	}

	var stdoutBuffer, stderrBuffer bytes.Buffer
	if _, err := stdoutBuffer.ReadFrom(stdoutReader); err != nil {
		t.Fatal(err)
	}
	if _, err := stderrBuffer.ReadFrom(stderrReader); err != nil {
		t.Fatal(err)
	}
	return stdoutBuffer.String(), stderrBuffer.String()
}

func resetOutputChannel() {
	outputCh = make(chan output, 10000)
}

func TestDebugWritesToStderr(t *testing.T) {
	resetOutputChannel()
	Init("debug", false)

	stdout, stderr := captureOutputs(t, func() {
		Debug(DebugMessage{Err: "hello debug"})
		Close()
	})
	if stdout != "" {
		t.Fatalf("debug wrote to stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "hello debug") {
		t.Fatalf("debug missing from stderr: %q", stderr)
	}
}

func TestTraceWritesToStderr(t *testing.T) {
	resetOutputChannel()
	Init("trace", false)

	stdout, stderr := captureOutputs(t, func() {
		Trace(TraceMessage{Message: "hello trace"})
		Close()
	})
	if stdout != "" {
		t.Fatalf("trace wrote to stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "hello trace") {
		t.Fatalf("trace missing from stderr: %q", stderr)
	}
}

func TestInfoWritesToStdout(t *testing.T) {
	resetOutputChannel()
	Init("info", false)

	stdout, stderr := captureOutputs(t, func() {
		Info(ErrorMessage{Err: "hello info"})
		Close()
	})
	if !strings.Contains(stdout, "hello info") {
		t.Fatalf("info missing from stdout: %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("info wrote to stderr: %q", stderr)
	}
}
