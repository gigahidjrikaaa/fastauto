package config

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"

    "fastauto/internal/safeio"
)

type RepoConfig struct {
    Mode     string          `yaml:"mode"`
    RepoPath string          `yaml:"repo_path"`
    Branches []string        `yaml:"branches"`
    Webhook  *WebhookConfig  `yaml:"webhook,omitempty"`
    Runner   *RunnerConfig   `yaml:"runner,omitempty"`
}

type WebhookConfig struct {
    Address string `yaml:"address"`
    TLSCertFile string `yaml:"tls_cert_file,omitempty"`
    TLSKeyFile  string `yaml:"tls_key_file,omitempty"`
    // future: path, tls client auth
}

type RunnerConfig struct {
    Labels []string `yaml:"labels,omitempty"`
}

type GlobalConfig struct {
    WebhookSecret string `yaml:"webhook_secret"`
}

func (g *GlobalConfig) GenerateWebhookSecret() error {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil { return err }
    g.WebhookSecret = hex.EncodeToString(b)
    return nil
}

func EnsureGlobalConfig() error {
    p := globalPath()
    dir := filepath.Dir(p)
    if err := os.MkdirAll(dir, 0o700); err != nil { return err }
    if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
        return SaveGlobalConfig(&GlobalConfig{})
    }
    return nil
}

func LoadGlobalConfig() (*GlobalConfig, error) {
    _ = EnsureGlobalConfig()
    p := globalPath()
    b, err := os.ReadFile(p)
    if err != nil { return &GlobalConfig{}, nil }
    var g GlobalConfig
    if err := yaml.Unmarshal(b, &g); err != nil { return &GlobalConfig{}, nil }
    return &g, nil
}

func SaveGlobalConfig(g *GlobalConfig) error {
    p := globalPath()
    b, err := yaml.Marshal(g)
    if err != nil { return err }
    return safeio.WriteFileAtomicWithBackup(p, b, 0o600)
}

func globalPath() string {
    xdg := os.Getenv("XDG_CONFIG_HOME")
    if xdg == "" {
        home, _ := os.UserHomeDir()
        xdg = filepath.Join(home, ".config")
    }
    return filepath.Join(xdg, "fastauto", "config.yml")
}

func SaveRepoConfig(path string, c *RepoConfig) error {
    b, err := yaml.Marshal(c)
    if err != nil { return err }
    return safeio.WriteFileAtomicWithBackup(path, b, 0o644)
}
