package net_ops

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func PutDataOnPaste(payload string) (string, error) {
	urlPost := "https://paste.netris.ru/add"

	headers := map[string]string{
		"Host":             "paste.netris.ru",
		"Origin":           "https://paste.netris.ru",
		"Referer":          "https://paste.netris.ru",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"X-Requested-With": "XMLHttpRequest",
	}

	data := url.Values{}
	data.Set("text", payload)

	request, err := http.NewRequest("POST", urlPost, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	pasteLink, _ := io.ReadAll(response.Body)

	return string(pasteLink), err
}
