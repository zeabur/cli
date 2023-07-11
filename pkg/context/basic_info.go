package context

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
	return b.id == "" || b.name == ""
}

var _ BasicInfo = &basicInfo{}
