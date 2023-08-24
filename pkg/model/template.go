package model

import (
	"time"
)

type Template struct {
	Code          string               `graphql:"code"`
	CreatedAt     time.Time            `graphql:"createdAt"`
	DeploymentCnt int                  `graphql:"deploymentCnt"`
	Description   string               `graphql:"description"`
	Name          string               `graphql:"name"`
	PreviewURL    string               `graphql:"previewURL"`
	Readme        string               `graphql:"readme"`
	Tags          []string             `graphql:"tags"`
	Services      []*ServiceInTemplate `graphql:"services"`
}

type TemplateConnection struct {
	PageInfo *PageInfo       `json:"pageInfo"`
	Edges    []*TemplateEdge `json:"edges"`
}

type TemplateEdge struct {
	Node   *Template `json:"node"`
	Cursor string    `json:"cursor"`
}

type Templates []*Template

func (t Templates) Header() []string {
	return []string{"Code", "Name", "Description", "Created At"}
}

func (t Templates) Rows() [][]string {
	rows := make([][]string, 0, len(t))
	headerLen := len(t.Header())
	for _, template := range t {
		row := make([]string, 0, headerLen)
		row = append(row, template.Code)
		row = append(row, template.Name)
		row = append(row, template.Description)
		row = append(row, template.CreatedAt.Format(time.RFC3339))

		rows = append(rows, row)
	}
	return rows
}

func (t *Template) Header() []string {
	return Templates{t}.Header()
}

func (t *Template) Rows() [][]string {
	return Templates{t}.Rows()
}

var (
	_ Tabler = (Templates)(nil)
	_ Tabler = (*Template)(nil)
)

type RepoConfig struct {
	IsPublicRepo bool   `json:"isPublicRepo"`
	RepoName     string `json:"repoName"`
	ServiceName  string `json:"serviceName"`
}

type RepoConfigs []*RepoConfig

type ServiceInTemplate struct {
	BranchName         string                 `json:"branchName"`
	CustomBuildCommand string                 `json:"customBuildCommand"`
	CustomStartCommand string                 `json:"customStartCommand"`
	DomainKey          string                 `json:"domainKey"`
	GitNamespaceID     int                    `json:"gitNamespaceID"`
	GitRepoID          int                    `json:"gitRepoID"`
	MarketplaceItem    MarketplaceItem        `json:"marketplaceItem"`
	Name               string                 `json:"name"`
	OutputDir          string                 `json:"outputDir"`
	RootDirectory      string                 `json:"rootDirectory"`
	Template           ServiceTemplate        `json:"template"`
	Variables          map[string]interface{} `json:"variables"`
	WatchPaths         []string               `json:"watchPaths"`
}
