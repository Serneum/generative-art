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
	JitterEnabled   bool
	OverlapEnabled  bool
}

type Dots struct {
	UserParams
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
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

	if !s.OverlapEnabled {
		s.InitialAlpha = 255
	}

	s.source = source
	s.dc = canvas
	return s
}

func (s *Dots) Update() {
	xOffset := (s.DestWidth - (int(s.DestWidth/s.Radius) * s.Radius)) / 2
	yOffset := (s.DestHeight - (int(s.DestHeight/s.Radius) * s.Radius)) / 2

	increment := s.Radius * 2
	if s.OverlapEnabled {
		increment = rand.Intn(s.Radius) + s.Radius
	}
	// i += s.Radius provides a cool overlapping effect. Maybe make a DotsOverlap
	for i := s.Radius + xOffset; i < s.sourceWidth; i += increment {
		for j := s.Radius + yOffset; j < s.sourceHeight; j += increment {
			r, g, b := rgb255(s.source.At(int(i), int(j)))

			destX := float64(i) * float64(s.DestWidth) / float64(s.sourceWidth)
			destY := float64(j) * float64(s.DestHeight) / float64(s.sourceHeight)
			if s.JitterEnabled {
				destX += float64(randRange(s.StrokeJitter))
				destY += float64(randRange(s.StrokeJitter))
			}

			s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
			s.dc.DrawCircle(destX, destY, float64(s.Radius))
			s.dc.FillPreserve()

			s.dc.Stroke()
			if s.OverlapEnabled {
				s.InitialAlpha = float64(rand.Intn(256))
				// s.InitialAlpha += s.AlphaIncrease
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
