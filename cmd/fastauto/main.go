package main

import (
    "log"
    "runtime"

    icmd "fastauto/internal/cmd"
)

var (
    version = "dev"
    commit  = ""
    date    = ""
)

func main() {
    // Ensure static-friendly defaults
    log.SetFlags(0)
    icmd.BuildInfo = icmd.Info{Version: version, Commit: commit, Date: date, Go: runtime.Version()}
    if err := icmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

