package fill

import (
	"fmt"
	"github.com/zeabur/cli/pkg/selector"
	"github.com/zeabur/cli/pkg/zcontext"
)

type ParamFiller interface {
	// Project fills the projectID if it is empty by asking user to select a project
	Project(projectID *string) (changed bool, err error)
	// ProjectByName makes sure either projectID or projectName is not empty
	// if necessary, it will ask user to select a project first
	ProjectByName(projectID, projectName *string) (changed bool, err error)
	// Environment fills the environmentID if it is empty by asking user to select an environment,
	// when the projectID is not empty, it will ask user to select a project first
	Environment(projectID, environmentID *string) (changed bool, err error)
	// Service fills the serviceID if it is empty by asking user to select a service,
	// when the projectID is not empty, it will ask user to select a project first
	Service(projectID, serviceID *string) (changed bool, err error)
	// ServiceByName makes sure either serviceID or serviceName is not empty by asking user to select a service,
	// if necessary, it will ask user to select a project first
	ServiceByName(projectCtx zcontext.Context, serviceID, serviceName *string) (changed bool, err error)
	// ServiceWithEnvironment fills the serviceID and environmentID if they are empty by asking user to select a service and an environment,
	// when the projectID is not empty, it will ask user to select a project first
	ServiceWithEnvironment(projectID, serviceID, environmentID *string) (changed bool, err error)
	// ServiceByNameWithEnvironment behaves like ServiceByName, but it will also fill the environmentID if it is empty
	ServiceByNameWithEnvironment(projectCtx zcontext.Context, serviceID, serviceName, environmentID *string) (changed bool, err error)
}

type paramFiller struct {
	selector selector.Selector
}

func NewParamFiller(selector selector.Selector) ParamFiller {
	return &paramFiller{selector: selector}
}

func (f *paramFiller) Project(projectID *string) (changed bool, err error) {
	if err = paramNilCheck(projectID); err != nil {
		return false, err
	}

	if *projectID != "" {
		return false, nil
	}

	_, project, err := f.selector.SelectProject()
	if err != nil {
		return false, err
	}

	*projectID = project.ID

	return true, nil
}

func (f *paramFiller) ProjectByName(projectID, projectName *string) (changed bool, err error) {
	if err = paramNilCheck(projectID, projectName); err != nil {
		return false, err
	}

	if *projectID != "" || *projectName != "" {
		return false, nil
	}

	_, project, err := f.selector.SelectProject()
	if err != nil {
		return false, err
	}

	*projectID = project.ID
	*projectName = project.Name

	return true, nil
}

func (f *paramFiller) Environment(projectID, environmentID *string) (changed bool, err error) {
	if err = paramNilCheck(projectID, environmentID); err != nil {
		return false, err
	}

	if *environmentID != "" {
		return false, nil
	}

	// if projectID is empty, ask user to select a project first
	if _, err = f.Project(projectID); err != nil {
		return false, err
	}

	_, environment, err := f.selector.SelectEnvironment(*projectID)
	if err != nil {
		return false, err
	}

	*environmentID = environment.ID

	return true, nil
}

func (f *paramFiller) Service(projectID, serviceID *string) (changed bool, err error) {
	if err = paramNilCheck(projectID, serviceID); err != nil {
		return false, err
	}

	// if projectID is empty, ask user to select a project first
	if _, err = f.Project(projectID); err != nil {
		return false, err
	}

	_, service, err := f.selector.SelectService(*projectID)
	if err != nil {
		return false, err
	}

	*serviceID = service.ID

	return true, nil
}

func (f *paramFiller) ServiceByName(projectCtx zcontext.Context, serviceID, serviceName *string) (changed bool, err error) {
	if err := paramNilCheck(serviceID, serviceName); err != nil {
		return false, err
	}

	// if serviceID is not empty, do nothing
	if *serviceID != "" {
		return false, nil
	}

	// 1. if service id is empty, service name is empty,
	// we should ask user to select a service by project id
	// 2. if service id is empty, service name is not empty,
	// we should use project id and service name to specify a service

	// Therefore, we should make sure project id is not empty
	if projectCtx.GetProject().Empty() {
		project, _, err := f.selector.SelectProject()
		if err != nil {
			return false, err
		}
		// set project to projectCtx, so that we can use it later
		projectCtx.SetProject(project)
	}

	// if service name is empty, ask user to select a service by project id
	if *serviceName == "" {
		service, _, err := f.selector.SelectService(projectCtx.GetProject().GetID())
		if err != nil {
			return false, err
		}

		*serviceID = service.GetID()
		*serviceName = service.GetName()
	}

	// if service name is not empty, do nothing,
	// we have already set project id

	return true, nil
}

func (f *paramFiller) ServiceWithEnvironment(projectID, serviceID, environmentID *string) (changed bool, err error) {
	if err := paramNilCheck(projectID, serviceID, environmentID); err != nil {
		return false, err
	}

	if *serviceID != "" && *environmentID != "" {
		return false, nil
	}

	// if projectID is empty, ask user to select a project first
	if _, err := f.Project(projectID); err != nil {
		return false, err
	}

	if _, err := f.Environment(projectID, environmentID); err != nil {
		return false, err
	}

	if _, err := f.Service(projectID, serviceID); err != nil {
		return false, err
	}

	return true, nil
}

func (f *paramFiller) ServiceByNameWithEnvironment(
	projectCtx zcontext.Context, serviceID, serviceName, environmentID *string) (changed bool, err error) {

	if err = paramNilCheck(serviceID, serviceName, environmentID); err != nil {
		return false, err
	}

	changed1, err := f.ServiceByName(projectCtx, serviceID, serviceName)
	if err != nil {
		return false, err
	}

	projectID := projectCtx.GetProject().GetID()

	changed2, err := f.Environment(&projectID, environmentID)
	if err != nil {
		return false, err
	}

	return changed1 || changed2, nil

}

func paramNilCheck(params ...*string) error {
	for _, param := range params {
		if param == nil {
			return fmt.Errorf("param cannot be nil")
		}
	}
	return nil
}
