package repositories

import "cold/models"

func CreateSession(session_id string, user_id string) error {
	_, err := DB.NamedExec("insert into sessions(id, session_user) values(:id, :session_user)", &models.Session{Id: session_id, SessionUser: user_id})
	return err
}
