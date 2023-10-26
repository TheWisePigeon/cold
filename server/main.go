package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load() // Dev only
	db, err := sqlx.Connect("sqlite3", "./.dev.cold.db")
	_ = db
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetReportCaller(true)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
			integration := chi.URLParam(r, "integration")
			if integration == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if integration != "supabase" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var payload = new(struct {
				Username               string            `json:"username"`
				Password               string            `json:"password"`
				IntegrationCredentials map[string]string `json:"integration_creds"`
			})
			err := json.NewDecoder(r.Body).Decode(payload)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf("%#v\n", payload)
			return

		})

		r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {

		})
	})

	fmt.Println("Server launched")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}
