package cli

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database related commands",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
