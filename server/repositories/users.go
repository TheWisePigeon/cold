package repositories

import "cold/models"

func GetUserByName(username string) error {
  var user = new(models.User)
  err = DB.Get(user, "select * from users where username=$1", username)
  return err
}

func GetUserById(id string){
}

func InsertUser( user *models.User) error {
  _, err := DB.NamedExec("insert into users(id, username, password) values( :id, :username, :password)", user)
  return err
}
