package model

// Tabler represents the model that can be displayed as a table
type Tabler interface {
	Header() []string
	Rows() [][]string
}
