// Package model provides the models for the Zeabur API.
package model

import "time"

// Environment is the simplest model of environment, which is used in most queries.
type Environment struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`
	ID        string    `bson:"_id" json:"id" graphql:"_id"`
	Name      string    `bson:"name" json:"name" graphql:"name"`
	ProjectID string    `bson:"projectID" json:"projectID" graphql:"projectID"`
}

type Environments []*Environment

func (e Environments) Header() []string {
	return []string{"ID", "Name"}
}

func (e Environments) Rows() [][]string {
	rows := make([][]string, 0, len(e))
	for _, env := range e {
		rows = append(rows, []string{env.ID, env.Name})
	}
	return rows
}
