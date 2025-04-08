package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ahmedsat/middleware/helpers"
	"github.com/ahmedsat/middleware/internals"
)

var debug bool = true

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fa := internals.FarmApplication{}
	err := fa.Scan(r.Body)
	if err != nil {
		httpError(w, http.StatusInternalServerError, errors.Join(
			errors.New("we can not parse the request body"),
			err,
		))
		return
	}

	err = fa.Validate()
	if err != nil {
		httpError(w, http.StatusBadRequest, errors.Join(
			errors.New("validation errors"),
			err,
		))
		return
	}

	faBytes, err := json.Marshal(fa)
	if err != nil {
		httpError(w, http.StatusInternalServerError, errors.Join(
			errors.New("can not marshal the data"),
			err,
		))
		return
	}

	res, err := helpers.ERPRequest("POST", "/api/resource/Farm Application", bytes.NewBuffer(faBytes))
	if err != nil {
		httpError(w, http.StatusInternalServerError, errors.Join(
			errors.New("can not send request to erp"),
			err,
		))
		return
	}

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		httpError(w, http.StatusInternalServerError, errors.Join(
			errors.New("there is an error with with erp"),
			err,
			errors.New(string(data)),
		))
		return
	}

	if debug {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
