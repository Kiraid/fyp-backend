package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)
//Oauth for Signin with google

//Meta-data for request to Google
var (
	clientID     = "406024045252-6q8slt53kok07c8hjuc84v0v2lbfuknu.apps.googleusercontent.com"
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURI  = "http://localhost:8080/auth/callback" // Ensure this matches the one in Google Console
)

//Function for exchaning authorization code with access token
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

// Fetches user info using the access token and client Secret
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
