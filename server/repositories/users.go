package repositories

import (
	"cold/models"
)

func GetUserByName(username string) error {
	var user = new(models.User)
	err = DB.Get(user, "select * from users where username=$1", username)
	return err
}

func GetUserById(id string) {
}

func InsertUser(user *models.User) error {
	_, err := DB.NamedExec("insert into users(id, username, password) values( :id, :username, :password)", user)
	return err
}

func RegisterNewUser(user *models.User, integration *models.Integration, user_integration *models.UserIntegration, credentials *[]models.Credential) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	//Insert user
	tx.NamedExec("insert into users(id, username, password) values(:id, :username, :password)", user)
	//Insert integration
	tx.NamedExec("insert into integrations(id, name) values(:id, :name)", integration)
	//Link user to integration
	tx.NamedExec("insert into user_integrations(id, integration, owner) values(:id, :integration, :owner)", user_integration)
	//Insert supabase creds
	tx.NamedExec(
		`
      insert into integration_credentials(id, integration, key, value)
      values(:id, :integration, :key, :value)
    `,
		&credentials,
	)
	err = tx.Commit()
	return err
}
