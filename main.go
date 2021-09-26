package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

const (
	BASE_URL = "https://www.bing.com"
	URL      = BASE_URL + "/HPImageArchive.aspx?format=js&idx=0&n=8"
)

func main() {
	cwd, _ := os.Getwd()
	response, err := http.Get(URL)
	check(err)

	body, _ := ioutil.ReadAll(response.Body)
	var data map[string]interface{}
	errr := json.Unmarshal(body, &data)
	check(errr)

	images := parseResponse(data)
	wg := new(sync.WaitGroup)

	for _, image := range images {
		wg.Add(1)
		go downloadImage(wg, image, path.Join(cwd, image.Slug+".jpeg"))
	}

	wg.Wait()
}
