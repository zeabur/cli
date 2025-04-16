package model

type GenericRegion interface {
	GetID() string
	IsAvailable() bool
	String() string
}

var (
	_ GenericRegion = &Region{}
	_ GenericRegion = &Server{}
)
