package entities

type UserQuota struct {
	UserID            string
	StorageLimit      int64
	StorageUsed       int64
	IsAllowedToUpload bool
	LastCalcAt        string
}
