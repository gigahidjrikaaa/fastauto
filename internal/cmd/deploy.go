package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/gigahidjrikaaa/fastauto/internal/deploy"
)

var deployNow bool

var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Trigger a deployment (runs deploy.sh)",
    RunE: func(cmd *cobra.Command, args []string) error {
        repo := viper.GetString("repo_path")
        if repo == "" { repo = "." }
        if deployNow {
            if err := deploy.Run(repo, nil); err != nil { return err }
            fmt.Println("Deployment completed")
            return nil
        }
        fmt.Println("Use --now to run deploy.sh immediately")
        return nil
    },
}

func init() {
    deployCmd.Flags().BoolVar(&deployNow, "now", false, "run deploy.sh now")
}
