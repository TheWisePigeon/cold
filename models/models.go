package models

type LoginPayload struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type User struct {
	Username           string `json:"username" db:"username"`
	Password           string `json:"password" db:"password"`
	FirstTimeLoggingIn bool   `db:"first_time_logging_in"`
}

type Session struct {
	Id   string `json:"id" db:"id"`
	User string `json:"user" db:"user"`
}

type GCP_Config struct {
	Id                        int    `json:"id" db:"id"`
	ServiceAccountKey         string `json:"service_account_key" db:"service_account_key"`
	ProjectId                 string `json:"project_id" db:"project_id"`
	BucketName                string `json:"bucket_name" db:"bucket_name"`
	LastUpdatedServiceAccount string `json:"last_updated_service_account" db:"last_updated_service_account"`
}

type GCPPayload struct {
	BucketName                string `form:"gcp_bucket_name"`
	ProjectId                 string `form:"gcp_project_id"`
	UploadedServiceAccountKey string `form:"gcp_service_key_uploaded"`
	ConfigId                  string `form:"gcp_config_id"`
}
