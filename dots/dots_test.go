package dots

import "testing"

func TestIterationCounts(t *testing.T) {
	tables := []struct {
		dots  Dots
		xIter float64
		yIter float64
	}{
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1}, float64(10), float64(10)},
		{Dots{sourceWidth: 2000, sourceHeight: 2000, increment: 50}, float64(40), float64(40)},
		{Dots{sourceWidth: 786, sourceHeight: 960, increment: 50}, float64(15.72), float64(19.2)},
	}

	for _, table := range tables {
		xIter, yIter := table.dots.iterationCounts()
		if xIter != table.xIter || yIter != table.yIter {
			t.Errorf("Incorrect iteration counts for width: %d, height: %d, increment size: %d.\n"+
				"Got xIter: %f, yIter: %f\n"+
				"Expected xIter: %f, yIter: %f\n",
				table.dots.sourceWidth, table.dots.sourceHeight, table.dots.increment,
				xIter, yIter, table.xIter, table.yIter)
		}
	}
}

func TestSetIncrementSize(t *testing.T) {
	tables := []struct {
		dots Dots
	}{
		{Dots{UserParams: UserParams{Radius: 1, Overlap: false}}},
		{Dots{UserParams: UserParams{Radius: 5, Overlap: false}}},
		{Dots{UserParams: UserParams{Radius: 25, Overlap: false}}},
		{Dots{UserParams: UserParams{Radius: 1, Overlap: true}}},
		{Dots{UserParams: UserParams{Radius: 5, Overlap: true}}},
		{Dots{UserParams: UserParams{Radius: 25, Overlap: true}}},
	}

	for _, table := range tables {
		table.dots.setIncrementSize()
		if table.dots.Overlap {
			if table.dots.increment < 0 || table.dots.increment > 2*table.dots.Radius {
				// TODO: Mock out rand.Intn, update this spec
				t.Errorf("Incorrect increment size returned for radius: %d, overlap: %v\n"+
					"Got %d\n"+
					"Expected %d\n",
					table.dots.Radius, table.dots.Overlap, table.dots.increment, 2*table.dots.Radius)
			}
		} else if table.dots.increment != 2*table.dots.Radius {
			t.Errorf("Incorrect increment size returned for radius: %d, overlap: %v\n"+
				"Got %d\n"+
				"Expected %d\n",
				table.dots.Radius, table.dots.Overlap, table.dots.increment, 2*table.dots.Radius)
		}
	}
}

func TestSetAlphaSettings(t *testing.T) {
	tables := []struct {
		dots          Dots
		initialAlpha  float64
		alphaIncrease float64
	}{
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "none", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 255, 0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 255, 0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "left", AlphaIncrease: float64(0), InitialAlpha: 0.1}}, 0.1, 2.549},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "left", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 0.1, 0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "top", AlphaIncrease: float64(0), InitialAlpha: 0.1}}, 0.1, 2.549},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "top", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 0.1, 0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "right", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 254.9, -0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "right", AlphaIncrease: float64(0), InitialAlpha: 0.1}}, 254.9, -2.549},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "bottom", AlphaIncrease: 0.06, InitialAlpha: 0.1}}, 254.9, -0.06},
		{Dots{sourceWidth: 10, sourceHeight: 10, increment: 1, UserParams: UserParams{Fade: "bottom", AlphaIncrease: float64(0), InitialAlpha: 0.1}}, 254.9, -2.549},
	}

	for _, table := range tables {
		table.dots.setAlphaSettings()
		if table.initialAlpha != table.dots.InitialAlpha || table.alphaIncrease != table.dots.AlphaIncrease {
			t.Errorf("Incorrect alpha values for fade: '%s', initial alpha: %f, alpha increase: %f\n"+
				"Got initial alpha: %f, alpha increase: %f\n"+
				"Expected initial alpha: %f, alpha increase: %f\n",
				table.dots.Fade, table.dots.InitialAlpha, table.dots.AlphaIncrease,
				table.initialAlpha, table.alphaIncrease, table.initialAlpha, table.alphaIncrease)
		}
	}
}
