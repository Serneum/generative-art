package cmd

import (
	"generative-art/sketch"
	"generative-art/util"
	"math/rand"

	"github.com/spf13/cobra"
)

var (
	iterations   = 0
	minEdgeCount = 0
	maxEdgeCount = 0
)

func init() {
	rootCmd.AddCommand(layerCmd)

	layerCmd.Flags().StringVarP(&file, "file", "l", "", "The path to a local file")
	layerCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	layerCmd.Flags().IntVarP(&width, "width", "", 2000, "Width of output")
	layerCmd.Flags().IntVarP(&height, "height", "", 2000, "Height of output")

	layerCmd.Flags().IntVarP(&iterations, "iterations", "i", 0, "Number of iterations")
	layerCmd.Flags().Float64VarP(&reduction, "reduction", "", 0.002, "Reduction per iteration")
	layerCmd.Flags().Float64VarP(&ratio, "ratio", "", 0.75, "Starting path size as a ratio of image width")
	layerCmd.Flags().Float64VarP(&alpha, "alpha", "a", 0.1, "Starting alpha")
	layerCmd.Flags().Float64VarP(&alphaIncrease, "alphaIncrease", "", 0.06, "Increase of alpha per iteration")
	layerCmd.Flags().IntVarP(&minEdgeCount, "minimumEdges", "", 0, "Minimum number of edges of path")
	layerCmd.Flags().IntVarP(&maxEdgeCount, "maximumEdges", "", 0, "Maximum number of edges of path")
	layerCmd.Flags().Float64VarP(&jitter, "jitter", "", 0, "Jitter multiplier")
	layerCmd.Flags().Float64VarP(&inversionThreshold, "inversion", "", 0.05, "Size at which to invert the color")
}

var layerCmd = &cobra.Command{
	Use:   "layer",
	Short: "Create sketches in the style of Preslav Rachev",
	Long:  `Create a sketch of overlapping shapes with various drawing options`,
	Run: func(cmd *cobra.Command, args []string) {
		img := util.LoadImage(file, url, width, height)

		if maxEdgeCount < minEdgeCount {
			maxEdgeCount = minEdgeCount
		}

		sketch := sketch.NewSketch(img, sketch.UserParams{
			DestWidth:                width,
			DestHeight:               height,
			StrokeRatio:              ratio,
			StrokeReduction:          reduction,
			StrokeInversionThreshold: inversionThreshold,
			StrokeJitter:             int(jitter * float64(width)),
			InitialAlpha:             alpha,
			AlphaIncrease:            alphaIncrease,
			MinEdgeCount:             minEdgeCount,
			MaxEdgeCount:             maxEdgeCount,
		})

		if iterations == 0 {
			iterations = rand.Intn(5001)
		}

		for i := 0; i < iterations; i++ {
			sketch.Update()
		}

		util.SaveImage(sketch.Output())
	},
}
