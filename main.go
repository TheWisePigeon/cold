package main

import (
	"cold/handlers"
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed views/*
var viewsfs embed.FS

func main() {
	db, err := sqlx.Connect("sqlite3", "./cold.db")
	_ = db
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")
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

		router.Get("/gcp-key-upload", func(c *fiber.Ctx) error {
			input_tag := `
        <input type="text" hidden name="gcp_service_key_uploaded" value="true">
        <div class="mt-2 flex gap-10">
          <input required class="block w-full rounded-md border-0 bg-white p-2 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" id="file_input" name="gcp_service_key" type="file" />
          <button class="rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50" hx-trigger="click" hx-target="#gcp_key_upload" hx-get="/api/cancel-gcp-key-upload">Cancel</button>
        </div>
      `
			return c.SendString(input_tag)
		})

    router.Get("/cancel-gcp-key-upload", func(c *fiber.Ctx) error {
      return c.SendString(
        `
          <input type="text" hidden name="gcp_service_key_uploaded" value="false">
          <div class="mt-2 flex items-center gap-x-3">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
              stroke="currentColor" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
            </svg>
            <button type="button" hx-trigger="click" hx-get="/api/gcp-key-upload" hx-target="#gcp_key_upload"
              class="rounded-md bg-white px-2.5 py-1.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
            >
              <h1>Upload new</h1>
            </button>
          </div>
        `,
      )
    })

		router.Route("/settings", func(router fiber.Router) {
			router.Post("/gcp", func(c *fiber.Ctx) error {
				return handlers.SaveGCPSettings(db, c)
			})
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
