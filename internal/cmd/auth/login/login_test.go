package login

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"golang.org/x/oauth2"

	mockapiclient "github.com/zeabur/cli/mocks/pkg/api"

	"github.com/zeabur/cli/internal/cmdutil"
	mockconfig "github.com/zeabur/cli/mocks/pkg/config"
	apiClient "github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/log"
	"github.com/zeabur/cli/pkg/model"
)

var _ = Describe("LoggedIn", func() {
	const (
		token = "this_is_a_token"
		user  = "Bird"

		refreshToken = "this_is_a_refresh_token"
		tokenType    = "Bearer"
	)

	var (
		expiry     = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		oauthToken = &oauth2.Token{
			AccessToken:  token,
			Expiry:       expiry,
			TokenType:    tokenType,
			RefreshToken: refreshToken,
		}
	)

	var (
		f      *cmdutil.Factory
		opts   *LoginOptions
		buffer *zaptest.Buffer

		expectedLogs []string
		gottenLogs   []string
	)
	When("the user already logged in", func() {

		// We use Ginkgo to write tests,
		// and here is the execution order of this block:

		// BeforeEach(in this block) -> BeforeEach(in child block) ->
		// JustBeforeEach(in this block) -> It() (in child block) ->
		// JustAfterEach(in this block) ...

		// Ginkgo also has many other blocks, like AfterAll, AfterEach, etc.
		// see https://onsi.github.io/ginkgo/ for more details

		// it will be executed before each test case
		BeforeEach(func() {
			f = cmdutil.NewFactory()

			// to mock config.GetTokenString and config.GetUser
			mc := mockconfig.NewConfig(GinkgoT())
			mc.On("GetTokenString").Return(token)
			mc.On("GetUser").Return(user)
			mc.On("GetToken").Return(oauthToken)
			f.Config = mc

			// to mock client.GetUserInfo
			client := mockapiclient.NewClient(GinkgoT())
			client.On("GetUserInfo", mock.Anything).Return(&model.User{
				Name: user,
			}, nil)
			f.ApiClient = client

			// to mock client.New
			newClientFunc := func(string) apiClient.Client {
				return client
			}
			opts = &LoginOptions{
				newClient: newClientFunc,
			}

			// reset the buffer
			buffer = &zaptest.Buffer{}
		})

		JustBeforeEach(func() {
			err := runLogin(f, opts)
			Expect(err).ToNot(HaveOccurred())
			gottenLogs = buffer.Lines()
		})

		JustAfterEach(func() {
			GinkgoT().Log("gotten logs:")
			for _, gotten := range gottenLogs {
				GinkgoT().Log(gotten)
			}

			Expect(len(expectedLogs)).To(Equal(len(gottenLogs)))
			for i, expected := range expectedLogs {
				Expect(gottenLogs[i]).To(Equal(expected))
			}
		})

		Context("log level is debug", func() {
			BeforeEach(func() {
				f.Log = log.NewForUT(buffer, zapcore.DebugLevel)

				expectedLogs = []string{
					`DEBUG	Running login in non-interactive mode`,
					`DEBUG	Already logged in	{"token string": "this_is_a_token", "token detail": {"access_token":"this_is_a_token","token_type":"Bearer","refresh_token":"this_is_a_refresh_token","expiry":"2020-01-01T00:00:00Z"}, "user": {"_id":"000000000000000000000000","name":"Bird","email":"","username":"","language":"","githubID":0,"avatarUrl":"","createdAt":"0001-01-01T00:00:00Z","bannedAt":null,"agreedAt":null,"lastCheckedInAt":null,"discordID":null}}`,
					`INFO	Already logged in as Bird, if you want to use a different account, please logout first`,
				}
			})

			It("should log the correct messages", func() {
				// main logic is in BeforeEach and JustAfterEach
			})

		})

		Context("log level is info", func() {
			BeforeEach(func() {
				f.Log = log.NewForUT(buffer, zapcore.InfoLevel)

				expectedLogs = []string{
					`INFO	Already logged in as Bird, if you want to use a different account, please logout first`,
				}
			})

			It("should log the correct messages", func() {
				// main logic is in BeforeEach and JustAfterEach
			})
		})
	})

	// not completed yet
	//When("the user is not logged in(log level is info)", func() {
	//	BeforeEach(func() {
	//		mc := mockconfig.NewConfig(GinkgoT())
	//		mc.On("GetTokenString").Return("")
	//		mc.On("GetUser").Return("")
	//		mc.On("GetToken").Return(nil)
	//		f.Config = mc
	//
	//		buffer = &zaptest.Buffer{}
	//		f.Log = log.NewForUT(buffer, zapcore.InfoLevel)
	//	})
	//
	//	JustBeforeEach(func() {
	//		err := runLogin(f, opts)
	//		Expect(err).ToNot(HaveOccurred())
	//		gottenLogs = buffer.Lines()
	//	})
	//
	//	JustAfterEach(func() {
	//		GinkgoT().Log("gotten logs:")
	//		for _, gotten := range gottenLogs {
	//			GinkgoT().Log(gotten)
	//		}
	//
	//		Expect(len(expectedLogs)).To(Equal(len(gottenLogs)))
	//		for i, expected := range expectedLogs {
	//			Expect(gottenLogs[i]).To(Equal(expected))
	//		}
	//	})
	//
	//	Context("login with token in flag, env, or config file", func() {
	//		It("should log the correct messages", func() {
	//
	//		})
	//	})
	//
	//	Context("login with browser", func() {
	//		BeforeEach(func() {
	//			expectedLogs = []string{
	//				`INFO	A browser window will be opened for you to login, please confirm`,
	//				`INFO	LoggedIn successful!`,
	//				`INFO	Logged in as	{"user": "Bird", "email": "1606251360@qq.com"}`,
	//			}
	//		})
	//
	//		It("should log the correct messages", func() {
	//
	//		})
	//
	//	})
	//})
})
