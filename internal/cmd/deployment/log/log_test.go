package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	deploymentLog "github.com/zeabur/cli/internal/cmd/deployment/log"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/printer"
)

// stubClient stubs only the calls a deployment-id lookup should need. Any
// other Client method hits the nil embedded interface and panics.
type stubClient struct {
	api.Client

	deployment *model.Deployment

	gotServiceID, gotEnvironmentID, gotDeploymentID string
}

func (c *stubClient) GetDeployment(_ context.Context, id string) (*model.Deployment, error) {
	if c.deployment == nil || c.deployment.ID != id {
		return nil, nil
	}
	return c.deployment, nil
}

func (c *stubClient) GetRuntimeLogs(_ context.Context, serviceID, environmentID, deploymentID string) (model.Logs, error) {
	c.gotServiceID = serviceID
	c.gotEnvironmentID = environmentID
	c.gotDeploymentID = deploymentID
	return model.Logs{}, nil
}

func TestLogResolvesServiceFromDeploymentID(t *testing.T) {
	t.Parallel()

	for _, interactive := range []bool{false, true} {
		client := &stubClient{deployment: &model.Deployment{
			ID:            "6aa50431fbf9c810b64d3a2f",
			ProjectID:     "6aa50431fbf9c810b64d3a20",
			ServiceID:     "6aa50431fbf9c810b64d3a21",
			EnvironmentID: "6aa50431fbf9c810b64d3a22",
		}}
		f := &cmdutil.Factory{
			ApiClient: client,
			Log:       zap.NewNop().Sugar(),
			Printer:   printer.New(),
		}
		f.Interactive = interactive
		f.JSON = true

		cmd := deploymentLog.NewCmdLog(f)
		cmd.SetArgs([]string{"--deployment-id", "6aa50431fbf9c810b64d3a2f"})
		require.NoError(t, cmd.Execute(), "interactive=%v", interactive)

		require.Equal(t, "6aa50431fbf9c810b64d3a21", client.gotServiceID, "interactive=%v", interactive)
		require.Equal(t, "6aa50431fbf9c810b64d3a22", client.gotEnvironmentID, "interactive=%v", interactive)
		require.Equal(t, "6aa50431fbf9c810b64d3a2f", client.gotDeploymentID, "interactive=%v", interactive)
	}
}

func TestLogUnknownDeploymentID(t *testing.T) {
	t.Parallel()

	f := &cmdutil.Factory{
		ApiClient: &stubClient{},
		Log:       zap.NewNop().Sugar(),
		Printer:   printer.New(),
	}

	cmd := deploymentLog.NewCmdLog(f)
	cmd.SetArgs([]string{"--deployment-id", "6aa50431fbf9c810b64d3a2f"})
	require.ErrorContains(t, cmd.Execute(), "not found")
}
