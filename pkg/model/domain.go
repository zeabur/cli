package model

import "time"

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
