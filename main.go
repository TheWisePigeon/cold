package main

import (
	"cold/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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
    return handlers.GetLoginPage(db, c)
	})

  app.Get("/logout", func(c *fiber.Ctx) error {
    return handlers.Logout(db, c)
  })

  app.Get("/home", func(c *fiber.Ctx) error {
    return handlers.GotoHomePage(db, c)
  })

	app.Get("/settings", func(c *fiber.Ctx) error {
    return handlers.GotoSettingsPage(db, c)
	})


  app.Route("/api", func(router fiber.Router) {
    router.Post("/login", func(c *fiber.Ctx) error {
      return handlers.Login(db, c)
    })
  })

	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Render("error", fiber.Map{})
	})

	println("Server launched")
	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
