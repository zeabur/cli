// Package model provides the models for the Zeabur API.
package model

import "time"

// Environment is the simplest model of environment, which is used in most queries.
type Environment struct {
	CreatedAt time.Time `graphql:"createdAt"`
	ID        string    `graphql:"_id"`
	Name      string    `graphql:"name"`
	ProjectID string    `graphql:"projectID"`
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

func (e *Environment) Header() []string {
	return Environments{e}.Header()
}

func (e *Environment) Rows() [][]string {
	return Environments{e}.Rows()
}

var (
	_ Tabler = (*Environments)(nil)
	_ Tabler = (*Environment)(nil)
)
