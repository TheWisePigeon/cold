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
    return c.Redirect("/home")
	})

	println("Server launched")
	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
