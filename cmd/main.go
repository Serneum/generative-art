package main

import (
	"fmt"
	// "generative-art/sketch"
	"generative-art/dots"
	"image"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
)

func main() {
	rand.Seed(time.Now().Unix())

	// img, err := gg.LoadImage("source.jpg")
	img, err := loadRandomUnsplashImage(2000, 2000)
	if err != nil {
		log.Panicln(err)
	}

	destWidth := 2000
	destHeight := 2000
	radius := rand.Intn(51) + 25
	// minEdgeCount := rand.Intn(3) + 3
	// maxEdgeCount := rand.Intn(3) + minEdgeCount + 1
	sketch := dots.NewDots(img, dots.UserParams{
		DestWidth:       destWidth,
		DestHeight:      destHeight,
		StrokeRatio:     0.75,
		StrokeReduction: 0.002,
		InitialAlpha:    0.1,
		Radius:          radius,
		Overlap:         true,
		Fade:            true,
	})

	// sketch := sketch.NewSketch(img, sketch.UserParams{
	// 	DestWidth:                destWidth,
	// 	DestHeight:               destHeight,
	// 	StrokeRatio:              0.75,
	// 	StrokeReduction:          0.002,
	// 	StrokeInversionThreshold: 0.05,
	// 	StrokeJitter:             int(0.1 * float64(destWidth)),
	// 	InitialAlpha:             0.1,
	// 	AlphaIncrease:            0.06,
	// 	MinEdgeCount:             minEdgeCount,
	// 	MaxEdgeCount:             maxEdgeCount,
	// })

	// totalCycleCount := rand.Intn(5001)
	// for i := 0; i < totalCycleCount; i++ {
	// 	sketch.Update()
	// }
	sketch.Update()

	outputImgName := fmt.Sprintf("../output/%s.png", uuid.New())
	fmt.Printf("Output: %s\n", outputImgName)
	gg.SavePNG(outputImgName, sketch.Output())
}

func loadRandomUnsplashImage(width, height int) (image.Image, error) {
	url := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Image: %s\n", res.Request.URL)
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}
