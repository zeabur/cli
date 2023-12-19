package util

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func PackZip() ([]byte, error) {
	zipBytes, err := wrapNodeFunction(os.Getenv("PWD"), map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("wrap node function: %w", err)
	}

	return zipBytes, nil
}

func wrapNodeFunction(baseFolder string, envVars map[string]string) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	err := filepath.Walk(baseFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking to %s: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		// This will ensure only the content inside baseFolder is included at the root of the ZIP.
		relativePath, err := filepath.Rel(baseFolder, path)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		lstat, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("lstat: %w", err)
		}

		if lstat.Mode()&os.ModeSymlink == os.ModeSymlink {

			zipFile, err := w.Create(relativePath + ".link")
			if err != nil {
				return fmt.Errorf("creating zip file: %w", err)
			}

			target, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("read symlink: %w", err)
			}

			_, err = zipFile.Write([]byte(target))
			if err != nil {
				return fmt.Errorf("writing zip file: %w", err)
			}

		} else {

			zipFile, err := w.Create(relativePath)
			if err != nil {
				return fmt.Errorf("creating zip file: %w", err)
			}

			fileContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			_, err = zipFile.Write(fileContent)
			if err != nil {
				return fmt.Errorf("writing zip file: %w", err)
			}

		}

		zipFile, err := w.Create(".zeabur-env.json")
		if err != nil {
			return fmt.Errorf("creating zip file: %w", err)
		}

		envJsonStr, err := json.Marshal(envVars)
		if err != nil {
			return fmt.Errorf("marshaling env vars: %w", err)
		}

		_, err = zipFile.Write(envJsonStr)
		if err != nil {
			return fmt.Errorf("writing zip file: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking function directory: %w", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("closing zip writer: %w", err)
	}

	return buf.Bytes(), nil
}
