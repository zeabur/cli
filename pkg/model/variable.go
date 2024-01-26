package model

type Variable struct {
	Key       string `json:"key" graphql:"key"`
	Value     string `json:"value" graphql:"value"`
	ServiceID string `json:"serviceID" graphql:"serviceID"`
}

type Variables []*Variable
