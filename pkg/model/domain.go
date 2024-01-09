package model

import "time"

type Domain struct {
	CreatedAt   time.Time `json:"createdAt" graphql:"createdAt"`
	ID          string    `json:"id" graphql:"_id"`
	Domain      string    `json:"domain" graphql:"domain"`
	ProjectID   string    `json:"projectID" graphql:"projectID"`
	RedirectTo  string    `json:"redirectTo" graphql:"redirectTo"`
	Status      string    `json:"status" graphql:"status"`
	IsGenerated bool      `json:"isGenerated" graphql:"isGenerated"`
}
