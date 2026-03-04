package cmd

import (
	"github.com/spf13/cobra"

	"cheaptrick/internal/fixture"
)

var outputDir string

var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "Generate 30 predefined text and tool call fixtures with a MANIFEST.md",
	Run: func(cmd *cobra.Command, args []string) {
		fixture.GenerateFromPrompts(outputDir)
	},
}

func init() {
	fixturesCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "fixtures", "Directory to output fixtures")

	rootCmd.AddCommand(fixturesCmd)
}
