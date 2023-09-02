package path

import (
	"fmt"
	"os"
	"path/filepath"
)

// CurrentDir returns the directory of the binary
func CurrentDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Dir(exe)
	return dirPath, nil
}

// AbsDir returns the absolute path of the dir.
// If the directory is not absolute, then it will make an absolute path from the current directory.
func AbsDir(currentDir string, dirPath string) string {
	if filepath.IsAbs(dirPath) {
		return dirPath
	}

	return filepath.Join(currentDir, dirPath)
}

// FileName returns the file name by removing the directory path
func FileName(filePath string) string {
	return filepath.Base(filePath)
}

// NoExtension returns the file name without an extension
func NoExtension(filename string) string {
	var extension = filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

// DirAndFileName returns the directory, file name without extension part.
//
// The function doesn't validate the path.
// Therefore, call this function after validateServicePath()
func DirAndFileName(fileDir string) (string, string) {
	dir, fileName := filepath.Split(fileDir)

	if len(dir) == 0 {
		dir = "."
	} else {
		dir = dir[0 : len(dir)-1]
	}

	ext := filepath.Ext(fileName)
	if len(ext) > 0 {
		fileName = fileName[:len(fileName)-len(ext)]
	}

	return dir, fileName
}

// FileExist returns true if the file exists. if the path is a directory, it will return an error.
func FileExist(fileDir string) (bool, error) {
	info, err := os.Stat(fileDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("os.Stat('%s'): %w", fileDir, err)
		}
	}

	if info.IsDir() {
		return false, fmt.Errorf("fileDir('%s') is directory", fileDir)
	}

	return true, nil
}

// DirExist returns true if the directory exists. if the path is a file, it will return an error
func DirExist(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("os.Stat('%s'): %w", dir, err)
		}
	}

	if !info.IsDir() {
		return false, fmt.Errorf("dir('%s') is not directory", dir)
	}

	return true, nil
}

// MakeDir creates all the directories, including the nested ones.
// If the directories exist, it will skip it.
// If the directory exists, and it includes the file name, it will throw error.
func MakeDir(path string) error {
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
