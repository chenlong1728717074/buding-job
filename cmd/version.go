package cmd

import (
	"buding-job/common/constant"
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constant.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
