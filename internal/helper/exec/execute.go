package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var DryRun bool

func prepareCommand(command string, args ...string) (*exec.Cmd, bool) {
	if DryRun {
		fmt.Printf("Executing: %s %s\n", command, strings.Join(args, " "))
		return nil, false
	}
	return exec.Command(command, args...), true
}

func ExecuteCommand(command string, args ...string) error {
	cmd, proceed := prepareCommand(command, args...)
	if !proceed {
		return nil
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute '%s %s': %w", command, strings.Join(args, " "), err)
	}
	return nil
}

func ExecuteCommandWithStdin(command string, args ...string) error {
	cmd, proceed := prepareCommand(command, args...)
	if !proceed {
		return nil
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute with stdin '%s %s': %w", command, strings.Join(args, " "), err)
	}
	return nil
}

func ExecuteCommandWithStderr(command string, args ...string) (string, error) {
	cmd, proceed := prepareCommand(command, args...)
	if !proceed {
		return "", nil
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stderr.String(), err
}

func ExecuteCommandWithOutput(command string, args ...string) (string, error) {
	cmd, proceed := prepareCommand(command, args...)
	if !proceed {
		return "", nil
	}

	output, err := cmd.Output()
	return string(output), err
}

func RunCommand(command string, args ...string) ([]string, error) {
	output, err := ExecuteCommandWithOutput(command, args...)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var cleaned []string
	for _, line := range lines {
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned, nil
}
