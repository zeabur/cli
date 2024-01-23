package util

import (
	"fmt"
	"os"
	"os/exec"
)

func PackZip() ([]byte, string, error) {
	fileName, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Current working directory: ", fileName)

	if _, err := os.Stat(".git"); err == nil {
		bytes, err := PackGitRepo()
		if err != nil {
			return nil, "", err
		}
		return bytes, fileName, nil
	}

	bytes, err := PackZipFile()
	if err != nil {
		return nil, "", err
	}
	return bytes, fileName, nil
}

func PackGitRepo() ([]byte, error) {
	format := "--format=zip"
	output := "--output=zeabur.zip"

	cmd := exec.Command("git", "archive", format, output, "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func PackZipFile() ([]byte, error) {
	cmd := exec.Command("zip", "-r", "zeabur.zip", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}
