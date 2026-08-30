package handlers

import (
	"RestFulApi-ManajemenKeuangan/config"
	"RestFulApi-ManajemenKeuangan/internal/database"
	"RestFulApi-ManajemenKeuangan/internal/models"
	"RestFulApi-ManajemenKeuangan/internal/utils"

	"context"
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

	url := googleOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

type googleUserInfo struct {
	ID			string `json:"id"`
	Email		string `json:"email"`
	Name		string `json:"name"`
	Picture		string `json:"picture"`
}

// GoogleCallback - GET /auth/google/callback
// Google redirect ke sini dengan ?code=...&state=...
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if !oauthStateStore[state] {
		utils.Error(c, http.StatusBadRequest, "state OAuth tidak valid (kemungkinan CSRF atau expired)")
		return
	}
	delete(oauthStateStore, state)

	if code == "" {
		utils.Error(c, http.StatusBadRequest, "parameter code tidak ditemukan")
		return
	}

	oauthCfg := googleOAuthConfig()
	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal exchange code ke Google: " + err.Error())
		return
	}

	client := oauthCfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal mengambil profile Google: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal membaca response Google")
		return
	}

	var gUser googleUserInfo
	if err := json.Unmarshal(body, &gUser); err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal parse profil Google")
		return
	}

	// Cari user berdasarkan GoogleID, kalau belum ada -> buat baru (upsert sederhana)
	var user models.User
	result := database.DB.Where("google_id = >", gUser.ID).First(&user)
	if result.Error != nil {
		user = models.User{
			GoogleID: 		gUser.ID,
			Email: 			gUser.Email,
			Name: 			gUser.Name,
			Avatar: 		gUser.Picture,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "gagal membuat user baru: " + err.Error())
			return
		}
	} else {
		// update data profil terbaru dari Google
		user.Name = gUser.Name
		user.Avatar = gUser.Picture
		database.DB.Save(&user)
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, config.App.JWTSecret, config.App.JWTExpiryHrs)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "gagal membuat JWT")
		return
	}

	// Redirect ke frontend membawa token, atau bisa juga langsung return JSON
	// jika frontend & backend dites terpisah (misal lewat Postman).
	if config.App.FrontendURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, config.App.FrontendURL+"?token="+jwtToken)
		return
	}

	utils.Success(c, http.StatusOK, "login berhasil", gin.H{
		"token": jwtToken,
		"user": user,
	})
}

// Me - GET /api/auth/me (protected)
func Me(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil profil", user)
}