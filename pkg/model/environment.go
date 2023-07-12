package model

import "time"

type Environment struct {
	ID        string    `bson:"_id" json:"id" graphql:"_id"`
	Name      string    `bson:"name" json:"name" graphql:"name"`
	ProjectID string    `bson:"projectID" json:"projectID" graphql:"projectID"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`
}
