package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "fastauto/internal/webhook"
)

var internalServeWebhookCmd = &cobra.Command{
    Use:    "internal-serve-webhook",
    Hidden: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        repo := viper.GetString("repo_path")
        if repo == "" { repo = "." }
        // Ensure repo config is loaded
        if repo != "" {
            v := viper.GetViper()
            v.SetConfigFile(repo + "/.fastauto.yml")
            _ = v.ReadInConfig()
        }
        return webhook.Serve(repo)
    },
}
