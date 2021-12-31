// Package path deals with local file paths
package path

import (
	"errors"
	"io/fs"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Clean(path string) string {
	path = stdpath.Clean(path)

	if strings.HasPrefix(path, "~") {
		homedir, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", homedir, 1)
	}

	return path
}

func CollapseUserDir(path string) string {
	homedir, _ := os.UserHomeDir()
	if homedir == "" {
		return path
	}
	if !strings.HasPrefix(path, homedir) {
		return path
	}

	path = strings.Replace(path, homedir, "~", 1)

	return path
}

func FindDir(base, dir string) ([]string, error) {
	dirs, err := GetDirs(base)
	if err != nil {
		return nil, err
	}

	results := make(chan []string, len(dirs))
	paths := []string{}

	read := &errgroup.Group{}
	read.Go(func() error {
		for r := range results {
			paths = append(paths, r...)
		}
		return nil
	})

	// search concurrently through each top level folder
	write := &errgroup.Group{}
	for i := range dirs {
		d := dirs[i]
		write.Go(func() error {
			p, err := findDir(filepath.Join(base, d), dir)
			results <- p
			return err
		})
	}
	if err := write.Wait(); err != nil {
		return nil, err
	}
	close(results)

	read.Wait()
	return paths, nil
}

func findDir(base, dir string) ([]string, error) {
	paths := []string{}
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, er error) error {
		if d == nil || !d.IsDir() {
			return nil
		}

		dirs, err := GetDirs(path)
		if errors.Is(err, fs.ErrPermission) {
			// skip folder if there are permission errors
			return filepath.SkipDir
		}
		if errors.Is(err, fs.ErrNotExist) {
			// skip folder if it does not exist
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}

		for _, d := range dirs {
			if d == dir {
				paths = append(paths, stdpath.Join(path, d))
				return filepath.SkipDir
			}
		}
		return nil
	})
	return paths, err
}

func GetDirs(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dirs, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return dirs, nil
}
