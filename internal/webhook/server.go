package webhook

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/spf13/viper"

    "fastauto/internal/config"
    "fastauto/internal/deploy"
    "fastauto/internal/gitutil"
    "fastauto/internal/match"
)

type pushEvent struct {
    Ref string `json:"ref"`
}

// Serve starts the webhook server using configuration
func Serve(repo string) error {
    g, _ := config.LoadGlobalConfig()
    addr := viper.GetString("webhook.address")
    if addr == "" { addr = ":8080" }
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request){ w.WriteHeader(200); _, _ = w.Write([]byte("ok")) })
    mux.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil { http.Error(w, "read", 400); return }
        sig := r.Header.Get("X-Hub-Signature-256")
        if sig == "" { sig = r.Header.Get("X-Hub-Signature") }
        if !VerifyHMAC(body, g.WebhookSecret, sig) {
            http.Error(w, "invalid signature", http.StatusUnauthorized)
            return
        }
        // Only handle push
        if e := r.Header.Get("X-GitHub-Event"); e != "push" {
            w.WriteHeader(204); return
        }
        var ev pushEvent
        if err := json.Unmarshal(body, &ev); err != nil { http.Error(w, "bad payload", 400); return }
        branch := gitutil.SanitizeRefToBranch(ev.Ref)
        // branch filter from viper
        branches := viper.GetStringSlice("branches")
        allowed := false
        if len(branches) == 0 {
            allowed = true
        } else {
            for _, p := range branches {
                if match.Glob(p, branch) { allowed = true; break }
            }
        }
        if !allowed {
            w.WriteHeader(204); return
        }
        // set env and run deployment
        env := map[string]string{
            "FA_REPO": repo,
            "FA_BRANCH": branch,
            "FA_EVENT": "push",
        }
        go func(){
            if err := deploy.Run(repo, env); err != nil {
                log.Printf("deploy error: %v", err)
            }
        }()
        w.WriteHeader(202)
        _, _ = w.Write([]byte("queued"))
    })
    log.Printf("fastauto webhook listening on %s, repo=%s", addr, repo)
    return http.ListenAndServe(addr, mux)
}

// VerifyHMAC checks GitHub signature headers using sha256
func VerifyHMAC(body []byte, secret, header string) bool {
    if secret == "" || header == "" { return false }
    // header can be "sha256=..." or "sha1=..."
    parts := strings.SplitN(header, "=", 2)
    if len(parts) != 2 { return false }
    algo := parts[0]
    sig := parts[1]
    var calc string
    switch algo {
    case "sha256":
        calc = signSHA256(body, secret)
    case "sha1":
        calc = signSHA1(body, secret)
    default:
        return false
    }
    return subtleCompare(calc, sig)
}

// for tests we expose via separate file

