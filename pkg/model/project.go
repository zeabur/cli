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
	ID          string `bson:"_id" json:"id" graphql:"_id"`
	Name        string `bson:"name" json:"name" graphql:"name"`
	Description string `bson:"description" json:"description" graphql:"description"`
	//Environments []Environment     `bson:"environments" json:"environments" graphql:"environments"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`
	//Owner         User      `bson:"owner" json:"owner" graphql:"owner"`
	//Collaborators []User    `bson:"collaborators" json:"collaborators" graphql:"collaborators"`
	IconURL string `bson:"iconUrl" json:"iconUrl" graphql:"iconURL"`
	//Services []Service `bson:"services" json:"services" graphql:"services"`
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
