package fill

import (
	"fmt"
	"github.com/zeabur/cli/pkg/selector"
)

type ParamFiller interface {
	// Project fills the projectID if it is empty by asking user to select a project
	Project(projectID *string) (changed bool, err error)
	// Environment fills the environmentID if it is empty by asking user to select an environment,
	// when the projectID is not empty, it will ask user to select a project first
	Environment(projectID, environmentID *string) (changed bool, err error)
	// Service fills the serviceID if it is empty by asking user to select a service,
	// when the projectID is not empty, it will ask user to select a project first
	Service(projectID, serviceID *string) (changed bool, err error)
	// ServiceByName makes sure ①project is not empty, ②serviceID and serviceName are not empty at the same time
	// 1. If serviceID and serviceName are both empty, it will ask user to select a service
	// (when projectID is not empty, it will ask user to select a project first)
	// 2. If serviceID is empty and serviceName is not empty: if projectID is not empty, do nothing,
	// otherwise ask user to select a project first
	ServiceByName(projectID, serviceID, serviceName *string) (changed bool, err error)
	// ServiceWithEnvironment fills the serviceID and environmentID if they are empty by asking user to select a service and an environment,
	// when the projectID is not empty, it will ask user to select a project first
	ServiceWithEnvironment(projectID, serviceID, environmentID *string) (changed bool, err error)
	// ServiceByNameWithEnvironment makes sure
	// 1. projectID and environmentID are not empty,
	// 2. serviceID and serviceName are not empty at the same time
	ServiceByNameWithEnvironment(projectID, serviceID, serviceName, environmentID *string) (changed bool, err error)
}

type paramFiller struct {
	selector selector.Selector
}

func NewParamFiller(selector selector.Selector) ParamFiller {
	return &paramFiller{selector: selector}
}

func (f *paramFiller) Project(projectID *string) (changed bool, err error) {
	if err := paramNilCheck(projectID); err != nil {
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

func (f *paramFiller) Environment(projectID, environmentID *string) (changed bool, err error) {
	if err := paramNilCheck(projectID, environmentID); err != nil {
		return false, err
	}

	if *environmentID != "" {
		return false, nil
	}

	// if projectID is empty, ask user to select a project first
	if _, err := f.Project(projectID); err != nil {
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
	if err := paramNilCheck(projectID, serviceID); err != nil {
		return false, err
	}

	// if projectID is empty, ask user to select a project first
	if _, err := f.Project(projectID); err != nil {
		return false, err
	}

	_, service, err := f.selector.SelectService(*projectID)
	if err != nil {
		return false, err
	}

	*serviceID = service.ID

	return true, nil
}

func (f *paramFiller) ServiceByName(projectID, serviceID, serviceName *string) (changed bool, err error) {
	if err := paramNilCheck(projectID, serviceID, serviceName); err != nil {
		return false, err
	}

	if *serviceID != "" {
		return false, nil
	}

	if *serviceName != "" {
		if *projectID != "" {
			return false, nil
		}

		// if projectID is empty, ask user to select a project first
		if _, err := f.Project(projectID); err != nil {
			return false, err
		}
	}

	// service name && service id are both empty

	if *projectID == "" {
		if _, err := f.Project(projectID); err != nil {
			return false, err
		}
	}

	service, _, err := f.selector.SelectService(*projectID)
	if err != nil {
		return false, err
	}

	*serviceID = service.GetID()
	*serviceName = service.GetName()

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

func (f *paramFiller) ServiceByNameWithEnvironment(projectID, serviceID, serviceName, environmentID *string) (changed bool, err error) {
	if err := paramNilCheck(projectID, serviceID, serviceName, environmentID); err != nil {
		return false, err
	}

	if _, err := f.ServiceByName(projectID, serviceID, serviceName); err != nil {
		return false, err
	}

	if _, err := f.Environment(projectID, environmentID); err != nil {
		return false, err
	}

	return true, nil
}

func paramNilCheck(params ...*string) error {
	for _, param := range params {
		if param == nil {
			return fmt.Errorf("param cannot be nil")
		}
	}
	return nil
}
