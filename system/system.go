package system

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileList will recover all the file names in the given path
func ReadFileList(path string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Adds the file to the list if it is a mp3 file
		if strings.HasSuffix(path, ".mp3") {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	if len(fileList) == 0 {
		return []string{}, fmt.Errorf("no files found in path: %s", path)
	}
	return fileList, nil
}

// ReadFile reads the setup file and returns the file pointer
func ReadFile(filePath string) (*os.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// ReadFileToData reads the setup file and unmarshal the json to the data interface
func ReadFileToData(filePath string, data interface{}) error {
	f, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(data); err != nil {
		return err
	}
	return nil
}

// LinkFile creates a hard link to the file at src named dst
func LinkFile(src, dst string) (err error) {

	// Check if the file exists
	_, err = os.Stat(dst)
	if err == nil {
		return nil
	}

	// Open the source file
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.Link(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// Copy copies the contents of the file at src to a regular file at dst.
// If the file named by dst already exists, it ignores.
// The function does not copy the file mode, file permission bits, or file attributes.
func CopyFile(src, dst string) (err error) {

	// Check if the file exists
	_, err = os.Stat(dst)
	if err == nil {
		return nil
	}

	// Open the source file
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the destination file
	w, err := os.Create(dst)
	if err != nil {
		return err
	}

	// Close the destination file
	defer func() {
		// Report the error, if any, from Close, but do so
		// only if there isn't already an outgoing error.
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	// Copy the file
	_, err = io.Copy(w, r)
	return err
}
