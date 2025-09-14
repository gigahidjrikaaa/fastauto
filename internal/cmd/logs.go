package cmd

import (
    "fmt"
    "os"
    "os/exec"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var logsFollow bool

var logsCmd = &cobra.Command{
    Use:   "logs",
    Short: "Show journald logs for fastauto services",
    RunE: func(cmd *cobra.Command, args []string) error {
        mode := viper.GetString("mode")
        unit := "fastauto-webhook.service"
        if mode == "runner" { unit = "fastauto-runner.service" }
        args := []string{"-u", unit}
        if logsFollow { args = append(args, "-f") }
        c := exec.Command("journalctl", args...)
        c.Stdout = os.Stdout
        c.Stderr = os.Stderr
        fmt.Printf("journalctl %v\n", args)
        return c.Run()
    },
}

func init() {
    logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "follow logs")
}

