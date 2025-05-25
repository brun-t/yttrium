package yttrium

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
)

// A extremely tiny Wrapper for os/exec
type CommandRunner struct {
	OS    string
	Shell string
}

// The constructor for CommandRunner
func NewCommandRunner() *CommandRunner {
	return &CommandRunner{
		OS: runtime.GOOS,
	}
}

// It set a shell that is not the default one for the OS
func (cr *CommandRunner) SetShell(shell string) {
	cr.Shell = shell
}

func (cr *CommandRunner) getCommand(cmd string, flags ...string) ([]string, error) {
	if cr.Shell == "" {
		switch cr.OS {
		case "windows":
			cr.Shell = "powershell.exe"
		case "darwin":
			cr.Shell = "zsh"
		case "linux":
			cr.Shell = "bash"
		case "freebsd", "openbsd", "netbsd", "dragonfly", "solaris":
			cr.Shell = "/bin/sh"
		case "plan9":
			cr.Shell = "rc"
		default:
			return nil, fmt.Errorf("unsupported operating system: %s", cr.OS)
		}
	}

	var shellArgs []string

	switch cr.Shell {
	case "cmd", "cmd.exe":
		shellArgs = []string{"/C", cmd}
	case "powershell", "powershell.exe":
		shellArgs = []string{"-Command", cmd}
	default:
		// Assume POSIX shell
		shellArgs = []string{"-c", cmd}
	}

	shellArgs = append(shellArgs, flags...)

	return shellArgs, nil
}

// It execs a command it returns both stdout, and stderr in same byte and the error
// it uses the "default" shell for the OS
func (cr *CommandRunner) Exec(cmd string, flags ...string) ([]byte, error) {
	shellArgs, err := cr.getCommand(cmd, flags...)

	if err != nil {
		return nil, err
	}

	command := exec.Command(cr.Shell, shellArgs...)

	output, err := command.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error executing command (%s %v): %w, output: %s", cr.Shell, shellArgs, err, output)
	}
	return output, nil
}

func (cr *CommandRunner) AsyncExec(cmd string, flags ...string) (<-chan string, error) {
	shellArgs, err := cr.getCommand(cmd, flags...)

	if err != nil {
		return nil, err
	}

	pipeChan := make(chan string)

	command := exec.Command(cr.Shell, shellArgs...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	command.Stderr = command.Stdout

	if err := command.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer close(pipeChan)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			pipeChan <- scanner.Text()
		}
		command.Wait()
	}()

	return pipeChan, nil
}
