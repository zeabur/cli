package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Service is the simplest model of service, which is used in most queries.
type Service struct {
	ID   string `graphql:"_id"`
	Name string `graphql:"name"`
	//Template  ServiceTemplate    `graphql:"template"`
	//Project *Project `graphql:"project"`

	CreatedAt time.Time `graphql:"createdAt"`

	// MarketItemCode is the code of the item in the marketplace. Only used if Template is ServiceTemplateMarketplace.
	MarketItemCode *string `graphql:"marketItemCode"`

	CustomBuildCommand *string  `graphql:"customBuildCommand"`
	CustomStartCommand *string  `graphql:"customStartCommand"`
	OutputDir          *string  `graphql:"outputDir"`
	RootDirectory      string   `graphql:"rootDirectory"`
	Template           string   `graphql:"template"`
	WatchPaths         []string `graphql:"watchPaths"`
}

type Services []*Service

func (s Services) Header() []string {
	return []string{"ID", "Name", "Type", "CreatedAt"}
}

func (s Services) Rows() [][]string {
	rows := make([][]string, 0, len(s))
	headerLen := len(s.Header())

	for _, service := range s {
		row := make([]string, 0, headerLen)
		row = append(row, service.ID)
		row = append(row, service.Name)
		row = append(row, service.Template)
		row = append(row, service.CreatedAt.Format(time.RFC3339))

		rows = append(rows, row)
	}

	return rows
}

func (s *Service) Header() []string {
	return Services{s}.Header()
}

func (s *Service) Rows() [][]string {
	return Services{s}.Rows()
}

var (
	_ Tabler = (Services)(nil)
	_ Tabler = (*Service)(nil)
)

// ServiceDetail has more information related to specific environment.
type ServiceDetail struct {
	Service    `graphql:"... on Service"`
	GitTrigger *GitTrigger `graphql:"gitTrigger(environmentID: $environmentID)"`
	ConsoleURL string      `graphql:"consoleURL(environmentID: $environmentID)"`
	Status     string      `graphql:"status(environmentID: $environmentID)"`
	Domains    []Domain    `graphql:"domains(environmentID: $environmentID)"`
}

type ServiceDetails []*ServiceDetail

func (s ServiceDetails) Header() []string {
	return []string{"ID", "Name", "Status", "Domains", "Type", "GitTrigger", "CreatedAt"}
}

func (s ServiceDetails) Rows() [][]string {
	rows := make([][]string, 0, len(s))
	headerLen := len(s.Header())

	for _, service := range s {
		row := make([]string, 0, headerLen)
		row = append(row, service.ID)
		row = append(row, service.Name)
		row = append(row, service.Status)
		domains := make([]string, len(service.Domains))
		for i, domain := range service.Domains {
			domains[i] = domain.Domain
		}
		row = append(row, strings.Join(domains, ","))
		row = append(row, service.Template)
		gitTrigger := ""
		if service.GitTrigger != nil {
			gitTrigger = fmt.Sprintf("%s(%s)", service.GitTrigger.BranchName, service.GitTrigger.Provider)
		} else {
			gitTrigger = "None"
		}
		row = append(row, gitTrigger)
		row = append(row, service.CreatedAt.Format(time.RFC3339))

		rows = append(rows, row)
	}

	return rows
}

func (s *ServiceDetail) Header() []string {
	return ServiceDetails{s}.Header()
}

func (s *ServiceDetail) Rows() [][]string {
	return ServiceDetails{s}.Rows()
}

var (
	_ Tabler = (ServiceDetails)(nil)
	_ Tabler = (*ServiceDetail)(nil)
)

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
		Timestamp time.Time `json:"timestamp" graphql:"timestamp"`
		Value     float64   `json:"value" graphql:"value"`
	} `json:"metrics" graphql:"metrics(environmentID: $environmentID, metricType: $metricType, startTime: $startTime, endTime: $endTime, projectID: $projectID)"`
}

// MetricType is the type of metric.
type MetricType string

// valid metric types
const (
	MetricTypeCPU     MetricType = "CPU"
	MetricTypeMemory  MetricType = "MEMORY"
	MetricTypeNetwork MetricType = "NETWORK"
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
	case MetricTypeNetwork:
		return formatFloat64(v) + "MB"
	default:
		return formatFloat64(v)
	}
}

func formatFloat64(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

type ServiceTemplate string

const (
	ServiceTemplateGit         ServiceTemplate = "GIT"
	ServiceTemplateMarketplace ServiceTemplate = "MARKETPLACE"
)

type CreateServiceInput struct {
	ProjectID string          `json:"projectID" graphql:"projectID"`
	Name      string          `json:"name" graphql:"name"`
	Template  ServiceTemplate `graphql:"template"`
}

type MarketplaceItem struct {
	Code        string `graphql:"code"`
	Description string `graphql:"description"`
	IconURL     string `graphql:"iconURL"`
	Name        string `graphql:"name"`
	NetworkType string `graphql:"networkType"`
}

type PrebuiltItem struct {
	ID          string `graphql:"id"`
	Name        string `graphql:"name"`
	Description string `graphql:"description"`
}

type GitRepo struct {
	Name  string `graphql:"name"`
	Owner string `graphql:"owner"`
	URL   string `graphql:"url"`
	ID    int    `graphql:"id"`
}
