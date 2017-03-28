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
	webhookURL = "/v1/%v/webhooks"
	paymentURL = "/v1/%v/payments/%v"
	//webhookURL  = "v1/me/webhooks"
	locationURL = "v1/me/locations"
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

	data := TokenRequest{
		ClientID:     v.ClientID,
		ClientSecret: v.ClientSecret,
		Code:         v.StoreCode,
	}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	fmt.Printf("AccessToken %v %v\n", urlStr, data)

	request, _ := json.Marshal(data)

	fmt.Printf("AccessToken Request %v \n", string(request))

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(request))

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(request)))
	r.Header.Add("Authorization", "Client "+data.ClientSecret)

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

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("Authorization", "Client "+v.ClientSecret)

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

// UpdateWebHook will init the sales hook for the Square store
func (v *Square) UpdateWebHook(token string, company string, location string, paymentUpdated bool) error {

	fmt.Println("UpdateWebHook", token, company, location)

	//"PAYMENT_UPDATED"
	//"INVENTORY_UPDATED"
	//TIMECARD_UPDATED

	var body string

	if paymentUpdated {
		body = strconv.Quote("PAYMENT_UPDATED")
	}

	body = "[" + body + "]"

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return err
	}

	u.Path = fmt.Sprintf(webhookURL, location)
	urlStr := fmt.Sprintf("%v", u)

	fmt.Println("URL", urlStr)

	client := &http.Client{}
	r, err := http.NewRequest("PUT", urlStr, bytes.NewBufferString(body))
	if err != nil {
		return err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+token)

	res, _ := client.Do(r)
	fmt.Println(res.Status)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("BODY", string(rawResBody))

	return nil
}

// GetLocations will return the categories of the authenticated token
func (v *Square) GetLocations(token string) (Locations, error) {
	client := &http.Client{}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = locationURL
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "Bearer "+token) //v.ClientSecret)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("GetLocations Body", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Locations

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Square Locations %s", res.Status)

}

// GetPayment will return the categories of the authenticated token
func (v *Square) GetPayment(token string, locationID string, paymentID string) (*Payment, error) {
	client := &http.Client{}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(paymentURL, locationID, paymentID)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "Bearer "+token) //v.ClientSecret)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println("GetPayment Body", string(rawResBody))

	if res.StatusCode == 200 {
		resp := Payment{}

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return &resp, nil
	}
	return nil, fmt.Errorf("Failed to get Square Payment %s", res.Status)

}
