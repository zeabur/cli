package model

import (
	"github.com/zeabur/cli/pkg/util"
	"time"
)

type Deployment struct {
	ID string `graphql:"_id"`

	ProjectID     string `json:"projectID" graphql:"projectID"`
	ServiceID     string `json:"serviceID" graphql:"serviceID"`
	EnvironmentID string `json:"environmentID" graphql:"environmentID"`
	GitProvider   string `json:"gitProvider" graphql:"gitProvider"`
	RepoOwner     string `json:"repoOwner" graphql:"repoOwner"`
	RepoName      string `json:"repoName" graphql:"repoName"`
	Ref           string `json:"ref" graphql:"ref"`
	CommitSHA     string `json:"commitSHA" graphql:"commitSHA"`
	CommitMessage string `json:"commitMessage" graphql:"commitMessage"`

	PlanType string `json:"planType" graphql:"planType"`
	//PlanMeta string `json:"planMeta" graphql:"planMeta" scala:"true"`

	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	ScheduledAt time.Time `json:"scheduledAt" graphql:"scheduledAt"`
	StartedAt   time.Time `json:"startedAt" graphql:"startedAt"`
	FinishedAt  time.Time `json:"finishedAt" graphql:"finishedAt"`

	Status string `json:"status" graphql:"status"`
}

type Deployments []*Deployment

func (d Deployments) Header() []string {
	return []string{"ID", "RepoName", "Status", "Ref", "CommitMessage", "PlanType", "CreatedAt", "CommitSHA"}
}

func (d Deployments) Rows() [][]string {
	headerLen := len(d.Header())
	rows := make([][]string, 0, len(d))
	for _, deployment := range d {
		row := make([]string, 0, headerLen)
		row = append(row, deployment.ID)
		row = append(row, deployment.RepoName)
		row = append(row, deployment.Status)
		row = append(row, deployment.Ref)
		row = append(row, truncateString(deployment.CommitMessage, 20))
		row = append(row, deployment.PlanType)
		row = append(row, util.ConvertTimeAgoString(deployment.CreatedAt))
		row = append(row, truncateString(deployment.CommitSHA, 8))

		rows = append(rows, row)
	}

	return rows
}

func (d *Deployment) Header() []string {
	return Deployments{d}.Header()
}

func (d *Deployment) Rows() [][]string {
	return Deployments{d}.Rows()
}

func truncateString(s string, maxLen int) string {
	// convert string to rune slice
	rs := []rune(s)
	if len(rs) > maxLen {
		return string(rs[:maxLen]) + "..."
	}

	return s
}

var (
	_ Tabler = (Deployments)(nil)
	_ Tabler = (*Deployment)(nil)
)
