// Package zcontext is used to store the context information for the application.
package zcontext

type basicInfo struct {
	id   string
	name string
}

func (b *basicInfo) GetID() string {
	return b.id
}

func (b *basicInfo) GetName() string {
	return b.name
}

func (b *basicInfo) Empty() bool {
	return b == nil || b.id == "" || b.name == ""
}

func NewBasicInfo(id, name string) BasicInfo {
	return &basicInfo{
		id:   id,
		name: name,
	}
}

var _ BasicInfo = &basicInfo{}
