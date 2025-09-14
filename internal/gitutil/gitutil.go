package gitutil

import (
    "bytes"
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// DiscoverRepo returns the repo root directory
func DiscoverRepo(dir string) (string, error) {
    if dir == "" { wd, _ := os.Getwd(); dir = wd }
    // try git toplevel
    c := exec.Command("git", "rev-parse", "--show-toplevel")
    c.Dir = dir
    out, err := c.Output()
    if err != nil {
        return "", errors.New("not a git repository")
    }
    return strings.TrimSpace(string(out)), nil
}

func DetectRemote(dir string) (string, error) {
    c := exec.Command("git", "remote", "get-url", "origin")
    c.Dir = dir
    out, err := c.Output()
    if err != nil { return "", err }
    return strings.TrimSpace(string(out)), nil
}

func CurrentBranch(dir string) (string, error) {
    c := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
    c.Dir = dir
    out, err := c.Output()
    if err != nil { return "", err }
    return strings.TrimSpace(string(out)), nil
}

func GitPull(dir string) error {
    c := exec.Command("git", "pull", "--ff-only")
    c.Dir = dir
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    return c.Run()
}

func WriteFileIfMissing(path string, data []byte, mode os.FileMode) error {
    if _, err := os.Stat(path); err == nil { return nil }
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, mode); err != nil { return err }
    if err := os.Rename(tmp, path); err != nil { return err }
    return nil
}

func SanitizeRefToBranch(ref string) string {
    // refs/heads/main -> main
    if strings.HasPrefix(ref, "refs/heads/") {
        return strings.TrimPrefix(ref, "refs/heads/")
    }
    return ref
}

// RunScript runs a shell script at path with env vars
func RunScript(path string, env map[string]string, dir string) error {
    var envs []string
    envs = append(envs, os.Environ()...)
    for k, v := range env { envs = append(envs, k+"="+v) }
    var shell string
    if p, err := exec.LookPath("bash"); err == nil { shell = p } else { shell = "/bin/sh" }
    c := exec.Command(shell, path)
    c.Dir = dir
    c.Env = envs
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    return c.Run()
}

func CommandOutput(dir string, name string, args ...string) (string, error) {
    c := exec.Command(name, args...)
    c.Dir = dir
    var out bytes.Buffer
    c.Stdout = &out
    c.Stderr = &out
    err := c.Run()
    return out.String(), err
}

func RepoNameFromURL(remote string) string {
    // supports https://github.com/org/repo.git or git@github.com:org/repo.git
    s := strings.TrimSuffix(remote, ".git")
    s = strings.TrimPrefix(s, "git@")
    s = strings.TrimPrefix(s, "https://")
    if i := strings.LastIndex(s, "/"); i != -1 { return s[i+1:] }
    if i := strings.LastIndex(s, ":"); i != -1 { return s[i+1:] }
    return filepath.Base(s)
}

