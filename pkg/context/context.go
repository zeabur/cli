package context

import "github.com/spf13/viper"

const (
	KeyContext = "context"

	KeyProject     = KeyContext + ".project"
	KeyProjectID   = KeyProject + ".id"
	KeyProjectName = KeyProject + ".name"

	KeyEnvironment     = KeyContext + ".environment"
	KeyEnvironmentID   = KeyEnvironment + ".id"
	KeyEnvironmentName = KeyEnvironment + ".name"

	KeyService     = KeyContext + ".service"
	KeyServiceID   = KeyService + ".id"
	KeyServiceName = KeyService + ".name"
)

type viperContext struct {
	viper *viper.Viper
}

func NewViperContext(viper *viper.Viper) Context {
	return &viperContext{viper: viper}
}

func (c *viperContext) GetProject() BasicInfo {
	return &basicInfo{
		id:   c.viper.GetString(KeyProjectID),
		name: c.viper.GetString(KeyProjectName),
	}
}

func (c *viperContext) SetProject(project BasicInfo) {
	c.viper.Set(KeyProjectID, project.GetID())
	c.viper.Set(KeyProjectName, project.GetName())
}

func (c *viperContext) ClearProject() {
	c.viper.Set(KeyProjectID, "")
	c.viper.Set(KeyProjectName, "")
}

func (c *viperContext) GetEnvironment() BasicInfo {
	return &basicInfo{
		id:   c.viper.GetString(KeyEnvironmentID),
		name: c.viper.GetString(KeyEnvironmentName),
	}
}

func (c *viperContext) SetEnvironment(environment BasicInfo) {
	c.viper.Set(KeyEnvironmentID, environment.GetID())
	c.viper.Set(KeyEnvironmentName, environment.GetName())
}

func (c *viperContext) ClearEnvironment() {
	c.viper.Set(KeyEnvironmentID, "")
	c.viper.Set(KeyEnvironmentName, "")
}

func (c *viperContext) GetService() BasicInfo {
	return &basicInfo{
		id:   c.viper.GetString(KeyServiceID),
		name: c.viper.GetString(KeyServiceName),
	}
}

func (c *viperContext) SetService(service BasicInfo) {
	c.viper.Set(KeyServiceID, service.GetID())
	c.viper.Set(KeyServiceName, service.GetName())
}

func (c *viperContext) ClearService() {
	c.viper.Set(KeyServiceID, "")
	c.viper.Set(KeyServiceName, "")
}

func (c *viperContext) ClearAll() {
	c.ClearProject()
	c.ClearEnvironment()
	c.ClearService()
}

var _ Context = &viperContext{}
