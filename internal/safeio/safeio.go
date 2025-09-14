package safeio

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "time"
)

// WriteFileAtomicWithBackup writes file atomically and creates a timestamped backup if exists
func WriteFileAtomicWithBackup(path string, data []byte, perm fs.FileMode) error {
    // backup existing
    if st, err := os.Stat(path); err == nil && st.Mode().IsRegular() {
        ts := time.Now().Format("20060102T150405")
        bak := fmt.Sprintf("%s.bak.%s", path, ts)
        _ = os.Rename(path, bak)
    } else {
        // ensure dir
        _ = os.MkdirAll(filepath.Dir(path), 0o755)
    }
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, perm); err != nil { return err }
    return os.Rename(tmp, path)
}

