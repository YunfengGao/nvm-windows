package file

import (
	"archive/zip"
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Function courtesy http://stackoverflow.com/users/1129149/swtdrgn
func Unzip(src, dest string) (string, error) {
	var folderName string
	var firstMatch = true
	r, err := zip.OpenReader(src)
	if err != nil {
		return folderName, err
	}
	defer r.Close()

	for _, f := range r.File {
		if !strings.Contains(f.Name, "..") {
			rc, err := f.Open()
			if err != nil {
				return folderName, err
			}
			defer rc.Close()

			fpath := filepath.Join(dest, f.Name)
			if f.FileInfo().IsDir() {
				if firstMatch {
					var fdir string
					if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
						fdir = fpath[lastIndex:]
						if strings.HasPrefix(fdir, "\\npm-cli-") {
							folderName = fdir
							firstMatch = false
						}
					}
				}
				os.MkdirAll(fpath, f.Mode())
			} else {
				var fdir string
				if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
					fdir = fpath[:lastIndex]
				}

				err = os.MkdirAll(fdir, f.Mode())
				if err != nil {
					log.Fatal(err)
					return folderName, err
				}
				f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					return folderName, err
				}
				defer f.Close()

				_, err = io.Copy(f, rc)
				if err != nil {
					return folderName, err
				}
			}
		} else {
			log.Printf("failed to extract file: %s (cannot validate)\n", f.Name)
		}
	}

	return folderName, nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Rename(src, dest string) error {
	if src == dest {
		return nil
	}
	if Exists(dest) {
		log.Printf("\nfile %s already exists, please delete\n", dest)
		return nil
	}

	err := os.Rename(src, dest)
	if err != nil {
		log.Printf("\nfailed to rename file: %s to %s, err=%s\n", src, dest, err)
		return err
	}
	log.Printf("\nsuccess to rename file: %s to %s\n", src, dest)
	return nil
}
