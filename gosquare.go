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
	baseURL        = "https://connect.squareup.com"
	tokenURL       = "oauth2/token"
	//refreshURL     = "oauth2/clients/%v/access-token/renew"
	webhookURL     = "/v1/%v/webhooks"
	paymentURL     = "/v1/%v/payments"
	paymentByIDURL = "/v1/%v/payments/%v"
	locationURL    = "v2/locations"
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
func (v *Square) AccessToken() (string, string, string, time.Time, error) {

	data := TokenRequest{
		ClientID:     v.ClientID,
		ClientSecret: v.ClientSecret,
		Code:         v.StoreCode,
		GrantType:    "authorization_code",
	}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	request, _ := json.Marshal(data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(request))

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(request)))
	r.Header.Add("Authorization", "Client "+data.ClientSecret)

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", "", time.Now(), err
		}

		return resp.AccessToken, resp.RefreshToken, resp.MerchantID, resp.ExpiresAt, nil
	}

	return "", "", "", time.Now(), fmt.Errorf("Failed to get access token: %s", res.Status)
}

// RefreshToken will get a new refresh token
func (v *Square) RefreshToken(refreshToken string) (string, string, string, time.Time, error) {

	data := TokenRequest{
		ClientID:     v.ClientID,
		ClientSecret: v.ClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", "", "", time.Now(), err
	}
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	request, _ := json.Marshal(data)

	client := &http.Client{}
	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(request))
	if err != nil {
		return "", "", "", time.Now(), err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(request)))

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", "", time.Now(), err
		}

		return resp.AccessToken, resp.RefreshToken, resp.MerchantID, resp.ExpiresAt, nil
	}

	resp := &ErrorResponse{}
	if err := json.Unmarshal(rawResBody, resp); err != nil {
		return "", "", "", time.Now(), fmt.Errorf("Error requesting access token: %v", err)
	}
	return "", "", "", time.Now(), fmt.Errorf("Error requesting access token: %v", resp.Message)
}

// UpdateWebHook will init the sales hook for the Square store
func (v *Square) UpdateWebHook(token string, company string, location string, paymentUpdated bool) error {
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

	client := &http.Client{}
	r, err := http.NewRequest("PUT", urlStr, bytes.NewBufferString(body))
	if err != nil {
		return err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+token)

	res, _ := client.Do(r)

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

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
	r.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp LocationList

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp.Locations, nil
	}
	return nil, fmt.Errorf("Failed to get Square Locations %s", res.Status)

}

// GetPayment will return the details of the payment
func (v *Square) GetPayment(token string, locationID string, paymentID string) (*Payment, error) {
	client := &http.Client{}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(paymentByIDURL, locationID, paymentID)
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//	fmt.Println(string(rawResBody))

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

// GetPayments will return the details of the payment
func (v *Square) GetPayments(token string, locationID string, start string, end string) ([]Payment, error) {
	client := &http.Client{}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = fmt.Sprintf(paymentURL, locationID)
	urlStr := fmt.Sprintf("%v", u)

	urlStr = urlStr + "?order=DESC&limit=200"
	if start != "" {
		urlStr = urlStr + fmt.Sprintf("&begin_time=%v&end_time=%v", start, end)
	}

	//urlStr = "https://connect.squareup.com/v1/V0CKA1WZJP1W0/payments?batch_token=Ti1kg6T7YQ63F7YcVoNJdkjPvrJtZt0B5Vi8HU1OoS8Oi51DCvGX0DxDeep69A5P1jCvo7vJYQ9Q3imPmQ4ojGusiQiAf0Mg8gTLK5xmUdrzkW5lyaJ0U1ppb4HaJL6a9Wu8uK&begin_time=2017-12-13T21%3A43%3A16.364795388Z&end_time=2018-12-13T21%3A43%3A16.364795388Z&limit=200&order=DESC"
	fmt.Println(urlStr)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp []Payment

		fmt.Println(res.Header["Link"])

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Square Payment %s", res.Status)

}
