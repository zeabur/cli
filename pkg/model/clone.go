package model

import "time"

// CloneProjectResult is the result of the cloneProject mutation.
type CloneProjectResult struct {
	NewProjectID string `graphql:"newProjectId"`
}

// CloneProjectEvent is a single event emitted during project cloning.
type CloneProjectEvent struct {
	Type      string    `graphql:"type"`
	CreatedAt time.Time `graphql:"createdAt"`
	Message   string    `graphql:"message"`
}

// CloneProjectStatusResult is the result of the cloneProjectStatus query.
type CloneProjectStatusResult struct {
	NewProjectID *string             `graphql:"newProjectId"`
	Events       []CloneProjectEvent `graphql:"events"`
	Error        *string             `graphql:"error"`
}
