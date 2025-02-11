package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetRequestData(r *http.Request, requestData interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestData); err != nil {
		return err
	}

	return nil
}

func WriteResponseData(w http.ResponseWriter, responseData interface{}, successStatusCode int) error {
	body, err := json.Marshal(responseData)
	if err != nil {
		return fmt.Errorf("error in marshalling response body: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	w.WriteHeader(successStatusCode)
	_, _ = w.Write(body)

	return nil
}

func WriteErrorMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"message":"%s"}`, message)
}
