package cmd

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/gigahidjrikaaa/fastauto/internal/assets"
    "github.com/gigahidjrikaaa/fastauto/internal/config"
    "github.com/gigahidjrikaaa/fastauto/internal/gitutil"
    "github.com/gigahidjrikaaa/fastauto/internal/safeio"
)

var (
    initMode   string
    initPort   int
    initBranch string
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize fastauto in this repository",
    RunE: func(cmd *cobra.Command, args []string) error {
        repo, err := gitutil.DiscoverRepo(repoDir)
        if err != nil {
            return err
        }
        remote, _ := gitutil.DetectRemote(repo)
        branch, _ := gitutil.CurrentBranch(repo)
        if initBranch != "" {
            branch = initBranch
        }

        fmt.Printf("Detected repo: %s\n", repo)
        if remote != "" { fmt.Printf("Remote: %s\n", remote) }
        if branch != "" { fmt.Printf("Branch: %s\n", branch) }

        // Prompt if not provided
        if initMode == "" {
            fmt.Print("Select mode [webhook|runner] (default webhook): ")
            r := bufio.NewReader(os.Stdin)
            s, _ := r.ReadString('\n')
            s = strings.TrimSpace(s)
            if s == "runner" { initMode = "runner" } else { initMode = "webhook" }
        }

        // Prepare repo config
        rcfg := config.RepoConfig{
            Mode:     initMode,
            RepoPath: repo,
            Branches: []string{branch},
        }
        if initMode == "webhook" {
            if initPort == 0 { initPort = 8080 }
            rcfg.Webhook = &config.WebhookConfig{Address: fmt.Sprintf(":%d", initPort)}
        } else {
            rcfg.Runner = &config.RunnerConfig{}
        }

        // Ensure deploy.sh exists
        dep := filepath.Join(repo, "deploy.sh")
        if _, err := os.Stat(dep); os.IsNotExist(err) {
            data, _ := assets.Asset("deploy.sh")
            if err := safeio.WriteFileAtomicWithBackup(dep, data, 0o755); err != nil {
                return err
            }
            fmt.Println("Wrote deploy.sh")
        }

        // Write workflow for runner mode
        if initMode == "runner" {
            wfPath := filepath.Join(repo, ".github", "workflows")
            _ = os.MkdirAll(wfPath, 0o755)
            wf := filepath.Join(wfPath, "fastauto.yml")
            data, _ := assets.Asset("workflows/fastauto.yml")
            // Replace branches placeholder
            b := strings.Join(rcfg.Branches, "\n        - ")
            data = []byte(strings.ReplaceAll(string(data), "{{BRANCHES}}", b))
            if err := safeio.WriteFileAtomicWithBackup(wf, data, 0o644); err != nil {
                return err
            }
            fmt.Println("Wrote .github/workflows/fastauto.yml")
        }

        // Persist repo config in .fastauto.yml
        repoCfgPath := filepath.Join(repo, ".fastauto.yml")
        if err := config.SaveRepoConfig(repoCfgPath, &rcfg); err != nil { return err }
        fmt.Println("Wrote .fastauto.yml")

        // Write global config dir exists
        _ = config.EnsureGlobalConfig()

        // Offer to install systemd units now
        fmt.Print("Install systemd service now? [Y/n]: ")
        r := bufio.NewReader(os.Stdin)
        s, _ := r.ReadString('\n')
        s = strings.TrimSpace(strings.ToLower(s))
        if s == "" || s == "y" || s == "yes" {
            // Reuse install command logic
            installMode = initMode
            if err := installCmd.RunE(installCmd, nil); err != nil { return err }
        } else {
            fmt.Println("You can install later via 'fastauto install'.")
        }

        // Update viper runtime
        viper.Set("mode", initMode)
        viper.Set("repo_path", repo)
        return nil
    },
}

func init() {
    initCmd.Flags().StringVar(&initMode, "mode", "", "mode: webhook or runner")
    initCmd.Flags().IntVar(&initPort, "port", 0, "webhook listen port")
    initCmd.Flags().StringVar(&initBranch, "branch", "", "default branch to watch")
}
