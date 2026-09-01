package login_test

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/zeabur/cli/internal/cmd/auth/login"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/auth"
	"github.com/zeabur/cli/pkg/log"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/zcontext"
)

type stubConfig struct {
	token string
}

func (s *stubConfig) GetTokenString() string       { return s.token }
func (s *stubConfig) SetTokenString(token string)  { s.token = token }
func (s *stubConfig) GetUser() string              { return "" }
func (s *stubConfig) SetUser(string)               {}
func (s *stubConfig) GetUsername() string          { return "" }
func (s *stubConfig) SetUsername(string)           {}
func (s *stubConfig) GetContext() zcontext.Context { return zcontext.NewViperContext(viper.New()) }
func (s *stubConfig) Write() error                 { return nil }

type stubClient struct {
	api.Client
}

func (stubClient) GetUserInfo(context.Context) (*model.User, error) {
	return &model.User{Name: "tester"}, nil
}

func runLoginWithStoredToken(t *testing.T, token string, interactive bool) string {
	t.Helper()

	buf := &zaptest.Buffer{}
	f := &cmdutil.Factory{
		Log:    log.NewForUT(buf, zapcore.InfoLevel),
		Config: &stubConfig{token: token},
	}
	f.Interactive = interactive

	opts := &login.Options{NewClient: func(string) api.Client { return stubClient{} }}
	if err := login.RunLogin(f, opts); err != nil {
		t.Fatalf("RunLogin returned error: %v", err)
	}
	return buf.String()
}

func TestRunLogin_NonInteractiveLegacyTokenWarns(t *testing.T) {
	out := runLoginWithStoredToken(t, "sk-legacy", false)

	if !strings.Contains(out, "Already logged in as tester") {
		t.Errorf("expected already-logged-in notice, got:\n%s", out)
	}
	if !strings.Contains(out, auth.LegacyAPIKeyDeprecationMessage) {
		t.Errorf("expected legacy deprecation warning, got:\n%s", out)
	}
}

func TestRunLogin_NonInteractiveAccessTokenDoesNotWarn(t *testing.T) {
	out := runLoginWithStoredToken(t, "zat_user_secret", false)

	if strings.Contains(out, auth.LegacyAPIKeyDeprecationMessage) {
		t.Errorf("unexpected legacy deprecation warning, got:\n%s", out)
	}
}
