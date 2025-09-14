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

var secretShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Print the current webhook HMAC secret",
    RunE: func(cmd *cobra.Command, args []string) error {
        g, _ := config.LoadGlobalConfig()
        if g.WebhookSecret == "" {
            fmt.Println("No secret set. Generate one with: fastauto secret rotate")
            fmt.Printf("Config: %s\n", config.GlobalPath())
            return nil
        }
        fmt.Printf("Secret: %s\n", g.WebhookSecret)
        fmt.Printf("Config: %s\n", config.GlobalPath())
        return nil
    },
}

func init() {
    secretCmd.AddCommand(secretRotateCmd)
    secretCmd.AddCommand(secretShowCmd)
}
