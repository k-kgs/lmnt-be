package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	kmw "kayam-be/internal/middleware"
	"kayam-be/internal/storage"
)

// maxUploadBytes is the server-side backstop to the FE's client-side
// compression — never trust the client's declared size alone (root design
// doc §5).
const maxUploadBytes = 2 * 1024 * 1024 // 2MB

var allowedContentTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

type UploadHandler struct {
	Presigner *storage.Presigner
}

type requestUploadRequest struct {
	ContentType       string `json:"content_type"`
	DeclaredSizeBytes int64  `json:"declared_size_bytes"`
}

type requestUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

var (
	errUnsupportedContentType = errors.New("unsupported content_type (allowed: image/jpeg, image/png, image/webp)")
	errInvalidSize            = errors.New("declared_size_bytes must be a positive number")
	errFileTooLarge           = fmt.Errorf("declared_size_bytes exceeds the %dMB limit", maxUploadBytes/1024/1024)
)

func (h *UploadHandler) RequestURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := kmw.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	var req requestUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ext, ok := allowedContentTypes[req.ContentType]
	if !ok {
		writeError(w, http.StatusBadRequest, errUnsupportedContentType)
		return
	}
	if req.DeclaredSizeBytes <= 0 {
		writeError(w, http.StatusBadRequest, errInvalidSize)
		return
	}
	if req.DeclaredSizeBytes > maxUploadBytes {
		writeError(w, http.StatusBadRequest, errFileTooLarge)
		return
	}

	objectKey := fmt.Sprintf("checkins/%s/%s.%s", userID.String(), uuid.New().String(), ext)

	uploadURL, err := h.Presigner.PresignPut(r.Context(), objectKey, req.ContentType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, requestUploadResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	})
}
