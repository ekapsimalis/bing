package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gosimple/slug"
)

type Image struct {
	Startdate time.Time
	Url       string
	Name      string
	Slug      string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseResponse(data map[string]interface{}) []*Image {
	var returnImages []*Image
	var images []interface{} = data["images"].([]interface{})

	for _, image := range images {
		var parsed_image map[string]interface{} = image.(map[string]interface{})
		date, _ := time.Parse("20060102", parsed_image["startdate"].(string))
		name := parsed_image["copyright"].(string)
		url := BASE_URL + parsed_image["url"].(string)
		slug := slug.Make(name)

		returnImages = append(returnImages, &Image{date, url, name, slug})
	}

	return returnImages
}

func downloadImage(wg *sync.WaitGroup, image *Image, path string) {
	defer wg.Done()

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	response, _ := http.Get(image.Url)
	io.Copy(file, response.Body)
	fmt.Println("Downloaded " + image.Name)
}
