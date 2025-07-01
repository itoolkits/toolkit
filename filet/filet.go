package filet

import (
	"bufio"
	"os"
)

// ReadAll read all content
func ReadAll(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, 0)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

// MkDir make dir
func MkDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// RemoveFile remove file
func RemoveFile(path string) error {
	return os.RemoveAll(path)
}
