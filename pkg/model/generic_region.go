package model

type GenericRegion interface {
	GetID() string
	String() string
}

var (
	_ GenericRegion = &Region{}
	_ GenericRegion = &Server{}
)
