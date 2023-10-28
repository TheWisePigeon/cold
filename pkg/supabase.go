package pkg

import (
	"fmt"
	"net/http"
)

func CheckSupabaseCreds(project_url, bucket_name, api_key string) (ok bool, status_code int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/storage/v1/bucket/%s", project_url, bucket_name), nil)
	if err != nil {
		Logger.Error(err)
		return false, 500
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		Logger.Error(err)
		return false, 500
	}
	if response.StatusCode == 200 {
		return true, 200
	}
	if response.StatusCode == 400 {
		return false, 400
	}
	return false, response.StatusCode
}
