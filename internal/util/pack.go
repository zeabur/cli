package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	zipBytes, err := os.ReadFile("zeabur.zip")
	if err != nil {
		return nil, err
	}

	return zipBytes, nil
}

func PackZipFile() ([]byte, error) {
	buf := new(bytes.Buffer)
	srcDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	zipWriter := zip.NewWriter(buf)

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(".", path)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
