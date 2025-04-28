package reviewbot

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetGitDiff returns the differences between the current branch and the specified target branch.
//
// The function runs the `git diff` command to compare the current branch with the provided
// target branch and returns the output as a string. If an error occurs during the command
// execution, the function returns an empty string and the error encountered.
//
// Parameters:
//   - targetBranch (string): The name of the branch to compare with the current branch.
//
// Returns:
//   - string: The output of the `git diff` command, which shows the differences between
//     the current branch and the target branch.
//   - error: Any error encountered during the execution of the command. If no error occurs,
//     this will be nil.
//
// Example usage:
//
//	diff, err := GetGitDiff("main")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Git Diff:", diff)
//	}
//
// Note:
//   - The function captures both stdout and stderr of the `git diff` command.
func GetGitDiff(targetBranch string) (string, error) {
	cmd := exec.Command("git", "diff", targetBranch)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error executing git diff: %w", err)
	}

	return output.String(), nil
}
