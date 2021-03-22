package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mp3tag",
	Long:  "Print the version number of mp3tag",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("no parameters for this command")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}
