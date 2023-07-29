package model

import (
	"time"
)

// Note: it's not recommended to embed other models in a model,
// because graphql will query them recursively.
// If you want to take advantage of graphql's nested query, you can define a new type,
// such as ProjectWithEnvironments, or xxxProjectResponse.

// Therefor, if the name of model doesn't has any prefix or suffix,
// we regard it as a basic model, to which we can add some basic methods and reuse them.

// if the name of model has a prefix or suffix, we only use it in the specific query.

// Project is the simplest model of project, which is used in most queries.
type Project struct {
	ID          string `graphql:"_id"`
	Name        string `graphql:"name"`
	Description string `graphql:"description"`
	//Environments []Environment     `graphql:"environments"`
	CreatedAt time.Time `graphql:"createdAt"`
	//Owner         User      `graphql:"owner"`
	//Collaborators []User    `graphql:"collaborators"`
	IconURL string `graphql:"iconURL"`
	//Services []Service `graphql:"services"`
}

// ProjectConnection is a connection to a list of items.
type ProjectConnection struct {
	PageInfo *PageInfo      `json:"pageInfo"`
	Edges    []*ProjectEdge `json:"edges"`
}

// ProjectEdge is an edge in a connection.
type ProjectEdge struct {
	Node   *Project `json:"node"`
	Cursor string   `json:"cursor"`
}

// ProjectUsage is a summary of the usage of a project in a given time period.
type ProjectUsage struct {
	// the project this usage is for
	Project *Project `json:"project"`
	// the beginning of the time period
	From time.Time `json:"from"`
	// the end of the time period
	To time.Time `json:"to"`
	// the total number of cpu used (in vCPU-minutes)
	CPU float64 `json:"cpu"`
	// the total number of memory used (in MB-minutes)
	Memory int `json:"memory"`
}

type Projects []*Project

func (p Projects) Header() []string {
	return []string{"ID", "Name", "Description", "Created At"}
}

func (p Projects) Rows() [][]string {
	rows := make([][]string, 0, len(p))
	headerLen := len(p.Header())
	for _, project := range p {
		row := make([]string, 0, headerLen)
		row = append(row, project.ID)
		row = append(row, project.Name)
		row = append(row, project.Description)
		row = append(row, project.CreatedAt.Format(time.RFC3339))

		rows = append(rows, row)
	}
	return rows
}

func (p *Project) Header() []string {
	return Projects{p}.Header()
}

func (p *Project) Rows() [][]string {
	return Projects{p}.Rows()
}

var (
	_ Tabler = (Projects)(nil)
	_ Tabler = (*Project)(nil)
)
