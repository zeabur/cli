package model

import (
	"time"
)

type Project struct {
	ID          string `bson:"_id" json:"id" graphql:"_id"`
	Name        string `bson:"name" json:"name" graphql:"name"`
	Description string `bson:"description" json:"description" graphql:"description"`
	// todo: model environment
	//Environments []Environment     `bson:"environments" json:"environments" graphql:"environments"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`
	Owner         User      `bson:"owner" json:"owner" graphql:"owner"`
	Collaborators []User    `bson:"collaborators" json:"collaborators" graphql:"collaborators"`
	IconUrl       string    `bson:"iconUrl" json:"iconUrl" graphql:"iconURL"`
	// todo: model service
	//Services []Service `bson:"services" json:"services" graphql:"services"`
}

type ProjectConnection struct {
	PageInfo *PageInfo      `json:"pageInfo"`
	Edges    []*ProjectEdge `json:"edges"`
}

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
