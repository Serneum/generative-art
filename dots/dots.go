package dots

import (
	"generative-art/util"
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
	Overlap         bool
	Fade            string
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
	currentAlpha      float64
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
	canvas.Stroke()

	s.setIncrementSize()
	s.setAlphaSettings()

	s.source = source
	s.dc = canvas
	return s
}

func (s *Dots) Update() {
	xOffset := float64(s.DestWidth-(int(s.DestWidth/s.increment)*s.increment)) / 2
	yOffset := float64(s.DestHeight-(int(s.DestHeight/s.increment)*s.increment)) / 2

	xIter, yIter := s.iterationCounts()

	for i := s.Radius; i < s.sourceWidth; i += s.increment {
		for j := s.Radius; j < s.sourceHeight; j += s.increment {
			r, g, b := util.Rgb255(s.source.At(int(i), int(j)))

			destX := float64(i) * float64(s.DestWidth) / float64(s.sourceWidth)
			destX += float64(xOffset)

			destY := float64(j) * float64(s.DestHeight) / float64(s.sourceHeight)
			destY += float64(yOffset)

			if s.StrokeJitter != 0 {
				destX += float64(util.RandRange(s.StrokeJitter))
				destY += float64(util.RandRange(s.StrokeJitter))
			}

			if s.Fade == "random" {
				s.currentAlpha = float64(rand.Intn(256))
			} else if s.Fade != "none" {
				s.setCurrentAlpha(i, j, xIter, yIter)
			}

			s.dc.SetRGBA255(r, g, b, int(s.currentAlpha))
			s.dc.DrawCircle(destX, destY, float64(s.Radius))
			s.dc.Fill()
			s.dc.Stroke()
		}
	}
}

func (s *Dots) Output() image.Image {
	return s.dc.Image()
}

func (s *Dots) iterationCounts() (xIter, yIter float64) {
	xIter = float64(s.sourceWidth) / float64(s.increment)
	yIter = float64(s.sourceHeight) / float64(s.increment)
	return xIter, yIter
}

func (s *Dots) setIncrementSize() {
	s.increment = s.Radius * 2
	if s.Overlap {
		s.increment = rand.Intn(s.Radius) + s.Radius
	}
}

func (s *Dots) setAlphaSettings() {
	if s.Fade == "none" || s.Fade == "" {
		s.InitialAlpha = 255
	} else if s.AlphaIncrease == 0 {
		totalIncrease := 255 - s.InitialAlpha
		xIter, yIter := s.iterationCounts()
		s.AlphaIncrease = float64(totalIncrease) / (float64(xIter * yIter))
	}

	if s.Fade == "right" || s.Fade == "bottom" {
		s.InitialAlpha = 255 - s.InitialAlpha
		s.AlphaIncrease *= -1
	}

	s.currentAlpha = s.InitialAlpha
}

func (s *Dots) setCurrentAlpha(iterCol, iterRow int, xIter, yIter float64) {
	row, col, iter := iterRow, iterCol, yIter
	if s.Fade == "top" || s.Fade == "bottom" {
		row, col, iter = iterCol, iterRow, xIter
	}

	gridX := float64(col) / float64(s.Radius*2)
	gridY := float64(row) / float64(s.Radius*2)
	iterationsBeforePos := gridX*iter + gridY

	s.currentAlpha = s.InitialAlpha + iterationsBeforePos*s.AlphaIncrease
	if s.currentAlpha > 255 {
		s.currentAlpha = 255
	} else if s.currentAlpha < 0 {
		s.currentAlpha = 0
	}
}
