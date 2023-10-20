package handlers

import (
	"cold/models"
	"database/sql"
	"strings"
	"time"

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
		"Location":    "Settings",
		"GCP_Config":  gcp_config,
		"EmptyConfig": err == sql.ErrNoRows,
	}, "layout")
}

func SaveGCPSettings(db *sqlx.DB, c *fiber.Ctx) error {
	error_tag := "<p class='text-sm text-red-500'>Something went wrong. Please try again</p>"
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
		return c.SendString(error_tag)
	}
	last_updated := ""
	if payload.UploadedServiceAccountKey == "true" {
		file, err := c.FormFile("gcp_service_key")
		if !strings.HasSuffix(file.Filename, ".json") {
			error_tag = "<p class='text-sm text-red-500'>The service account key must be a JSON file</p>"
			return c.SendString(error_tag)
		}
		err = c.SaveFile(file, "./gcp_service_key.json")
		if err != nil {
			println(err.Error())
			return c.SendString("Something went wrong. Please try again")
		}
		last_updated = time.Now().Format("Monday 2 2006, 15:04")
		_, err = db.Exec(
			"update gcp_configs set last_updated_service_account=$1 where id=1",
			last_updated,
		)
	}
	//Because only one GC Storage setting is supported at the moment
	gcp_config := new(models.GCP_Config)
	gcp_config.Id = 1
	gcp_config.BucketName = payload.BucketName
	gcp_config.ProjectId = payload.ProjectId
	gcp_config.ServiceAccountKey = "./gcp_service_key.json"
	_, err = db.NamedExec(
		`
    insert or replace into gcp_configs(id, bucket_name, project_id, service_account_key, last_updated_service_account) 
    values(:id, :bucket_name, :project_id, :service_account_key, :last_updated_service_account)
    `,
		gcp_config,
	)
	if err != nil {
		println(err.Error())
		return c.SendString(error_tag)
	}
	return c.Redirect("/settings")
}
