package webhook

import (
    "crypto/hmac"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/hex"
)

func signSHA256(body []byte, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    return hex.EncodeToString(mac.Sum(nil))
}

func signSHA1(body []byte, secret string) string {
    mac := hmac.New(sha1.New, []byte(secret))
    mac.Write(body)
    return hex.EncodeToString(mac.Sum(nil))
}

func subtleCompare(a, b string) bool {
    // constant time compare
    if len(a) != len(b) { return false }
    var v byte
    for i := 0; i < len(a); i++ {
        v |= a[i] ^ b[i]
    }
    return v == 0
}

