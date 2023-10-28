package main

import (
	"cold/handlers"
	"cold/pkg"
	"cold/repositories"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // Dev only
	pkg.InitLogger()
	err := repositories.ConnectToDB()
	if err != nil {
		pkg.Logger.Fatal(err)
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handlers.RegisterHandlers(r)

	fmt.Println("Server launched")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
