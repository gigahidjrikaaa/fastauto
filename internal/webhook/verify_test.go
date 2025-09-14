package webhook

import "testing"

func TestVerifyHMAC_SHA256(t *testing.T) {
    body := []byte(`{"ref":"refs/heads/main"}`)
    secret := "supersecret"
    hdr := "sha256=" + signSHA256(body, secret)
    if !VerifyHMAC(body, secret, hdr) {
        t.Fatal("expected verify true")
    }
}

func TestVerifyHMAC_SHA1(t *testing.T) {
    body := []byte("hi")
    secret := "s"
    hdr := "sha1=" + signSHA1(body, secret)
    if !VerifyHMAC(body, secret, hdr) {
        t.Fatal("expected verify true")
    }
}

func TestVerifyHMAC_Bad(t *testing.T) {
    body := []byte("hi")
    secret := "s"
    hdr := "sha256=deadbeef"
    if VerifyHMAC(body, secret, hdr) {
        t.Fatal("expected verify false")
    }
}

