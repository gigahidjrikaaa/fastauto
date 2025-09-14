//go:build !windows

package systemd

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "text/template"
)

func isRoot() bool { return os.Geteuid() == 0 }

func systemctlArgs() []string {
    if isRoot() { return []string{} }
    return []string{"--user"}
}

func unitDir() string {
    if isRoot() { return "/etc/systemd/system" }
    xdg := os.Getenv("XDG_CONFIG_HOME")
    if xdg == "" { home, _ := os.UserHomeDir(); xdg = filepath.Join(home, ".config") }
    return filepath.Join(xdg, "systemd", "user")
}

func RenderUnit(tmpl string, vars map[string]string) string {
    t := template.Must(template.New("u").Parse(tmpl))
    var b bytes.Buffer
    _ = t.Execute(&b, vars)
    return b.String()
}

func InstallUnit(name string, content []byte) (string, error) {
    dir := unitDir()
    if err := os.MkdirAll(dir, 0o755); err != nil { return "", err }
    path := filepath.Join(dir, name)
    if err := os.WriteFile(path, content, 0o644); err != nil { return "", err }
    // daemon-reload
    args := append(systemctlArgs(), "daemon-reload")
    c := exec.Command("systemctl", args...)
    _ = c.Run()
    return path, nil
}

func EnableAndStart(name string) error {
    for _, sub := range [][]string{{"enable", name}, {"start", name}} {
        args := append(systemctlArgs(), sub...)
        c := exec.Command("systemctl", args...)
        c.Stdout = os.Stdout
        c.Stderr = os.Stderr
        if err := c.Run(); err != nil { return fmt.Errorf("systemctl %s: %w", strings.Join(args, " "), err) }
    }
    return nil
}

func StopAndDisable(name string) error {
    for _, sub := range [][]string{{"stop", name}, {"disable", name}} {
        args := append(systemctlArgs(), sub...)
        c := exec.Command("systemctl", args...)
        if err := c.Run(); err != nil { return err }
    }
    return nil
}

func RemoveUnit(name string) (string, error) {
    path := filepath.Join(unitDir(), name)
    if err := os.Remove(path); err != nil { return "", err }
    _ = exec.Command("systemctl", append(systemctlArgs(), "daemon-reload")...).Run()
    return path, nil
}

func Status(name string) (string, error) {
    args := append(systemctlArgs(), "status", name, "--no-pager")
    out, err := exec.Command("systemctl", args...).CombinedOutput()
    return string(out), err
}
