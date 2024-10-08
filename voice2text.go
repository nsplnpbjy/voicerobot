package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	Voice2TextUrl = "https://vop.baidu.com/server_api"
	CUID          = "ijuqLN2knJfJ2BeQ0f3k5upefxbv529o"
)

func GetFileContentAsBase64(path string) (string, int) {
	srcByte, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return "", 0
	}
	return base64.StdEncoding.EncodeToString(srcByte), len(srcByte)
}

func Voice2Text(audioData []byte) (string, error) {
	fmt.Println("Audio data length:", len(audioData))

	if !TokenInfoVar.IsAuth() {
		TokenInfoVar.VoiceAuth()
	}
	type Payload struct {
		Format  string `json:"format"`
		Rate    int    `json:"rate"`
		Channel int    `json:"channel"`
		Cuid    string `json:"cuid"`
		Speech  string `json:"speech"`
		Len     int    `json:"len"`
		Token   string `json:"token"`
	}
	speech := base64.StdEncoding.EncodeToString(audioData)
	payload := Payload{
		Format:  "m4a",
		Rate:    16000,
		Channel: 1,
		Cuid:    CUID,
		Speech:  speech,
		Len:     len(audioData),
		Token:   TokenInfoVar.Access_token,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", Voice2TextUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var bodyinfo map[string]interface{}
	json.Unmarshal(body, &bodyinfo)
	if resArray, ok := bodyinfo["result"].([]interface{}); ok {
		for _, item := range resArray {
			if str, ok := item.(string); ok {
				return str, nil
			}
		}
	}
	return "", fmt.Errorf("unexpected response format: %s", string(body))
}
