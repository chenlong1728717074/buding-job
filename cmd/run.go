package cmd

import (
	"buding-job/app"
	"errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "start",
	Short: "job start",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "cluster" {
			Error(cmd, args, errors.New("not currently supported"))
		}
		jobApp := app.NewBuDingJobApp()
		jobApp.Start()
	},
}

func Execute() {
	rootCmd.Execute()
}
