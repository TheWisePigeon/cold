package handlers

import (
	"cold/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func GetLoginPage(db *sqlx.DB, c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func Login(db *sqlx.DB, c *fiber.Ctx) error {
	payload := new(models.LoginPayload)
	err := c.BodyParser(payload)
	if err != nil {
		return c.Render("login", fiber.Map{
			"BadRequest": true,
		})
	}
	user := new(models.User)
	err = db.Get(
		user,
		"select * from users where username=$1",
		payload.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Render("login", fiber.Map{
				"BadRequest": true,
			})
		}
		return c.Render("login", fiber.Map{
			"InternalError": true,
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Render("login", fiber.Map{
			"BadRequest": true,
		})
	}
	session_id := uuid.NewString()
	_, err = db.Exec(
		"insert into sessions(id, user) values($1, $2)",
		session_id,
		user.Username,
	)
	if err != nil {
		return c.Render("login", fiber.Map{
			"InternalError": true,
		})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = session_id
	cookie.Path = "/"
	c.Cookie(cookie)
	c.Set("hx-redirect", "/home")
	return c.Render("/login", fiber.Map{})
}

func Logout(db *sqlx.DB, c *fiber.Ctx) error {
	session_id := c.Cookies("session_id", "")
	if session_id == "" {
		return c.Redirect("login")
	}
	session := new(models.Session)
	err := db.Get(
		session,
		"select * from sessions where id=$1",
		session_id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect("/login")
		}
		return c.Redirect("/error")
	}
	_, err = db.Exec("delete from sessions where id=$1", session.Id)
	if err != nil {
		return c.Redirect("/error")
	}
	return c.Redirect("/login")
}

func GotoHomePage(db *sqlx.DB, c *fiber.Ctx) error {
	session_id := c.Cookies("session_id", "")
	if session_id == "" {
		return c.Redirect("login")
	}
	session := new(models.Session)
	err := db.Get(
		session,
		"select * from sessions where id=$1",
		session_id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect("/login")
		}
		return c.Redirect("/error")
	}
	//As only GC Storage is supported at the moment
	gcp_config := new(models.GCP_Config)
	err = db.Get(
		gcp_config,
		"select * from gcp_configs where id=1",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect("/settings")
		}
		return c.Redirect("/error")
	}
	return c.Render("home", fiber.Map{
		"Username": session.User,
		"Location": "Home",
	}, "layout")
}

func GotoSettingsPage(db *sqlx.DB, c *fiber.Ctx) error {
	session_id := c.Cookies("session_id", "")
	if session_id == "" {
		return c.Redirect("login")
	}
	session := new(models.Session)
	err := db.Get(
		session,
		"select * from sessions where id=$1",
		session_id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect("/login")
		}
		return c.Redirect("/error")
	}
	//Fetch GCP_config
	gcp_config := new(models.GCP_Config)
	err = db.Get(
		gcp_config,
		"select * from gcp_configs where id=1",
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return c.Redirect("/error")
		}
	}
	return c.Render("settings", fiber.Map{
		"Location":   "Settings",
		"GCP_Config": gcp_config,
	}, "layout")
}

func SaveGCPSettings(db *sqlx.DB, c *fiber.Ctx) error {
	session_id := c.Cookies("session_id", "")
	if session_id == "" {
		return c.Redirect("login")
	}
	session := new(models.Session)
	err := db.Get(
		session,
		"select * from sessions where id=$1",
		session_id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect("/login")
		}
		return c.Redirect("/error")
	}
	payload := new(models.GCPPayload)
  err = c.BodyParser(payload)
  if err != nil {
    return c.Render("settings", fiber.Map{
      "Location":   "Settings",
      "GCP_Config": payload,
      "BadRequest": true,
    })
  }
	file, err := c.FormFile("gcp_service_key")
	if err != nil {
		println(err.Error())
		return c.Render("settings", fiber.Map{
			"Location":   "Settings",
			"GCP_Config": payload,
			"BadRequest": true,
		})
	}
  println(file.Filename)
	return c.Render("settings", fiber.Map{
		"Location":   "Settings",
		"GCP_Config": payload,
	})
}
