package auth

import (
	"cold/models"
	"cold/pkg"
	"cold/repositories"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
  "html/template"
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
		err = repositories.GetUserByName(payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				new_userid := uuid.NewString()
				err, hashed_pwd := pkg.Hash(payload.Password)
				if err != nil {
					pkg.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				integration_id := uuid.NewString()
				supabase_creds := []models.Credential{
					{
						Id:          uuid.NewString(),
						Integration: integration_id,
						Key:         "project_url",
						Value:       payload.SupabaseCredentials.ProjectUrl,
					},
					{
						Id:          uuid.NewString(),
						Integration: integration_id,
						Key:         "api_key",
						Value:       payload.SupabaseCredentials.ApiKey,
					},
					{
						Id:          uuid.NewString(),
						Integration: integration_id,
						Key:         "bucket_name",
						Value:       payload.SupabaseCredentials.BucketName,
					},
					{
						Id:          uuid.NewString(),
						Integration: integration_id,
						Key:         "folder",
						Value:       payload.SupabaseCredentials.Folder,
					},
				}
				err = repositories.RegisterNewUser(
					&models.User{
						Id:       new_userid,
						Username: payload.Username,
						Password: hashed_pwd,
					},
					&models.Integration{
						Id:   integration_id,
						Name: "supabase",
					},
					&models.UserIntegration{
						Id:          uuid.NewString(),
						Integration: integration_id,
						Owner:       new_userid,
					},
					supabase_creds,
				)
				if err != nil {
					pkg.Logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				new_sessionid := uuid.NewString()
				err = repositories.CreateSession(new_sessionid, new_userid)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Add("redirect-to", "login")
					return
				}
				cookie := http.Cookie{
					Name:     "session_id",
					Value:    new_sessionid,
					HttpOnly: true,
					Path:     "/",
				}
				http.SetCookie(w, &cookie)
				w.WriteHeader(http.StatusCreated)
				return
			}
			pkg.Logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func GetAuthPage(w http.ResponseWriter, r *http.Request){
  t, _ := template.New("auth.gohtml").ParseFS(pkg.Views, "views/auth.gohtml")
  t.Execute(w, nil)
}
