package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var IsVersion bool
var MatchBy []string
var CaseSensitive bool

var rootCmd = &cobra.Command{
	Use:   "dedupthingy",
	Short: "Music dedup thingy",
	Run: func(cmd *cobra.Command, args []string) {
		if IsVersion {
			runVersion()
		} else {
			if len(args) > 0 {
				if err := runDeduplicate(args); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

			} else {
				cmd.Root().Help()
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&IsVersion, "version", "v", false, "show version")
	rootCmd.PersistentFlags().StringSliceVarP(&MatchBy, "match", "m", []string{"title", "artist"}, "Specify what tags have to match. Available: title,artist,album")
	rootCmd.PersistentFlags().BoolVarP(&CaseSensitive, "sensitive", "s", false, "Case sensitive")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
