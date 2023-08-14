//go:build !windows
// +build !windows

package path

import (
	"path/filepath"
)

// BinPath Creates a binary file
func BinPath(dir string, name string) string {
	return filepath.Join(dir, name)
}
