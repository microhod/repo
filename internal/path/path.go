package path

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

func Clean(path string) string {
	// expand home dir
	if strings.HasPrefix(path, "~") {
		homedir, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", homedir, 1)
	}
	return stdpath.Clean(path)
}

func CollapseHomeDir(path string) string {
	homedir, _ := os.UserHomeDir()
	if homedir == "" {
		return path
	}
	if !strings.HasPrefix(path, homedir) {
		return path
	}
	return strings.Replace(path, homedir, "~", 1)
}

// FindDir finds directories with the name 'dir' in 'base'.
// Nested matches are ignored e.g. if dir=b nd /a/b/ matches, then /a/b/c/b will be ignored
func FindDir(base, dir string) ([]string, error) {
	// get top level directories
	dirs, err := getDirs(base)
	if err != nil {
		return nil, fmt.Errorf("getting dirs in [%s]: %w", base, err)
	}

	results := make(chan []string, len(dirs))
	var paths []string

	read := new(sync.WaitGroup)
	read.Go(func() {
		for r := range results {
			paths = append(paths, r...)
		}
	})
	// search concurrently through each top level directory
	write := new(errgroup.Group)
	for i := range dirs {
		d := dirs[i]
		write.Go(func() error {
			p, err := findDir(filepath.Join(base, d), dir)
			results <- p
			return err
		})
	}
	if err := write.Wait(); err != nil {
		return nil, fmt.Errorf("searching for [%s] in [%s]: %w", dir, base, err)
	}
	close(results)

	read.Wait()
	return paths, nil
}

func findDir(base, dir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, er error) error {
		if d == nil || !d.IsDir() {
			return nil
		}

		dirs, err := getDirs(path)
		if errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrNotExist) {
			// skip folder if there are permission errors / folder does not exist
			return filepath.SkipDir
		}
		if err != nil {
			return fmt.Errorf("searching [%s]: %w", path, err)
		}

		for _, d := range dirs {
			if d == dir {
				paths = append(paths, stdpath.Join(path, d))
				// ignore nested matches
				return filepath.SkipDir
			}
		}
		return nil
	})
	return paths, err
}

func getDirs(path string) ([]string, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		paths = append(paths, d.Name())

	}
	return paths, nil
}
