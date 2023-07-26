package model

import (
	"time"
)

// Service is the simplest model of service, which is used in most queries.
type Service struct {
	ID   string `bson:"_id" json:"id" graphql:"_id"`
	Name string `bson:"name" json:"name" graphql:"name"`
	//Template  ServiceTemplate    `bson:"template" json:"template" graphql:"template"`
	//Project *Project `bson:"project" json:"project" graphql:"project"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// MarketItemCode is the code of the item in the marketplace. Only used if Template is ServiceTemplateMarketplace.
	MarketItemCode *string `bson:"marketItemCode" json:"marketItemCode" graphql:"marketItemCode"`

	CustomBuildCommand *string  `bson:"customBuildCommand" json:"customBuildCommand" graphql:"customBuildCommand"`
	CustomStartCommand *string  `bson:"customStartCommand" json:"customStartCommand" graphql:"customStartCommand"`
	OutputDir          *string  `bson:"outputDir" json:"outputDir" graphql:"outputDir"`
	RootDirectory      string   `bson:"rootDirectory" json:"rootDirectory" graphql:"rootDirectory"`
	WatchPaths         []string `bson:"watchPaths" json:"watchPaths" graphql:"watchPaths"`
	Template           string   `bson:"template" json:"template" graphql:"template"`
}

// ServiceDetail has more information related to specific environment.
type ServiceDetail struct {
	Service    `bson:",inline" graphql:"... on Service"`
	ConsoleURL string      `bson:"consoleURL" json:"consoleURL" graphql:"consoleURL(environmentID: $environmentID)"`
	Domains    []Domain    `bson:"domains" json:"domains" graphql:"domains(environmentID: $environmentID)"`
	GitTrigger *GitTrigger `bson:"gitTrigger" json:"gitTrigger" graphql:"gitTrigger(environmentID: $environmentID)"`
	Status     string      `bson:"status" json:"status" graphql:"status(environmentID: $environmentID)"`
}

// ServiceConnection is a connection to a list of items.
type ServiceConnection struct {
	PageInfo *PageInfo      `json:"pageInfo" graphql:"pageInfo"`
	Edges    []*ServiceEdge `json:"edges" graphql:"edges"`
}

// ServiceEdge is an edge in a connection.
type ServiceEdge struct {
	Node   *Service `json:"node" graphql:"node"`
	Cursor string   `json:"cursor" graphql:"cursor"`
}

type ServiceDetailConnection struct {
	PageInfo *PageInfo            `json:"pageInfo" graphql:"pageInfo"`
	Edges    []*ServiceDetailEdge `json:"edges" graphql:"edges"`
}

type ServiceDetailEdge struct {
	Node   *ServiceDetail `json:"node" graphql:"node"`
	Cursor string         `json:"cursor" graphql:"cursor"`
}

type TempTCPPort struct {
	ServiceID     string `json:"serviceID" graphql:"serviceID"`
	EnvironmentID string `json:"environmentID" graphql:"environmentID"`
	TargetPort    int    `json:"targetPort" graphql:"targetPort"`
	NodePort      int    `json:"nodePort" graphql:"nodePort"`
	RemainSeconds int    `json:"remainSeconds" graphql:"remainSeconds"`
}

type GitTrigger struct {
	BranchName string `json:"branchName" graphql:"branchName"`
	Provider   string `json:"provider" graphql:"provider"`
	RepoID     int    `json:"repoID" graphql:"repoID"`
}
