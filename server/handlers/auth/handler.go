package auth

import (
	"cold/pkg"
	"encoding/json"
	"fmt"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	integration := r.URL.Query().Get("integration")
	if integration == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if integration == "supabase" {
		var payload = new(struct {
			Username            string `json:"username"`
			Password            string `json:"password"`
			SupabaseCredentials struct {
				ProjectUrl string `json:"project_url"`
				BucketName string `json:"bucket_name"`
				ApiKey     string `json:"api_key"`
				Folder     string `json:"folder"`
			} `json:"creds"`
		})
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			pkg.Logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ok, status_code := pkg.CheckSupabaseCreds(
			payload.SupabaseCredentials.ProjectUrl,
			payload.SupabaseCredentials.BucketName,
			payload.SupabaseCredentials.ApiKey,
		)
		if !ok {
			switch status_code {
			case 400:
				w.WriteHeader(http.StatusBadRequest)
				return
			case 404:
				w.WriteHeader(http.StatusNotFound)
				return
			case 500:
				w.WriteHeader(http.StatusInternalServerError)
				return
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}
