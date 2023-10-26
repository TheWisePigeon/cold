package models

type User struct {
  Id string `json:"id" db:"id"`
  Username string `json:"username" db:"username"`
  Password string `json:"password" db:"password"`
}

type Integration struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Owner string `json:"owner" db:"owner"`
}

type Credentials struct {
	Id          string `json:"id" db:"id"`
	Integration string `json:"integration" db:"integration"`
	Key         string `json:"key" db:"key"`
	Value       string `json:"value" db:"value"`
}
