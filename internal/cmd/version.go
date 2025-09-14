package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show fastauto version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("fastauto %s\n", versionString())
        if BuildInfo.Commit != "" {
            fmt.Printf("commit: %s\n", BuildInfo.Commit)
        }
        if BuildInfo.Date != "" {
            fmt.Printf("built:  %s\n", BuildInfo.Date)
        }
        if BuildInfo.Go != "" {
            fmt.Printf("go:     %s\n", BuildInfo.Go)
        }
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

