package cmd

import (
	"github.com/spf13/cobra"

	"cheaptrick/internal/fixture"
)

var count int
var outputDir string
var promptType string

var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "Generate predefined text and tool call fixtures with a MANIFEST.md",
	Run: func(cmd *cobra.Command, args []string) {
		fixture.GenerateFromPrompts(outputDir, promptType, count)
	},
}

func init() {
	fixturesCmd.Flags().IntVarP(&count, `count`, `c`, 10, `Number of fixtures to generate`)
	fixturesCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "fixtures", "Directory to output fixtures")
	fixturesCmd.Flags().StringVarP(&promptType, "prompt-type", "p", "text", "Type of prompts to generate (text, tool-call)")

	rootCmd.AddCommand(fixturesCmd)
}
