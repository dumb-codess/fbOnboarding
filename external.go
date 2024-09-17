package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int
}

type AppHandler func(*http.Request) (*Response, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := fn(r)

	if err != nil {
		errResponse := writeErrorResponse(ctx, err, w)
		if _, err := w.Write(errResponse); err != nil {
			logger.Fatal().Msgf("http response writing failed: %v", err)
		}
		return
	}

	w.WriteHeader(resp.StatusCode)
	httpResponse := make(map[string]interface{})
	httpResponse["data"] = resp.Data
	httpResponse["status"] = resp.StatusCode

	response, err := json.Marshal(httpResponse)
	if err != nil {
		errResponse := writeErrorResponse(ctx, err, nil)

		if _, err := w.Write(errResponse); err != nil {
			logger.Error().Msgf("failed to marshal response: %v", err)
		}

		return
	}

	if _, err := w.Write(response); err != nil {
		logger.Error().Msgf("failed to write response: %v", err)

	}

}

func writeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) []byte {
	w.WriteHeader(http.StatusInternalServerError)

	errResponse := make(map[string]interface{})
	errResponse["message"] = "Internal Server Error"
	errResponse["error_message"] = err.Error()
	errResponse["code"] = http.StatusInternalServerError

	response, er := json.Marshal(errResponse)
	if er != nil {
		return []byte(http.StatusText(http.StatusInternalServerError))
	}
	return response
}

func writeBadResponse(ctx context.Context, err error, w http.ResponseWriter) {
	errResp := writeErrorResponse(ctx, err, w)
	if _, err := w.Write(errResp); err != nil {
		logger.Error().Msgf("failed to marshal response: %v", err)
	}
}
