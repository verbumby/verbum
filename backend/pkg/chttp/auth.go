package chttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"
)

// Principal represents current user
type Principal struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	Locale            string `json:"locale"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Zoneinfo          string `json:"zoneinfo"`
	UpdatedAt         uint32 `json:"updated_at"`
	EmailVerified     bool   `json:"email_verified"`
}

var cookieManager *securecookie.SecureCookie

// InitCookieManager inializes cookie manager
func InitCookieManager() {
	cookieManager = securecookie.New(
		[]byte(viper.GetString("cookie.hashKey")),
		[]byte(viper.GetString("cookie.blockKey")),
	)
}

// AuthHandler handler to authenticate user
func AuthHandler(w http.ResponseWriter, ctx *Context) error {
	defer ctx.R.Body.Close()

	if err := ctx.R.ParseForm(); err != nil {
		http.Error(w, fmt.Errorf("parse form: %w", err).Error(), http.StatusBadRequest)
		return nil
	}

	stateCookie, err := ctx.R.Cookie(viper.GetString("cookie.nameState"))
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return nil
	}
	if stateCookie.Value != ctx.R.Form.Get("state") {
		http.Error(w, "", http.StatusUnauthorized)
		return nil
	}
	stateCookie.MaxAge = -1
	stateCookie.Path = "/admin/"
	http.SetCookie(w, stateCookie)

	// exchange code for access token
	code := ctx.R.Form.Get("code")
	bodyValues := url.Values{}
	bodyValues.Set("grant_type", "authorization_code")
	bodyValues.Set("code", code)
	bodyValues.Set("redirect_uri", viper.GetString("https.canonicalAddr")+"/admin/auth")
	bodyValues.Set("client_id", viper.GetString("oauth.clientID"))
	bodyValues.Set("client_secret", viper.GetString("oauth.clientSecret"))

	req, err := http.NewRequest(http.MethodPost, viper.GetString("oauth.endpointToken"), strings.NewReader(bodyValues.Encode()))
	if err != nil {
		return fmt.Errorf("token request create: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("token request: %d %s", resp.StatusCode, string(respBody))
	}

	respData := struct {
		AccessToken string `json:"access_token"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return fmt.Errorf("token request: response decode: %w", err)
	}
	accessToken := respData.AccessToken

	// get user info by the access token
	req, err = http.NewRequest(http.MethodGet, viper.GetString("oauth.endpointUserinfo"), nil)
	if err != nil {
		return fmt.Errorf("userinfo request create: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("get userinfo request: %d %s", resp.StatusCode, string(respBody))
	}

	p := Principal{}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return fmt.Errorf("userinfo request body read: %w", err)
	}

	allowedEmails := viper.GetStringSlice("editorsWhitelist")
	allow := false
	for _, ae := range allowedEmails {
		if ae == p.Email {
			allow = true
			break
		}
	}
	if !allow {
		http.Error(w, "", http.StatusUnauthorized)
		return nil
	}

	cookie, err := encodeCookie(p)
	if err != nil {
		return fmt.Errorf("auth set cookie: %w", err)
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, ctx.R, "/admin/", http.StatusFound)
	return nil
}

func encodeCookie(p Principal) (*http.Cookie, error) {
	encoded, err := cookieManager.Encode(viper.GetString("cookie.name"), p)
	if err != nil {
		return nil, fmt.Errorf("encode cookie: %w", err)
	}

	cookie := &http.Cookie{
		Name:     viper.GetString("cookie.name"),
		Value:    encoded,
		Path:     "/admin/",
		MaxAge:   viper.GetInt("cookie.maxAge"),
		HttpOnly: true,
	}
	return cookie, nil
}

// AuthMiddleware verifies that the current user is authenticated
func AuthMiddleware(f HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, ctx *Context) error {
		if cookie, err := ctx.R.Cookie(viper.GetString("cookie.name")); err == nil {
			p := &Principal{}
			if err = cookieManager.Decode(viper.GetString("cookie.name"), cookie.Value, p); err == nil {
				ctx.P = p
				return f(w, ctx)
			}
		}

		redirectURL := viper.GetString("https.canonicalAddr") + "/admin/auth"
		query := url.Values{}
		query.Set("response_mode", "form_post")
		query.Set("response_type", "code")
		query.Set("client_id", viper.GetString("oauth.clientID"))
		query.Set("scope", "openid profile email")
		query.Set("redirect_uri", redirectURL)
		state := fmt.Sprintf("%d", rand.Int())
		query.Set("state", state)
		query.Set("nonce", fmt.Sprintf("%d", rand.Int()))

		oauthURL := viper.GetString("oauth.endpointAuth") + "?" + query.Encode()
		http.SetCookie(w, &http.Cookie{
			Name:     viper.GetString("cookie.nameState"),
			Value:    state,
			Path:     "/admin/",
			HttpOnly: true,
			MaxAge:   60 * 5,
		})
		http.Redirect(w, ctx.R, oauthURL, http.StatusFound)
		return nil
	}
}
