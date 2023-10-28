package models

type User struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type Integration struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type UserIntegration struct {
	Id          string `db:"id"`
	Integration string `db:"integration"`
	Owner       string `db:"owner"`
}

type Credential struct {
	Id          string `json:"id" db:"id"`
	Integration string `json:"integration" db:"integration"`
	Key         string `json:"key" db:"key"`
	Value       string `json:"value" db:"value"`
}

type Session struct {
	Id          string `db:"id"`
	SessionUser string `db:"session_user"`
}
