package handlers

import (
	"RestFulApi-ManajemenKeuangan/config"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	googleOAuth "golang.org/x/oauth2/google"
)

// state disimpan sangat sederhana (in-memory) hanya untuk keperluan latihan.
// Untuk production sebaiknya simpan di session/cookie yang signed, atau di cache (redis) dengan TTL.
var oauthStateStore = map[string]bool{}

func googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:		config.App.GoogleClientID,
		ClientSecret:	config.App.GoogleClientSecret,
		RedirectURL:	config.App.GoogleRedirectURL,
		Scopes:			[]string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:		googleOAuth.Endpoint,
	}
}

func randomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GoogleLogin - GET /auth/google/login
// Meredirect user ke halaman consent Google.
func GoogleLogin(c *gin.Context) {
	state := randomState()
	oauthStateStore[state] = true

	url := googleOAuthConfig().AuthCodeURL()
}