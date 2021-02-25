package dots

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

type UserParams struct {
	DestWidth       int
	DestHeight      int
	StrokeRatio     float64
	StrokeReduction float64
	StrokeJitter    int
	InitialAlpha    float64
	AlphaIncrease   float64
	Radius          int
	Jitter          bool
	Overlap         bool
	Fade            bool
	FadeDirection   string
}

type Dots struct {
	UserParams
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
	increment         int
}

func NewDots(source image.Image, userParams UserParams) *Dots {
	s := &Dots{UserParams: userParams}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()

	s.increment = s.Radius * 2
	if s.Overlap {
		s.increment = rand.Intn(s.Radius) + s.Radius
	}

	if !s.Fade {
		s.InitialAlpha = 255
	} else {
		if s.AlphaIncrease == 0 {
			totalIncrease := 255 - s.InitialAlpha
			s.AlphaIncrease = float64(totalIncrease) / float64((s.sourceWidth/s.increment)*(s.sourceHeight/s.increment))
		}

		if s.FadeDirection == "" {
			s.FadeDirection = "left"
		}
	}

	if s.FadeDirection == "right" {
		s.InitialAlpha = 255 - s.InitialAlpha
		s.AlphaIncrease *= -1
	}

	if s.Jitter && s.StrokeJitter == 0 {
		s.StrokeJitter = int(0.1 * float64(s.DestWidth))
	}

	s.source = source
	s.dc = canvas
	return s
}

func (s *Dots) Update() {
	xOffset := float64(s.DestWidth-(int(s.DestWidth/s.increment)*s.increment)) / 2
	yOffset := float64(s.DestHeight-(int(s.DestHeight/s.increment)*s.increment)) / 2

	// i += s.Radius provides a cool overlapping effect. Maybe make a DotsOverlap
	for i := s.Radius; i < s.sourceWidth; i += s.increment {
		for j := s.Radius; j < s.sourceHeight; j += s.increment {
			r, g, b := rgb255(s.source.At(int(i), int(j)))

			destX := float64(i) * float64(s.DestWidth) / float64(s.sourceWidth)
			destY := float64(j) * float64(s.DestHeight) / float64(s.sourceHeight)
			destX += float64(xOffset)
			destY += float64(yOffset)
			if s.Jitter {
				destX += float64(randRange(s.StrokeJitter))
				destY += float64(randRange(s.StrokeJitter))
			}

			s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
			s.dc.DrawCircle(destX, destY, float64(s.Radius))
			s.dc.FillPreserve()

			s.dc.Stroke()
			if s.FadeDirection == "random" {
				s.InitialAlpha = float64(rand.Intn(256))
			} else {
				s.InitialAlpha += s.AlphaIncrease
			}
		}
	}
}

func (s *Dots) Output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}
