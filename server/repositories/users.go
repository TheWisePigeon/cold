package repositories

import "cold/models"

func GetUserByName(username string){
}

func GetUserById(id string){
}

func InsertUser( user *models.User) error {
  _, err := DB.NamedExec("insert into users(id, username, password) values( :id, :username, :password)", user)
  return err
}
