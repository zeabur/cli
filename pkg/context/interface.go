package context

type Context interface {
	GetProject() BasicInfo
	SetProject(project BasicInfo)
	ClearProject()

	GetEnvironment() BasicInfo
	SetEnvironment(environment BasicInfo)
	ClearEnvironment()

	GetService() BasicInfo
	SetService(service BasicInfo)
	ClearService()

	ClearAll()
}

type BasicInfo interface {
	GetID() string
	GetName() string
	Empty() bool
}
