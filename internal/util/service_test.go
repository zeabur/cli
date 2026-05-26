package util_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

type fakeServiceClient struct {
	api.Client

	getServiceCalls []getServiceCall
	getServiceRet   *model.Service
	getServiceErr   error

	listAllServicesProject string
	listAllServices        model.Services
	listAllServicesErr     error
}

type getServiceCall struct{ id, owner, projectName, name string }

func (c *fakeServiceClient) GetService(_ context.Context, id, owner, projectName, name string) (*model.Service, error) {
	c.getServiceCalls = append(c.getServiceCalls, getServiceCall{id, owner, projectName, name})
	return c.getServiceRet, c.getServiceErr
}

func (c *fakeServiceClient) ListAllServices(_ context.Context, projectID string) (model.Services, error) {
	c.listAllServicesProject = projectID
	return c.listAllServices, c.listAllServicesErr
}

// TestGetServiceByName_Personal: ownerID == "" preserves the personal
// `service(owner, projectName, name)` query exactly as before — same args,
// same backend call. Existing personal users see zero behavior change.
func TestGetServiceByName_Personal(t *testing.T) {
	want := &model.Service{ID: "65aa1234567890abcdef1234", Name: "web"}
	c := &fakeServiceClient{getServiceRet: want}

	got, err := util.GetServiceByName(c, "", "alice", "api", "65cc1234567890abcdef0000", "web")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if len(c.getServiceCalls) != 1 {
		t.Fatalf("expected 1 GetService call, got %d", len(c.getServiceCalls))
	}
	if c.getServiceCalls[0] != (getServiceCall{id: "", owner: "alice", projectName: "api", name: "web"}) {
		t.Fatalf("call args = %+v", c.getServiceCalls[0])
	}
	if c.listAllServicesProject != "" {
		t.Fatalf("personal path must not call ListAllServices, got project=%q", c.listAllServicesProject)
	}
}

// TestGetServiceByName_TeamFound: team workspace uses projectID-scoped
// ListAllServices and matches by name. Personal username argument is
// deliberately unused on this path.
func TestGetServiceByName_TeamFound(t *testing.T) {
	projectID := "65cc1234567890abcdef0000"
	svc1 := &model.Service{ID: "65aa1234567890abcdef1234", Name: "web"}
	svc2 := &model.Service{ID: "65bb5678901234abcdef5678", Name: "worker"}
	c := &fakeServiceClient{listAllServices: model.Services{svc1, svc2}}

	got, err := util.GetServiceByName(c, "65cc1234567890abcdefffff", "alice", "api", projectID, "worker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != svc2.ID {
		t.Fatalf("got ID %q, want %q", got.ID, svc2.ID)
	}
	if c.listAllServicesProject != projectID {
		t.Fatalf("ListAllServices called with project %q, want %q", c.listAllServicesProject, projectID)
	}
	if len(c.getServiceCalls) != 0 {
		t.Fatalf("personal GetService must not run in team workspace; got %d calls", len(c.getServiceCalls))
	}
}

// TestGetServiceByName_TeamWithoutProjectContext: a service lookup by name
// in a team workspace needs a project context (services are scoped to
// projects). Surface the actionable error rather than silently falling
// through to the personal account.
func TestGetServiceByName_TeamWithoutProjectContext(t *testing.T) {
	_, err := util.GetServiceByName(&fakeServiceClient{}, "65cc1234567890abcdefffff", "alice", "api", "", "web")
	if err == nil {
		t.Fatal("want error when team workspace has no project context")
	}
	if !strings.Contains(err.Error(), "without a project context") {
		t.Fatalf("error = %v, want 'without a project context' message", err)
	}
}

// TestGetServiceByName_TeamNotFound: name missing in the project's service
// list errors with a project-scoped message, not a personal-account 404.
func TestGetServiceByName_TeamNotFound(t *testing.T) {
	c := &fakeServiceClient{listAllServices: model.Services{{ID: "x", Name: "web"}}}

	_, err := util.GetServiceByName(c, "65cc1234567890abcdefffff", "alice", "api", "65cc1234567890abcdef0000", "missing")
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "no service named") {
		t.Fatalf("error = %v, want 'no service named ...'", err)
	}
}

// TestGetServiceByName_TeamListErr: a backend failure on the team path must
// propagate, not silently fall through to a personal-username lookup.
func TestGetServiceByName_TeamListErr(t *testing.T) {
	c := &fakeServiceClient{listAllServicesErr: errors.New("boom")}

	_, err := util.GetServiceByName(c, "65cc1234567890abcdefffff", "alice", "api", "65cc1234567890abcdef0000", "web")
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !strings.Contains(err.Error(), "list services in project") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %v, want wrapped boom", err)
	}
	if len(c.getServiceCalls) != 0 {
		t.Fatalf("personal fallback must not run; got %d calls", len(c.getServiceCalls))
	}
}
