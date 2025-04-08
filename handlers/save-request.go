package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func SaveRequest(w http.ResponseWriter, r *http.Request) {
	dirName := fmt.Sprintf("log/%s", time.Now().Format("2006-01-02_15-04-05"))
	err := os.MkdirAll(dirName, os.ModePerm)
	if httpError(w, http.StatusInternalServerError, err) {
		return
	}

	headersFile, err := os.Create(dirName + "/headers.txt")
	if httpError(w, http.StatusInternalServerError, err) {
		return
	}
	defer headersFile.Close()
	for key, value := range r.Header {
		fmt.Fprintf(headersFile, "%s: %s\n", key, value)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if httpError(w, http.StatusInternalServerError, err) {
		return
	}
	r.Body.Close() // close the original body

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	bodyFile, err := os.Create(dirName + "/body.json")
	if httpError(w, http.StatusInternalServerError, err) {
		return
	}
	defer bodyFile.Close()
	_, err = bodyFile.Write(bodyBytes)
	if httpError(w, http.StatusInternalServerError, err) {
		return
	}
}
