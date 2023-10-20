package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type LoginPayload struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
  FirstTimeLoggingIn bool `db:"first_time_logging_in"`
}

type Session struct {
	Id   string `json:"id" db:"id"`
	User string `json:"user" db:"user"`
}

func main() {
	db, err := sqlx.Connect("sqlite3", "./cold.db")
	_ = db
	engine := html.New("/home/thewisepigeon/code/cold/views", ".html")
	if err != nil {
		panic(err)
	}
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	app.Post("/api/login", func(c *fiber.Ctx) error {
		payload := new(LoginPayload)
		err := c.BodyParser(payload)
		if err != nil {
			return c.Render("login", fiber.Map{
				"BadRequest": true,
			})
		}
		user := new(User)
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
	})

  app.Get("/logout", func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id", "")
		if session_id == "" {
			return c.Redirect("login")
		}
		session := new(Session)
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
    if err!=nil{
      return c.Redirect("/error")
    }
    return c.Redirect("/login")
  })

	app.Get("/home", func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id", "")
		if session_id == "" {
			return c.Redirect("login")
		}
		session := new(Session)
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
    current_user := new(User)
    err = db.Get(
      current_user,
      "select * from users where username=$1",
      session.User,
    )
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Redirect("/login")
			}
			return c.Redirect("/error")
		}
    if current_user.FirstTimeLoggingIn {
      c.Redirect("/settings")
    }
		return c.Render("home", fiber.Map{
			"Username": session.User,
      "Location": "Home",
		}, "layout")
	})

  app.Get("/settings", func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id", "")
		if session_id == "" {
			return c.Redirect("login")
		}
		session := new(Session)
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
    return c.Render("settings", fiber.Map{
      "Location":"Settings",
    }, "layout")
  })

	println("Server launched")
	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
