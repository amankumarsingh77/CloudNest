package entities

import (
	"github.com/amankumarsingh77/cloudnest/internal/utils/json"
)

type File struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	FolderId            json.NullString `json:"folderId"`
	RemoteFileName      string          `json:"remote_file_name"`
	MimeType            string          `json:"mime_type"`
	Size                int64           `json:"size"`
	Checksum            json.NullString `json:"checksum"`
	EncryptionKey       json.NullString `json:"encryption_key"`
	Url                 string          `json:"url"`
	CreatedBy           string          `json:"created_by"`
	Path                string          `json:"path"`
	IsDeleted           bool            `json:"is_deleted"`
	DeletedAt           json.NullTime   `json:"deleted_at"`
	PermanentDeletionAt json.NullTime   `json:"permanent_deletion_at"`
	LastAccessedAt      json.NullTime   `json:"last_accessed_at"`
	CreatedAt           string          `json:"created_at"`
	UpdatedAt           json.NullTime   `json:"updated_at"`
	DownloadUrl         string          `json:"download_url"`
}
