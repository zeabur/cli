package model

import "time"

type Domain struct {
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt" graphql:"createdAt"`
	ID            string    `json:"id" bson:"_id" graphql:"_id"`
	Domain        string    `json:"domain" bson:"domain" graphql:"domain"`
	EnvironmentID string    `json:"environmentID" bson:"environmentID" graphql:"environmentID"`
	ProjectID     string    `json:"projectID" bson:"projectID" graphql:"projectID"`
	RedirectTo    string    `json:"redirectTo" bson:"redirectTo" graphql:"redirectTo"`
	ServiceID     string    `json:"serviceID" bson:"serviceID" graphql:"serviceID"`
	Status        string    `json:"status" bson:"status" graphql:"status"`
	IsGenerated   bool      `json:"isGenerated" bson:"isGenerated" graphql:"isGenerated"`
}
