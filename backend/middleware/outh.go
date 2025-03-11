package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

var (
	clientID     = "YOUR_CLIENT_ID"
	clientSecret = "YOUR_CLIENT_SECRET"
	redirectURI  = "http://localhost:8080/auth/callback" // Ensure this matches the one in Google Console
)

// ExchangeAuthCode exchanges the auth code for an access token
func ExchangeAuthCode(authCode string) (map[string]interface{}, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", authCode)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result, nil
}

// GetUserInfo fetches user details using the access token
func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(body, &userInfo)

	return userInfo, nil
}
