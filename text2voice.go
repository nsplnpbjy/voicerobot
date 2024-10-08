package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	Text2VoiceUrl = "https://tsn.baidu.com/text2audio"
)

func Text2Voice(text string) ([]byte, error) {
	formdata := url.Values{}
	if !TokenInfoVar.IsAuth() {
		TokenInfoVar.VoiceAuth()
	}
	formdata.Set("tok", TokenInfoVar.Access_token)
	formdata.Set("cuid", CUID)
	formdata.Set("ctp", "1")
	formdata.Set("lan", "zh")
	formdata.Set("per", "0")
	formdata.Set("aue", "3")
	formdata.Set("tex", text)

	body := bytes.NewBufferString(formdata.Encode())
	req, err := http.NewRequest("POST", Text2VoiceUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType[:5] == "audio" {
		fmt.Println("Audio synthesis successful!")
		audioBinary, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read audio binary: %w", err)
		}
		return audioBinary, nil
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		return nil, fmt.Errorf("API error: %s", string(bodyBytes))
	}
}
