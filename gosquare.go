package gosquare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL    = "https://connect.squareup.com"
	tokenURL   = "oauth2/token"
	refreshURL = "oauth2/clients/%v/access-token/renew"
)

var (
	defaultSendTimeout = time.Second * 30
)

// Square The main struct of this package
type Square struct {
	StoreCode    string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
}

// NewClient will create a Square client with default values
func NewClient(code string, clientID string, clientSecret string) *Square {
	return &Square{
		StoreCode:    code,
		Timeout:      defaultSendTimeout,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// AccessToken will get a new access token
func (v *Square) AccessToken() (string, string, time.Time, error) {

	data := url.Values{}
	data.Set("code", v.StoreCode)
	data.Add("client_secret", v.ClientSecret)
	data.Add("client_id", v.ClientID)

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	fmt.Printf("AccessToken %v %v\n", urlStr, data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)

	fmt.Printf("AccessToken Body %v \n", string(rawResBody))

	if err != nil {
		return "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", time.Now(), err
		}

		return resp.AccessToken, resp.MerchantID, resp.ExpiresAt, nil
	}

	return "", "", time.Now(), fmt.Errorf("Failed to get access token: %s", res.Status)
}

// RefreshToken will get a new refresh token
func (v *Square) RefreshToken(refreshtoken string) (string, string, time.Time, error) {

	data := url.Values{}
	data.Set("access_token", refreshtoken)

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", "", time.Now(), err
	}

	u.Path = fmt.Sprintf(refreshURL, v.ClientID)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", "", time.Now(), err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", time.Now(), err
	}

	fmt.Println("BODY", string(rawResBody))

	if res.StatusCode >= 400 {
		return "", "", time.Now(), fmt.Errorf("Failed to get refresh token: %s", res.Status)
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", time.Now(), err
		}

		return resp.AccessToken, resp.MerchantID, resp.ExpiresAt, nil
	}

	return "", "", time.Now(), fmt.Errorf("Error requesting access token")
}
