package entities

type Folder struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	ParentFolderId     string `json:"parent_folder_id"`
	Path               string `json:"path"`
	Color              string `json:"color"`
	Description        string `json:"description"`
	IsDeleted          bool   `json:"is_deleted"`
	PermanentDeletedAt string `json:"permanent_deletion_at"`
	CreatedBy          string `json:"created_by"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}
