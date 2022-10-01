package ht

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/piotrpersona/gorr/log"
)

func HandleJSON[B, R any](handler func(body B, request *http.Request) (response R, status int, err error)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writeError(writer, "error reading body bytes", err, http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var requestBody B
		err = json.Unmarshal(bodyBytes, &requestBody)
		if err != nil {
			writeError(writer, "cannot unmarshal body", err, http.StatusInternalServerError)
			return
		}

		response, status, err := handler(requestBody, request)
		if err != nil {
			writeError(writer, "error while processing", err, http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(status)
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			writeError(writer, "error encoding response controller", err, http.StatusInternalServerError)
			return
		}
	}
}

func writeError(writer http.ResponseWriter, msg string, err error, code int) {
	log.Errorf("%s: %s", msg, err)
	http.Error(writer, msg, http.StatusInternalServerError)
}
