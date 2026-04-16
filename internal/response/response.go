package response

import (
	"encoding/json"
	"net/http"
)

type genericFailedResponse struct {
	Success    bool   `json:"success"` // in failed response this will be false...
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type validationFailResponse struct {
	genericFailedResponse
	Message string            `json:"message"`
	Details []ErrFieldDetails `json:"details"`
}

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
}

type ErrFieldDetails struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

// Fail sends out a generic failed response with a message and status code
func Fail(w http.ResponseWriter, statusCode int, message string) {
	res := genericFailedResponse{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resJson)
}

// ValidationFail returns error response when validation error occurs
func ValidationFail(w http.ResponseWriter, statusCode int, message string, errs []ErrFieldDetails) {
	res := validationFailResponse{
		Details: errs,
	}
	res.Success = false
	res.StatusCode = statusCode
	if message != "" {
		res.Message = message
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resJson)
}

// Success sends a success response with or without data
func Success(w http.ResponseWriter, statusCode int, data any) {

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	res := successResponse{
		Success: true,
		Data:    data,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resJson)
}
