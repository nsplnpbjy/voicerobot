package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	Client_id     = "Ib8xBUo5S2LDpxepCndTkUe5"
	Client_secret = "XEWkAwSiTMi4iGhjtaS4YBDl2bq7esvU"
	TokenUrl      = "https://aip.baidubce.com/oauth/2.0/token"
)

type TokenInfo struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
	Last_in      int64
}

var (
	TokenInfoVar = TokenInfo{}
)

func (T *TokenInfo) IsAuth() bool {
	return time.Now().Unix()-T.Last_in < T.Expires_in
}

func (T *TokenInfo) VoiceAuth() {
	params := url.Values{}
	params.Add("client_id", Client_id)
	params.Add("client_secret", Client_secret)
	params.Add("grant_type", "client_credentials")
	url := TokenUrl + "?" + params.Encode()

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	err = json.Unmarshal(body, T)
	T.Last_in = time.Now().Unix()
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

}
