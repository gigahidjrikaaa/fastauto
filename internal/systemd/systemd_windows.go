//go:build windows

package systemd

import "fmt"

func InstallUnit(name string, content []byte) (string, error) { return "", fmt.Errorf("systemd not supported on Windows") }
func EnableAndStart(name string) error { return fmt.Errorf("systemd not supported on Windows") }
func StopAndDisable(name string) error { return fmt.Errorf("systemd not supported on Windows") }
func RemoveUnit(name string) (string, error) { return "", fmt.Errorf("systemd not supported on Windows") }
func Status(name string) (string, error) { return "", fmt.Errorf("systemd not supported on Windows") }
func RenderUnit(tmpl string, vars map[string]string) string { return "" }

