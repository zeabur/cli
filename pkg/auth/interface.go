package auth

import "golang.org/x/oauth2"

// Client represents the client which can help user login and refresh token
type Client interface {
	Login() (token *oauth2.Token, err error)
	RefreshToken(old *oauth2.Token) (new *oauth2.Token, err error)
}
