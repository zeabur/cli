package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/klauspost/compress/flate"
)

func PackZip() ([]byte, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Current working directory: ", currentDir)

	if _, err := os.Stat(".git"); err == nil {
		bytesString, err := PackGitRepo()
		if err != nil {
			fmt.Println("Error packing git repository to zip file!!!")
			return nil, "", err
		}
		return bytesString, currentDir, nil
	}

	bytesString, err := PackZipFile()
	if err != nil {
		fmt.Println("Error packing current directory to zip file!!!")
		return nil, "", err
	}

	return bytesString, currentDir, nil
}

func PackGitRepo() ([]byte, error) {
	format := "--format=zip"
	output := "--output=zeabur.zip"

	cmd := exec.Command("git", "archive", format, output, "HEAD")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = os.Remove("zeabur.zip")
		if err != nil {
			fmt.Println(err)
		}
	}()

	zipBytes, err := os.ReadFile("zeabur.zip")
	if err != nil {
		return nil, err
	}

	return zipBytes, nil
}

func PackZipFile() ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Register a custom compressor for better compression.
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the . directory.
		if path != "." {
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name, err = filepath.Rel(".", path)
			if err != nil {
				return err
			}

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
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Close the zip writer.
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
