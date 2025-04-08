package helpers

import (
	"io"
	"net/http"
	"strings"
)

func HttpRequest(
	method, url string,
	headers map[string]string,
	body io.Reader,
) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func ERPRequest(method, targetUrl string, body io.Reader) (resp *http.Response, err error) {
	return HttpRequest(
		method,
		"https://erp-sys.sekem.com/"+strings.TrimPrefix(targetUrl, "/"),
		map[string]string{
			"Authorization": "token 180a8380111619b:51eb64b6b260cf6",
			// "Content-Type":  "application/json",
			// "Accept":        "application/json",
		},
		body,
	)
}
