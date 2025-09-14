package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "fastauto/internal/systemd"
)

var uninstallCmd = &cobra.Command{
    Use:   "uninstall",
    Short: "Disable and remove services",
    RunE: func(cmd *cobra.Command, args []string) error {
        mode := viper.GetString("mode")
        unit := "fastauto-webhook.service"
        if mode == "runner" { unit = "fastauto-runner.service" }
        if err := systemd.StopAndDisable(unit); err != nil { return err }
        path, err := systemd.RemoveUnit(unit)
        if err != nil { return err }
        fmt.Println("Removed:", path)
        return nil
    },
}

