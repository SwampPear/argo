package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	verbose int
)

var rootCmd = &cobra.Command{
	Use:   "argo",
	Short: "Vuln discovery",
	Long:  "Vuln discovery",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "Increase verbosity (-v, -vv)")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}