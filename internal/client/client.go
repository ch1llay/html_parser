package client

import (
	"io/ioutil"
	"net/http"
)

func GetHtml(url string) string {
	var client = &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	response, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return string(body)

}
