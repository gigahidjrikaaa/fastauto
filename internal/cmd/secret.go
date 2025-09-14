package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "fastauto/internal/config"
)

var secretCmd = &cobra.Command{
    Use:   "secret",
    Short: "Manage secrets",
}

var secretRotateCmd = &cobra.Command{
    Use:   "rotate",
    Short: "Rotate the webhook HMAC secret",
    RunE: func(cmd *cobra.Command, args []string) error {
        g, _ := config.LoadGlobalConfig()
        if err := g.GenerateWebhookSecret(); err != nil { return err }
        if err := config.SaveGlobalConfig(g); err != nil { return err }
        fmt.Println("Rotated global webhook secret.")
        return nil
    },
}

func init() {
    secretCmd.AddCommand(secretRotateCmd)
}

