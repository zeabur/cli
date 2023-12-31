package webapp

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

//go:embed login_success_page.html
var loginSuccessHTML []byte

// CodeResponse represents the code received by the local server's callback handler.
type CodeResponse struct {
	Code  string
	State string
}

// bindLocalServer initializes a LocalServer that will listen on a randomly available TCP port.
func bindLocalServer() (*localServer, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	return &localServer{
		SuccessPath: "/success",
		listener:    listener,
		resultChan:  make(chan CodeResponse, 1),
	}, nil
}

type localServer struct {
	listener         net.Listener
	WriteSuccessHTML func(w io.Writer)
	resultChan       chan CodeResponse

	CallbackPath string
	SuccessPath  string

	once sync.Once
}

func (s *localServer) Port() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *localServer) Close() error {
	return s.listener.Close()
}

func (s *localServer) Serve() error {
	return http.Serve(s.listener, s)
}

func (s *localServer) WaitForCode(ctx context.Context) (CodeResponse, error) {
	select {
	case <-ctx.Done():
		return CodeResponse{}, ctx.Err()
	case code := <-s.resultChan:
		return code, nil
	}
}

// ServeHTTP implements http.Handler.
func (s *localServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("method: %s, url: %s\n", r.Method, r.URL.String())

	path := r.URL.Path
	switch path {
	case s.CallbackPath:
		s.ServeCallback(w, r)
	case s.SuccessPath:
		s.ServeSuccess(w, r)
	default:
		w.WriteHeader(404)
	}
}

func (s *localServer) ServeCallback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	s.resultChan <- CodeResponse{
		Code:  params.Get("code"),
		State: params.Get("state"),
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	successURL := fmt.Sprintf("http://localhost:%d%s", s.Port(), s.SuccessPath)
	http.Redirect(w, r, successURL, http.StatusFound)
}

func (s *localServer) ServeSuccess(w http.ResponseWriter, r *http.Request) {
	if s.SuccessPath != "" && r.URL.Path != s.SuccessPath {
		w.WriteHeader(404)
		return
	}

	defer func() {
		// if method is GET, close (to ignore Options)
		if r.Method == http.MethodGet {
			closeServer := func() {
				err := s.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
			s.once.Do(func() {
				time.AfterFunc(1*time.Second, closeServer)
			})
		}
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "text/html")

	if s.WriteSuccessHTML != nil {
		s.WriteSuccessHTML(w)
	} else {
		if err := defaultSuccessHTML(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func defaultSuccessHTML(w io.Writer) error {
	_, err := w.Write(loginSuccessHTML)
	return err
}
