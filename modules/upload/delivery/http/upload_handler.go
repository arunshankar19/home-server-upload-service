package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/response"
)

type uploadHandler struct {
	uploadUsecase domain.UploadUsecase
}

// NewUploadHandler returns a upload handler
func NewUploadHandler(uploadUsecase domain.UploadUsecase) domain.UploadHandler {
	return &uploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

// InitiateMultipartUpload starts an multipart upload
func (h *uploadHandler) InitiateMultipartUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var initiateUploadReq domain.InitiateUploadReqDTO

	err := json.NewDecoder(r.Body).Decode(&initiateUploadReq)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	uploadID, err := h.uploadUsecase.InitiateMultipartUpload(ctx, initiateUploadReq)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedUser) {
			response.Fail(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		response.Fail(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(w, http.StatusOK, domain.InitiateUploadResponseDTO{UploadID: uploadID})
}
