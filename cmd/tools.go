package cmd

import (
	"cheaptrick/internal/shell/generator"

	"github.com/spf13/cobra"
)

var toolsOutputDir string

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Generate 20 sample canned tool response files for use with `cheaptrick shell --tool-responses`",
	Long: `Generates a directory of canned tool response files that simulate common
external tool integrations. These files are used by the shell's tool-call
loop to automatically respond to FunctionCall parts without manual input.

The generated tools cover weather, search, email, databases, calendars,
file I/O, translation, ticketing, messaging, and more — providing a
realistic starting point for agent development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.GenerateTools(toolsOutputDir)
	},
}

func init() {
	toolsCmd.Flags().StringVarP(&toolsOutputDir, "output-dir", "o", "./mock_tools", "Output directory for generated tool files")
	rootCmd.AddCommand(toolsCmd)
}
