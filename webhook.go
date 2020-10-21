package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func sendWebhook(webhookUrl string, message string) {
	resp, err := http.Get(webhookUrl + url.QueryEscape(message))
	if err != nil {
		fmt.Println("error:", err)
	}
	resp.Body.Close()
}
