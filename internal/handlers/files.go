package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/amankumarsingh77/cloudnest/internal/domain/entities"
	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type FilePayload struct {
	Name           string          `json:"name" validate:"required"`
	Path           string          `json:"path" validate:"required"`
	FolderId       json.NullString `json:"folder_id"`
	RemoteFileName string          `json:"remote_file_name" validate:"required"`
	MimeType       string          `json:"mime_type" validate:"required"`
	Size           int64           `json:"size" validate:"required"`
	Checksum       json.NullString `json:"checksum" validate:"required"`
	Url            string          `json:"url" validate:"required"`
}

type updateFilePayload struct {
	Name           *string          `json:"name"`
	Path           *string          `json:"path"`
	FolderId       *json.NullString `json:"folder_id"`
	RemoteFileName *string          `json:"remote_file_name"`
	MimeType       *string          `json:"mime_type"`
	Size           *int64           `json:"size"`
	Checksum       *json.NullString `json:"checksum"`
	Url            *string          `json:"url"`
}

type presignedUrlPayload struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func (h *Handler) GetPresignedUrlHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r.Context())
	var payload presignedUrlPayload
	if err := json.ReadJson(w, r, &payload); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, "invalid payload")
		return
	}
	ctx := r.Context()

	filepath := fmt.Sprintf("%s/%d/%s", user.ID, time.Now().Unix(), payload.Filename)
	quota, err := h.Services.DB.UserQuota.GetUserQuota(ctx, user.ID)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !quota.IsAllowedToUpload {
		json.WriteJsonError(w, http.StatusForbidden, "quota exceeded : please upgrade")
		return
	}
	presignedUrl, err := h.Services.Storage.GetPresignedURL(ctx, filepath, payload.Size)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]string{
		"url":              presignedUrl,
		"remote_file_name": filepath,
		"created_at":       time.Now().Format(time.RFC3339),
	}
	if err := json.WriteJson(w, http.StatusOK, response); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	var filePayload FilePayload
	user := getUserFromCtx(r.Context())
	if err := json.ReadJson(w, r, &filePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(&filePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	file := &entities.File{
		Name:           filePayload.Name,
		Path:           filePayload.Path,
		FolderId:       filePayload.FolderId,
		RemoteFileName: filePayload.RemoteFileName,
		MimeType:       filePayload.MimeType,
		Size:           filePayload.Size,
		Checksum:       filePayload.Checksum,
		Url:            filePayload.Url,
		CreatedBy:      user.ID,
	}
	if err := h.Services.DB.File.CreateFile(ctx, file); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.WriteJson(w, http.StatusCreated, file); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	file := getFileFromCtx(r.Context())
	var updatePayload updateFilePayload
	if err := json.ReadJson(w, r, &updatePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(&updatePayload); err != nil {
		json.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if updatePayload.Name != nil {
		file.Name = *updatePayload.Name
	}
	if updatePayload.Path != nil {
		file.Path = *updatePayload.Path
	}
	if updatePayload.FolderId != nil {
		file.FolderId = *updatePayload.FolderId
	}
	if updatePayload.RemoteFileName != nil {
		file.RemoteFileName = *updatePayload.RemoteFileName
	}
	if updatePayload.MimeType != nil {
		file.MimeType = *updatePayload.MimeType
	}
	if updatePayload.Size != nil {
		file.Size = *updatePayload.Size
	}
	if updatePayload.Checksum != nil {
		file.Checksum = *updatePayload.Checksum
	}
	if updatePayload.Url != nil {
		file.Url = *updatePayload.Url
	}
	ctx := r.Context()

	if err := h.Services.DB.File.UpdateFile(ctx, file); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.WriteJson(w, http.StatusOK, file); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) GetFileHandler(w http.ResponseWriter, r *http.Request) {
	fileId := getFileFromCtx(r.Context()).ID
	ctx := r.Context()
	file, err := h.Services.DB.File.GetFileById(ctx, fileId)
	if err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}
	if err := json.WriteJson(w, http.StatusOK, file); err != nil {
		json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
	}

}

func (h *Handler) FileContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileId := chi.URLParam(r, "fileID")
		if fileId == "" {
			json.WriteJsonError(w, http.StatusBadRequest, "file id is required")
			return
		}
		if _, err := uuid.Parse(fileId); err != nil {
			json.WriteJsonError(w, http.StatusBadRequest, "file id is invalid")
			return
		}
		ctx := r.Context()
		file, err := h.Services.DB.File.GetFileById(ctx, fileId)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				json.WriteJsonError(w, http.StatusNotFound, err.Error())
			default:
				json.WriteJsonError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		ctx = context.WithValue(ctx, "file", file)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getFileFromCtx(ctx context.Context) *entities.File {
	return ctx.Value("file").(*entities.File)
}
