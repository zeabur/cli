package util

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func PackZip() ([]byte, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}

	bytesString, err := PackZipWithoutGitIgnoreFiles()
	if err != nil {
		return nil, "", err
	}

	return bytesString, currentDir, nil
}

func PackZipWithoutGitIgnoreFiles() ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Register a custom compressor for better compression
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	ignoreObject, err := gitignore.CompileIgnoreFile("./.gitignore")
	if err != nil {
		fmt.Println("Error compiling .gitignore file:", err)
	}

	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", path, err)
			return err
		}

		if path == "." {
			return nil
		}

		if strings.HasPrefix(path, ".git") {
			return nil
		}

		if ignoreObject != nil && ignoreObject.MatchesPath(path) {
			return nil
		}

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
