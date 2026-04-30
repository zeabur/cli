package api

import (
	"context"
)

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
