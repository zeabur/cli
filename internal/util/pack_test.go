package util_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeabur/cli/internal/util"
)

func TestPackZipWithZeaburIgnore(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "zeabur-pack-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create test files and directories
	testFiles := map[string]string{
		"main.go":                    "package main",
		"README.md":                  "# Test Project",
		".agent/skills/test.txt":     "should be ignored",
		".agents/config.json":        "should be ignored",
		"src/app.go":                 "package src",
		".gitignore":                 "*.log\n",
		"test.log":                   "log content",
		".git/config":                "[core]\n",
	}

	for path, content := range testFiles {
		dir := filepath.Dir(path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create dir %s: %v", dir, err)
			}
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Create .zeaburignore file
	zeaburignore := `.agent/
.agents/
.cursor/
`
	if err := os.WriteFile(".zeaburignore", []byte(zeaburignore), 0644); err != nil {
		t.Fatalf("Failed to create .zeaburignore: %v", err)
	}

	// Pack the zip
	zipBytes, err := util.PackZipWithoutGitIgnoreFiles()
	if err != nil {
		t.Fatalf("PackZipWithoutGitIgnoreFiles failed: %v", err)
	}

	// Read the zip and check contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("Failed to read zip: %v", err)
	}

	filesInZip := make(map[string]bool)
	for _, file := range zipReader.File {
		filesInZip[file.Name] = true
		t.Logf("File in zip: %s", file.Name)
	}

	// Check that expected files are included
	expectedFiles := []string{"main.go", "README.md", "src/app.go", ".zeaburignore", ".gitignore"}
	for _, file := range expectedFiles {
		if !filesInZip[file] {
			t.Errorf("Expected file %s not found in zip", file)
		}
	}

	// Check that .zeaburignore patterns are excluded
	excludedByZeaburIgnore := []string{".agent/skills/test.txt", ".agents/config.json", ".agent/", ".agents/"}
	for _, file := range excludedByZeaburIgnore {
		if filesInZip[file] {
			t.Errorf("File %s should be excluded by .zeaburignore but found in zip", file)
		}
	}

	// When .zeaburignore exists, .gitignore patterns are NOT applied
	// So test.log should be included in the zip
	if !filesInZip["test.log"] {
		t.Errorf("Expected test.log to be included in zip when .zeaburignore takes precedence over .gitignore")
	}

	// Check that .git directory is excluded
	for path := range filesInZip {
		if strings.HasPrefix(path, ".git/") || strings.HasPrefix(path, ".git\\") {
			t.Errorf("File %s should be excluded (.git directory) but found in zip", path)
		}
	}
}

func TestPackZipWithoutZeaburIgnore(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "zeabur-pack-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		"main.go":   "package main",
		"README.md": "# Test Project",
		"test.log":  "log content",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Create .gitignore file (no .zeaburignore)
	gitignore := "*.log\n"
	if err := os.WriteFile(".gitignore", []byte(gitignore), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Pack the zip
	zipBytes, err := util.PackZipWithoutGitIgnoreFiles()
	if err != nil {
		t.Fatalf("PackZipWithoutGitIgnoreFiles failed: %v", err)
	}

	// Read the zip and check contents
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("Failed to read zip: %v", err)
	}

	filesInZip := make(map[string]bool)
	for _, file := range zipReader.File {
		filesInZip[file.Name] = true
		t.Logf("File in zip: %s", file.Name)
	}

	// Check that expected files are included
	expectedFiles := []string{"main.go", "README.md", ".gitignore"}
	for _, file := range expectedFiles {
		if !filesInZip[file] {
			t.Errorf("Expected file %s not found in zip", file)
		}
	}

	// Check that .gitignore patterns are respected (should fallback to .gitignore)
	if filesInZip["test.log"] {
		t.Errorf("File test.log should be excluded by .gitignore but found in zip")
	}
}
