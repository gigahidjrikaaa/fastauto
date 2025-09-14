package cmd

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "fastauto/internal/assets"
    "fastauto/internal/config"
    "fastauto/internal/safeio"
    "fastauto/internal/systemd"
)

var installMode string

var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install and enable services for webhook or runner",
    RunE: func(cmd *cobra.Command, args []string) error {
        mode := installMode
        if mode == "" {
            mode = viper.GetString("mode")
        }
        if mode == "" { return errors.New("mode not set; run fastauto init or pass --mode") }

        repo := viper.GetString("repo_path")
        if repo == "" {
            wd, _ := os.Getwd(); repo = wd
        }

        if mode == "webhook" {
            // ensure secret
            g, _ := config.LoadGlobalConfig()
            if g.WebhookSecret == "" {
                if err := g.GenerateWebhookSecret(); err != nil { return err }
                if err := config.SaveGlobalConfig(g); err != nil { return err }
                fmt.Println("Generated webhook secret in global config")
            }
            // write unit
            content, _ := assets.Asset("systemd/fastauto-webhook.service.tmpl")
            exe, _ := os.Executable()
            unit := systemd.RenderUnit(string(content), map[string]string{
                "ExecStart": fmt.Sprintf("%s internal-serve-webhook --repo %s", exe, repo),
                "Description": "fastauto webhook server",
            })
            path, err := systemd.InstallUnit("fastauto-webhook.service", []byte(unit))
            if err != nil { return err }
            if err := systemd.EnableAndStart("fastauto-webhook.service"); err != nil { return err }
            fmt.Println("Installed:", path)
            fmt.Println("Enabled and started fastauto-webhook.service")
            return nil
        }

        if mode == "runner" {
            // Download/install runner via helper script into repo/.fastauto/runner
            runnerDir := filepath.Join(repo, ".fastauto", "runner")
            _ = os.MkdirAll(runnerDir, 0o755)
            script, _ := assets.Asset("scripts/install_runner.sh")
            scriptPath := filepath.Join(runnerDir, "install_runner.sh")
            if err := safeio.WriteFileAtomicWithBackup(scriptPath, script, 0o755); err != nil { return err }
            // Provide systemd unit that calls runner's run.sh
            content, _ := assets.Asset("systemd/fastauto-runner.service.tmpl")
            unit := systemd.RenderUnit(string(content), map[string]string{
                "ExecStart": fmt.Sprintf("/bin/bash -lc '%s'", filepath.Join(runnerDir, "run.sh")),
                "Description": "fastauto self-hosted GitHub Actions runner",
                "WorkingDirectory": runnerDir,
            })
            path, err := systemd.InstallUnit("fastauto-runner.service", []byte(unit))
            if err != nil { return err }
            if err := systemd.EnableAndStart("fastauto-runner.service"); err != nil { return err }
            fmt.Println("Installed:", path)
            fmt.Println("Enabled and started fastauto-runner.service")
            // Write workflow when runner mode
            wfPath := filepath.Join(repo, ".github", "workflows")
            _ = os.MkdirAll(wfPath, 0o755)
            wf := filepath.Join(wfPath, "fastauto.yml")
            if _, err := os.Stat(wf); err != nil {
                data, _ := assets.Asset("workflows/fastauto.yml")
                if err := safeio.WriteFileAtomicWithBackup(wf, data, 0o644); err != nil { return err }
                fmt.Println("Wrote .github/workflows/fastauto.yml")
            }
            return nil
        }

        return fmt.Errorf("unknown mode: %s", mode)
    },
}

func init() {
    installCmd.Flags().StringVar(&installMode, "mode", "", "mode: webhook or runner")
}
