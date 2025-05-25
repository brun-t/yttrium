package yttrium

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCommandFail(t *testing.T) {
	cr := NewCommandRunner()

	// totally not existing command
	_, err := cr.Exec("ujnannnna")
	if err == nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("Function working correctly, first CommandRunner test passed")
}

func TestCommandSuccess(t *testing.T) {
	cr := NewCommandRunner()
	output, err := cr.Exec("echo", "a")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	trimmedOutput := strings.TrimSpace(string(output)) // Remove leading/trailing spaces, including newlines

	fmt.Println("Output:", trimmedOutput) // Print the trimmed output

	if trimmedOutput != "a" {
		t.Errorf("output:%s", trimmedOutput)
		t.FailNow()
	}

	fmt.Println("Second CommandRunner test passed")
}

func TestAsyncExecCrossPlatform(t *testing.T) {
	cr := NewCommandRunner()

	var cmd string

	switch cr.Shell {
	case "powershell", "powershell.exe":
		cmd = `echo line1; Start-Sleep -Seconds 1; echo line2; Start-Sleep -Seconds 1; echo line3`
	default:
		cmd = `echo line1; sleep 1; echo line2; sleep 1; echo line3`
	}

	//cr.SetShell("bash") uncomment in case of that windows shells stop working again

	outChan, err := cr.AsyncExec(cmd)
	if err != nil {
		t.Fatalf("AsyncExec failed: %v", err)
	}

	var output []string
	start := time.Now()
	for line := range outChan {
		t.Logf("OUT: %s", line)
		output = append(output, line)
	}
	duration := time.Since(start)

	if len(output) < 3 {
		t.Errorf("Expected at least 3 lines of output, got %d", len(output))
	}

	if duration < 2*time.Second {
		t.Errorf("Expected output to be streamed over time, finished too quickly: %v", duration)
	}

	expectedLines := []string{"line1", "line2", "line3"}
	for _, expected := range expectedLines {
		found := false
		for _, line := range output {
			if strings.Contains(line, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find '%s' in output", expected)
		}
	}
}
