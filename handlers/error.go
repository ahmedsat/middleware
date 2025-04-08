package handlers

import (
	"fmt"
	"net/http"
)

func httpError(w http.ResponseWriter, statusCode int, err error) bool {
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}
