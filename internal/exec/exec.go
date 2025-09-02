package exec

import (
	e "os/exec"
)

// Run executes a system command with arguments and returns the result.
//
// This function provides a centralized command execution interface for the application,
// primarily used to run govulncheck and git commands. It captures both stdout and stderr
// in combined output and provides detailed error information for debugging.
//
// The function handles command execution errors gracefully and provides diagnostic
// information including exit codes and full command output. Error messages are
// printed to stdout for immediate visibility during scanning operations.
//
// Parameters:
//   - command: Full path to the executable to run
//   - args: Variable number of command line arguments
//
// Returns:
//   - exitCode: Process exit code (0 for success, non-zero for failure)
//   - output: Combined stdout and stderr from the command
//   - err: Non-nil if command execution failed (distinct from non-zero exit code)
//
// Security Note:
//
//	This function executes commands with user-provided arguments. Ensure all
//	inputs are properly validated before calling to prevent command injection.
//
// Example:
//
//	exitCode, output, err := Run("/usr/bin/git", "clone", "--depth", "1", "repo.git")
func Run(command string, args ...string) (exitCode int, output string, err error) {
	return RunAt("", command, args...)
}

func RunAt(dir, command string, args ...string) (exitCode int, output string, err error) {
	cmd := e.Command(command, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return cmd.ProcessState.ExitCode(), string(out), err
}
