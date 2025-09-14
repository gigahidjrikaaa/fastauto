package cmd

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/spf13/cobra"

    "github.com/gigahidjrikaaa/fastauto/internal/safeio"
)

var (
    selfBinDir string
    selfForce  bool
)

var selfCmd = &cobra.Command{
    Use:   "self",
    Short: "Manage the fastauto CLI installation",
}

var selfInstallCmd = &cobra.Command{
    Use:   "install",
    Short: "Install this fastauto binary into your PATH",
    RunE: func(cmd *cobra.Command, args []string) error {
        exe, err := os.Executable()
        if err != nil { return err }
        targetDir := selfBinDir
        if targetDir == "" {
            if os.Geteuid() == 0 {
                targetDir = "/usr/local/bin"
            } else {
                home, _ := os.UserHomeDir()
                targetDir = filepath.Join(home, ".local", "bin")
            }
        }
        if err := os.MkdirAll(targetDir, 0o755); err != nil { return err }
        dst := filepath.Join(targetDir, "fastauto")
        if !selfForce {
            if _, err := os.Stat(dst); err == nil {
                bak := fmt.Sprintf("%s.bak.%s", dst, time.Now().Format("20060102T150405"))
                if err := os.Rename(dst, bak); err != nil { return err }
            }
        }
        // copy file contents then make executable with atomic write
        data, err := os.ReadFile(exe)
        if err != nil { return err }
        if err := safeio.WriteFileAtomicWithBackup(dst, data, 0o755); err != nil { return err }

        // PATH guidance
        pathEnv := os.Getenv("PATH")
        inPath := false
        for _, p := range strings.Split(pathEnv, string(os.PathListSeparator)) {
            if p == targetDir { inPath = true; break }
        }
        fmt.Printf("Installed fastauto to %s\n", dst)
        if !inPath {
            fmt.Printf("Note: %s is not in PATH. Add it to your shell profile.\n", targetDir)
        }
        return nil
    },
}

var selfUninstallCmd = &cobra.Command{
    Use:   "uninstall",
    Short: "Remove the installed fastauto from your PATH",
    RunE: func(cmd *cobra.Command, args []string) error {
        targetDir := selfBinDir
        if targetDir == "" {
            if os.Geteuid() == 0 {
                targetDir = "/usr/local/bin"
            } else {
                home, _ := os.UserHomeDir()
                targetDir = filepath.Join(home, ".local", "bin")
            }
        }
        dst := filepath.Join(targetDir, "fastauto")
        if err := os.Remove(dst); err != nil {
            if os.IsNotExist(err) { fmt.Printf("Not found: %s\n", dst); return nil }
            return err
        }
        fmt.Printf("Removed %s\n", dst)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(selfCmd)
    selfCmd.AddCommand(selfInstallCmd)
    selfCmd.AddCommand(selfUninstallCmd)
    selfCmd.PersistentFlags().StringVar(&selfBinDir, "bin-dir", "", "install directory (default: ~/.local/bin or /usr/local/bin as root)")
    selfCmd.PersistentFlags().BoolVar(&selfForce, "force", false, "overwrite existing without backup")
}

