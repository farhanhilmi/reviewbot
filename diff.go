package reviewbot

import (
	"bytes"
	"fmt"
	"os/exec"
)

func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "main...HEAD", "--unified=0")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff error: %w", err)
	}
	return out.String(), nil
}
