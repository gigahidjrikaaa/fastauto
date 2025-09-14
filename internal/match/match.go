package match

import "path/filepath"

// Glob matches branch name against pattern supporting * and ?
func Glob(pattern, name string) bool {
    ok, err := filepath.Match(pattern, name)
    if err != nil { return false }
    return ok
}

