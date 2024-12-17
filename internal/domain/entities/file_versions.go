package entities

type FileVersion struct {
	ID             string `json:"id"`
	FileId         string `json:"fileId"`
	VersionNum     int64  `jons:"version_number"`
	RemoteFileName string `json:"remote_file_name"`
	Size           int64  `json:"size"`
	CreatedAt      string `json:"created_at"`
}
