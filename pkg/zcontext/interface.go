package zcontext

// Context represents the current context of the CLI, including the current project, environment, service, etc.
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

// BasicInfo represents the basic information of a resource.
type BasicInfo interface {
	GetID() string
	GetName() string
	Empty() bool
}
