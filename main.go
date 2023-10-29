package main

import (
	"cold/handlers"
	"cold/pkg"
	"cold/repositories"
	"fmt"
	"net/http"
	"os"
  "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)


//go:embed views/*
var views embed.FS

//go:embed static/output.css
var static embed.FS

func main() {
	godotenv.Load() // Dev only
	pkg.InitLogger()
  pkg.Views = views
	err := repositories.ConnectToDB()
	if err != nil {
		pkg.Logger.Fatal(err)
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

  r.Get("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
    css, err := static.ReadFile("static/output.css")
    if err!= nil {
      pkg.Logger.Error(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    w.Header().Set("Content-Type", "text/css")
    w.Write(css)
    return
  })
	handlers.RegisterHandlers(r)

	fmt.Println("Server launched")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
