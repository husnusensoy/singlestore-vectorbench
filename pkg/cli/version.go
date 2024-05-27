package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Vector Benchmark (vectbench) tag: %s (%s), build time %s\n", tag, sha1ver, buildTime)
	},
}
