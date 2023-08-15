package path

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetExecPath returns the current path of the executable
func GetExecPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// GetPath returns the path related to the execPath.
// If the path itself is absolute, then it's returned directly
func GetPath(execPath string, mainPath string) string {
	if filepath.IsAbs(mainPath) {
		return mainPath
	}

	return filepath.Join(execPath, mainPath)
}

// FileName returns the file name by removing the directory path
func FileName(path string) string {
	return filepath.Base(path)
}

// SplitServicePath returns the directory, file name without extension part.
//
// The function doesn't validate the path.
// Therefore, call this function after validateServicePath()
func SplitServicePath(servicePath string) (string, string) {
	dir, fileName := filepath.Split(servicePath)

	if len(dir) == 0 {
		dir = "."
	}

	fileName = fileName[0 : len(fileName)-4]

	return dir, fileName
}

// FileExists returns true if the file exists. if the path is a directory, it will return an error.
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("os.Stat('%s'): %w", path, err)
		}
	}

	if info.IsDir() {
		return false, fmt.Errorf("path('%s') is directory", path)
	}

	return true, nil
}

// DirExists returns true if the directory exists. if the path is a file, it will return an error
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("os.Stat('%s'): %w", path, err)
		}
	}

	if !info.IsDir() {
		return false, fmt.Errorf("path('%s') is not directory", path)
	}

	return true, nil
}

// MakePath creates all the directories, including the nested ones.
func MakePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0777); err != nil {
				return fmt.Errorf("failed to create a directory at '%s' path: %w", path, err)
			}
			return nil
		} else {
			return fmt.Errorf("failed to read '%s': %w", path, err)
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("the path '%s' is not a directory", path)
	}

	return nil
}
