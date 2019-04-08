package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "stops and starts deployments and statefulSets",
	Long:  `This command stops and starts kubernetes deployments and statefulSets`,
	Run: func(cmd *cobra.Command, args []string) {
		stopCmd.Run(cmd, args)
		startCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
