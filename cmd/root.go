package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context, version string) {
	rootCmd := &cobra.Command{
		Use:   "gerard",
		Short: "gerard your butler services....",
		Long:  `gerard your butler services....`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version:      version,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(aaExecute())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
