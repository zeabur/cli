package model

import (
	"github.com/zeabur/cli/pkg/util"
	"time"
)

type Domain struct {
	ID            string    `json:"id" graphql:"_id"`
	Domain        string    `json:"domain" graphql:"domain"`
	ServiceID     string    `json:"serviceID" graphql:"serviceID"`
	EnvironmentID string    `json:"environmentID" graphql:"environmentID"`
	ProjectID     string    `json:"projectID" graphql:"projectID"`
	PortName      string    `json:"portName" graphql:"portName"`
	RedirectTo    string    `json:"redirectTo" graphql:"redirectTo"`
	Status        string    `json:"status" graphql:"status"`
	IsGenerated   bool      `json:"isGenerated" graphql:"isGenerated"`
	CreatedAt     time.Time `json:"createdAt" graphql:"createdAt"`
}

type Domains []*Domain

func (d Domains) Header() []string {
	return []string{"Domain", "Status", "CreatedAt"}
}

func (d Domains) Rows() [][]string {
	rows := make([][]string, 0, len(d))
	headerLen := len(d.Header())
	for _, domain := range d {
		row := make([]string, 0, headerLen)
		if domain.RedirectTo != "" {
			row = append(row, domain.Domain+"(->"+domain.RedirectTo+")")
		} else {
			row = append(row, domain.Domain)
		}
		row = append(row, domain.Status)
		row = append(row, util.ConvertTimeAgoString(domain.CreatedAt))

		rows = append(rows, row)
	}
	return rows
}
