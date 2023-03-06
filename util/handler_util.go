package util

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

func GetEncryptedBodyOrFail(response http.ResponseWriter, request *http.Request) (string, bool) {
    if request.ContentLength == 0 {
    	http.Error(response, "Request body is required", http.StatusBadRequest)
    	return "", false
    }

	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Failed to read request body", http.StatusInternalServerError)
		return "", false
	}
	defer request.Body.Close()

	body := string(bodyBytes)

	if _, err := base64.StdEncoding.DecodeString(body); err != nil {
    	http.Error(response, "Invalid request body", http.StatusBadRequest)
    	return "", false
	}

	return body, true
}