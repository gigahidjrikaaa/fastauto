package deploy

import (
    "os"
    "path/filepath"

    "fastauto/internal/gitutil"
)

// Run pulls latest and runs deploy.sh in repo
func Run(repo string, env map[string]string) error {
    // autopull
    if err := gitutil.GitPull(repo); err != nil { return err }
    // run deploy.sh if present
    p := filepath.Join(repo, "deploy.sh")
    if _, err := os.Stat(p); err == nil {
        return gitutil.RunScript(p, env, repo)
    }
    return nil
}
