package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var binaryExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".bmp": true,
	".ico": true, ".svg": true, ".webp": true, ".avif": true,
	".woff": true, ".woff2": true, ".ttf": true, ".otf": true, ".eot": true,
	".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".7z": true, ".rar": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".exe": true, ".dll": true, ".so": true, ".dylib": true, ".wasm": true,
	".mp3": true, ".mp4": true, ".wav": true, ".ogg": true, ".webm": true,
	".bin": true, ".dat": true, ".db": true, ".sqlite": true,
}

// ListUploadFiles lists files in an uploaded project.
func (c *client) ListUploadFiles(ctx context.Context, uploadID string, path *string) ([]string, error) {
	var query struct {
		Files []string `graphql:"files(uploadID: $uploadID, path: $path)"`
	}

	err := c.Query(ctx, &query, V{
		"uploadID": ObjectID(uploadID),
		"path":     path,
	})
	if err != nil {
		return nil, err
	}

	return query.Files, nil
}

// ReadUploadFile reads the content of a file in an uploaded project.
func (c *client) ReadUploadFile(ctx context.Context, uploadID string, path string) (string, error) {
	var query struct {
		FileContent string `graphql:"fileContent(uploadID: $uploadID, path: $path)"`
	}

	err := c.Query(ctx, &query, V{
		"uploadID": ObjectID(uploadID),
		"path":     path,
	})
	if err != nil {
		return "", err
	}

	return query.FileContent, nil
}

// PullUploadFiles downloads all files from an upload to a local directory.
func (c *client) PullUploadFiles(ctx context.Context, uploadID string, targetDir string) (int, error) {
	return c.pullDir(ctx, uploadID, "", targetDir)
}

func (c *client) pullDir(ctx context.Context, uploadID string, remotePath string, localDir string) (int, error) {
	baseDir := filepath.Clean(localDir)

	var pathPtr *string
	if remotePath != "" {
		pathPtr = &remotePath
	}

	entries, err := c.ListUploadFiles(ctx, uploadID, pathPtr)
	if err != nil {
		return 0, fmt.Errorf("list %q: %w", remotePath, err)
	}

	count := 0
	for _, entry := range entries {
		fullRemote := remotePath + entry
		fullLocal := filepath.Join(baseDir, fullRemote)

		if !strings.HasPrefix(filepath.Clean(fullLocal), baseDir) {
			return count, fmt.Errorf("path traversal detected: %q", fullRemote)
		}

		if strings.HasSuffix(entry, "/") {
			if err := os.MkdirAll(fullLocal, 0o755); err != nil {
				return count, fmt.Errorf("mkdir %q: %w", fullLocal, err)
			}
			n, err := c.pullDir(ctx, uploadID, fullRemote, localDir)
			count += n
			if err != nil {
				return count, err
			}
		} else {
			ext := strings.ToLower(filepath.Ext(entry))
			if binaryExtensions[ext] {
				continue
			}
			content, err := c.ReadUploadFile(ctx, uploadID, fullRemote)
			if err != nil {
				return count, fmt.Errorf("read %q: %w", fullRemote, err)
			}
			dir := filepath.Dir(fullLocal)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return count, fmt.Errorf("mkdir %q: %w", dir, err)
			}
			if err := os.WriteFile(fullLocal, []byte(content), 0o644); err != nil {
				return count, fmt.Errorf("write %q: %w", fullLocal, err)
			}
			count++
		}
	}

	return count, nil
}
