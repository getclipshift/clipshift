package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print version",
	Long:    `Print version number`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		println(Version)
	},
}
