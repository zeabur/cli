package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http"
)

//go:embed login_success_page.html
var loginSuccessHTML []byte

type CallbackServer struct {
	listener net.Listener
	tokenCh  chan TokenResponse
}

type TokenResponse struct {
	Token string
	State string
}

func NewCallbackServer() (*CallbackServer, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	tokenCh := make(chan TokenResponse, 1)

	return &CallbackServer{
		listener: listener,
		tokenCh:  tokenCh,
	}, nil
}

func (s *CallbackServer) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /callback", func(w http.ResponseWriter, r *http.Request) {
		s.tokenCh <- TokenResponse{
			Token: r.FormValue("api_key"),
			State: r.FormValue("state"),
		}

		w.Header().Add("Content-Type", "text/html")

		w.WriteHeader(200)
		_, _ = w.Write(loginSuccessHTML)
	})

	return http.Serve(s.listener, mux)
}

func (s *CallbackServer) Close() error {
	err := s.listener.Close()
	close(s.tokenCh)

	return err
}

func (s *CallbackServer) Port() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *CallbackServer) GetCallbackURL() string {
	return fmt.Sprintf("http://localhost:%d/callback", s.Port())
}

func (s *CallbackServer) WaitForToken(ctx context.Context) (TokenResponse, error) {
	select {
	case <-ctx.Done():
		return TokenResponse{}, ctx.Err()
	case token, ok := <-s.tokenCh:
		if !ok {
			return TokenResponse{}, fmt.Errorf("token channel closed")
		}

		return token, nil
	}
}
