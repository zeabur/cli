package util

import (
	"os/exec"
)

func PackZip() ([]byte, string, error) {
	format := "--format=zip"
	output := "--output=zeabur.zip"

	cmd := exec.Command("git", "archive", format, output, "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", err
	}
	return out, output, nil
}
