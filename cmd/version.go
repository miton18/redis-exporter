package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	githash = "HEAD"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("redis-exporter %s (%s)\n", version, githash)
	},
}
