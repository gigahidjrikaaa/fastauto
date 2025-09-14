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
        jArgs := []string{"-u", unit}
        if logsFollow { jArgs = append(jArgs, "-f") }
        c := exec.Command("journalctl", jArgs...)
        c.Stdout = os.Stdout
        c.Stderr = os.Stderr
        fmt.Printf("journalctl %v\n", jArgs)
        return c.Run()
    },
}

func init() {
    logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "follow logs")
}
