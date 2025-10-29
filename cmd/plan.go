package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	scopePath string
	outPath   string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Create a plan (stub)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("[plan] stub — scope: %s, out: %s\n", scopePath, outPath)
	},
}

func init() {
	planCmd.Flags().StringVar(&scopePath, "scope", "", "Path to scope.yaml")
	_ = planCmd.MarkFlagRequired("scope")

	planCmd.Flags().StringVar(&outPath, "out", "", "Output plan path")
	_ = planCmd.MarkFlagRequired("out")
	
	rootCmd.AddCommand(planCmd)
}