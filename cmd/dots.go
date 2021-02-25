package cmd

import (
	"generative-art/dots"
	"generative-art/util"

	"github.com/spf13/cobra"
)

var (
	fade    = "none"
	overlap = false
	radius  = 5
)

func init() {
	rootCmd.AddCommand(dotsCmd)

	dotsCmd.Flags().StringVarP(&file, "file", "l", "", "The path to a local file")
	dotsCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	dotsCmd.Flags().IntVarP(&width, "width", "", 2000, "Width of output")
	dotsCmd.Flags().IntVarP(&height, "height", "", 2000, "Height of output")

	dotsCmd.Flags().Float64VarP(&alpha, "alpha", "a", 0.1, "Starting alpha")
	dotsCmd.Flags().Float64VarP(&alphaIncrease, "alphaIncrease", "i", 0.06, "Increase of alpha per iteration")
	dotsCmd.Flags().StringVarP(&fade, "fade", "f", "none", "Direction to fade the image in")
	dotsCmd.Flags().IntVarP(&jitter, "jitter", "j", 0, "Jitter multiplier")
	dotsCmd.Flags().BoolVarP(&overlap, "overlap", "o", false, "Allow the dots to overlap each other")
	dotsCmd.Flags().IntVarP(&radius, "radius", "r", 25, "The radius of the dots when creating the image")
	dotsCmd.Flags().Float64VarP(&ratio, "ratio", "t", 0.75, "Starting path size as a ratio of image width")
	dotsCmd.Flags().Float64VarP(&reduction, "reduction", "d", 0.002, "Reduction per iteration")
}

var dotsCmd = &cobra.Command{
	Use:   "dots",
	Short: "Create a dot-based version of the image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		img := util.LoadImage(file, url, width, height)
		dots := dots.NewDots(img, dots.UserParams{
			DestWidth:       width,
			DestHeight:      height,
			StrokeRatio:     ratio,
			StrokeReduction: reduction,
			InitialAlpha:    alpha,
			AlphaIncrease:   alphaIncrease,
			Radius:          radius,
			Overlap:         overlap,
			Fade:            fade,
			StrokeJitter:    jitter,
		})

		dots.Update()
		util.SaveImage(dots.Output())
	},
}
