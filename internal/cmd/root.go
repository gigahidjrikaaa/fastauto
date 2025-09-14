package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// Info holds build info injected via ldflags
type Info struct {
    Version string
    Commit  string
    Date    string
    Go      string
}

var BuildInfo Info

var (
    cfgFile string
    repoDir string
)

var rootCmd = &cobra.Command{
    Use:   "fastauto",
    Short: "Autopull + auto-deploy for any git repo",
    Long:  "fastauto automates autopull and auto-deploy via GitHub runner or webhook.",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        if repoDir != "" { viper.Set("repo_path", repoDir) }
    },
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default .fastauto.yml and $XDG_CONFIG_HOME/fastauto/config.yml)")
    rootCmd.PersistentFlags().StringVar(&repoDir, "repo", "", "path to repository (default: auto-detect)*)")

    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(installCmd)
    rootCmd.AddCommand(statusCmd)
    rootCmd.AddCommand(logsCmd)
    rootCmd.AddCommand(deployCmd)
    rootCmd.AddCommand(secretCmd)
    rootCmd.AddCommand(uninstallCmd)
    // hidden internal commands used by systemd
    rootCmd.AddCommand(internalServeWebhookCmd)
}

func initConfig() {
    v := viper.GetViper()
    v.SetConfigType("yaml")
    v.SetEnvPrefix("FASTAUTO")
    v.AutomaticEnv()

    // Load repo-local config if present
    if cfgFile != "" {
        v.SetConfigFile(cfgFile)
        _ = v.ReadInConfig()
    } else {
        // Repo local
        wd, _ := os.Getwd()
        if wd != "" {
            repoCfg := filepath.Join(wd, ".fastauto.yml")
            if _, err := os.Stat(repoCfg); err == nil {
                v.SetConfigFile(repoCfg)
                _ = v.ReadInConfig()
            }
        }
        // Global config
        if v.ConfigFileUsed() == "" {
            xdg := os.Getenv("XDG_CONFIG_HOME")
            if xdg == "" {
                home, _ := os.UserHomeDir()
                xdg = filepath.Join(home, ".config")
            }
            g := filepath.Join(xdg, "fastauto", "config.yml")
            if _, err := os.Stat(g); err == nil {
                v.SetConfigFile(g)
                _ = v.ReadInConfig()
            }
        }
    }
}

// Execute runs the root command
func Execute() error { return rootCmd.Execute() }

// Version string
func versionString() string {
    if BuildInfo.Version == "" {
        return "dev"
    }
    extra := ""
    if BuildInfo.Commit != "" { extra = fmt.Sprintf(" (%s)", BuildInfo.Commit[:min(7, len(BuildInfo.Commit))]) }
    return fmt.Sprintf("%s%s", BuildInfo.Version, extra)
}

func min(a,b int) int { if a<b {return a}; return b }
