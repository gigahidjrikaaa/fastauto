package assets

import (
    "embed"
    "fmt"
)

//go:embed systemd/* scripts/* deploy.sh workflows/*
var embedded embed.FS

// Asset returns embedded asset by path (relative to assets root)
func Asset(path string) ([]byte, error) {
    b, err := embedded.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("asset %s: %w", path, err)
    }
    return b, nil
}

