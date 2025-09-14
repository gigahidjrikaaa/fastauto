package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/gigahidjrikaaa/fastauto/internal/webhook"
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

func init() {
    // Optional TLS flags override config
    internalServeWebhookCmd.Flags().String("tls-cert", "", "path to TLS certificate (PEM)")
    internalServeWebhookCmd.Flags().String("tls-key", "", "path to TLS key (PEM)")
    _ = viper.BindPFlag("webhook.tls_cert_file", internalServeWebhookCmd.Flags().Lookup("tls-cert"))
    _ = viper.BindPFlag("webhook.tls_key_file", internalServeWebhookCmd.Flags().Lookup("tls-key"))
}
