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
    pkg.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("%#v\n", payload)
	return

}
