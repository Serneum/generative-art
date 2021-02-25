package util

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
)

func LoadImage(file, url string, width, height int) image.Image {
	var img image.Image
	var err error
	if file != "" {
		img, err = gg.LoadImage(file)
	} else {
		img, err = loadRandomUnsplashImage(url, width, height)
	}

	if err != nil {
		log.Panicln(err)
	}
	return img
}

func SaveImage(image image.Image) {
	outputImgName := fmt.Sprintf("output/%s.png", uuid.New())
	fmt.Printf("Output: %s\n", outputImgName)
	gg.SavePNG(outputImgName, image)
}

func loadRandomUnsplashImage(url string, width, height int) (image.Image, error) {
	if url == "" {
		url = fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Image: %s\n", res.Request.URL)
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}

func Rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func RandRange(max int) int {
	return -max + rand.Intn(2*max)
}
