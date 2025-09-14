package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/gigahidjrikaaa/fastauto/internal/systemd"
)

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show service status",
    RunE: func(cmd *cobra.Command, args []string) error {
        mode := viper.GetString("mode")
        units := []string{"fastauto-webhook.service"}
        if mode == "runner" { units = []string{"fastauto-runner.service"} }
        for _, u := range units {
            out, err := systemd.Status(u)
            if err != nil { fmt.Printf("%s: %v\n", u, err) }
            fmt.Printf("%s:\n%s\n", u, out)
        }
        return nil
    },
}
