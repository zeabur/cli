package model

import (
	"strconv"
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

// ServiceDetailConnection is a connection to a list of items.
type ServiceDetailConnection struct {
	PageInfo *PageInfo            `json:"pageInfo" graphql:"pageInfo"`
	Edges    []*ServiceDetailEdge `json:"edges" graphql:"edges"`
}

// ServiceDetailEdge is an edge in a connection.
type ServiceDetailEdge struct {
	Node   *ServiceDetail `json:"node" graphql:"node"`
	Cursor string         `json:"cursor" graphql:"cursor"`
}

// TempTCPPort is a temporary TCP port.
type TempTCPPort struct {
	ServiceID     string `json:"serviceID" graphql:"serviceID"`
	EnvironmentID string `json:"environmentID" graphql:"environmentID"`
	TargetPort    int    `json:"targetPort" graphql:"targetPort"`
	NodePort      int    `json:"nodePort" graphql:"nodePort"`
	RemainSeconds int    `json:"remainSeconds" graphql:"remainSeconds"`
}

// GitTrigger represents a git trigger.
type GitTrigger struct {
	BranchName string `json:"branchName" graphql:"branchName"`
	Provider   string `json:"provider" graphql:"provider"`
	RepoID     int    `json:"repoID" graphql:"repoID"`
}

// ServiceMetric is a metric of a service.
type ServiceMetric struct {
	Metrics []struct {
		Value     float64   `json:"value" graphql:"value"`
		Timestamp time.Time `json:"timestamp" graphql:"timestamp"`
	} `json:"metrics" graphql:"metrics(environmentID: $environmentID, metricType: $metricType, startTime: $startTime, endTime: $endTime)"`
}

// MetricType is the type of metric.
type MetricType string

// valid metric types
const (
	MetricTypeCPU     MetricType = "CPU"
	MetricTypeMemory  MetricType = "MEMORY"
	MetricTypeNetwork MetricType = "NETWORK"
	MetricTypeDisk    MetricType = "DISK"
)

func (m MetricType) GetGraphQLType() string {
	return "MetricType"
}

func (m MetricType) WithMeasureUnit(v float64) string {
	switch m {
	case MetricTypeCPU:
		return formatFloat64(v*100) + "%"
	case MetricTypeMemory:
		return formatFloat64(v) + "MB"
		//case MetricTypeNetwork:
		//	return formatFloat64(v) + "MB"
		//case MetricTypeDisk:
		//	return formatFloat64(v) + "MB"
	default:
		return formatFloat64(v)
	}
}

func formatFloat64(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
