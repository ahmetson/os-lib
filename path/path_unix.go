//go:build !windows
// +build !windows

package path

import (
	"path/filepath"
)

// BinPath Creates a binary file
func BinPath(url string, name string) string {
	return filepath.Join(c.Bin, url+"/bin")
}
