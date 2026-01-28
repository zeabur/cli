package util

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"errors"
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

	// .zeaburignore has higher priority than .gitignore
	// Try to load .zeaburignore first, fallback to .gitignore if not exists
	var ignoreObject *gitignore.GitIgnore
	var err error

	ignoreObject, err = gitignore.CompileIgnoreFile("./.zeaburignore")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// .zeaburignore not found, try .gitignore
			ignoreObject, err = gitignore.CompileIgnoreFile("./.gitignore")
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					fmt.Println("Error compiling .gitignore file:", err)
				}
				ignoreObject = nil
			}
		} else {
			fmt.Println("Error compiling .zeaburignore file:", err)
			ignoreObject = nil
		}
	}

	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			// Skip files/directories that cannot be accessed (e.g., symlinks to non-existent targets)
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if path == "." {
			return nil
		}

		// Skip .git directory but not .gitignore or other .git* files
		if path == ".git" || strings.HasPrefix(path, ".git"+string(filepath.Separator)) {
			if info != nil && info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check ignore patterns before processing
		if ignoreObject != nil {
			// For directories, we need to check with trailing slash for proper gitignore matching
			checkPath := path
			if info.IsDir() {
				checkPath = path + "/"
			}

			if ignoreObject.MatchesPath(checkPath) {
				// Skip ignored files/directories entirely
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip symlinks to avoid "is a directory" errors
		if info.Mode()&os.ModeSymlink != 0 {
			if info.IsDir() {
				return filepath.SkipDir
			}
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

		// For Windows, we should replace the backslashes with forward slashes
		header.Name = filepath.ToSlash(header.Name)

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
