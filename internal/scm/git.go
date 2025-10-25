package scm

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chainguard-dev/git-urls"

	"github.com/microhod/repo/internal/repo"
	"github.com/microhod/repo/internal/path"
)

var cmd func(command string, args ...string) *exec.Cmd = exec.Command

type Git struct {
	DefaultRemotePrefix string
}

// ParseRepoFromRemote parses a repo object from the remote URL (as a raw string)
func (git *Git) ParseRepoFromRemote(rawURL string) (*repo.Repo, error) {
	remote, err := giturls.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	// use default prefix if giturls could not parse the remote url
	if remote.Scheme == "file" {
		rawURL = joinURL(git.DefaultRemotePrefix, rawURL)
		remote, err = giturls.Parse(rawURL)
		if err != nil {
			return nil, err
		}
	}

	return git.parseRepoFromRemote(remote)
}

func (git *Git) parseRepoFromRemote(remote *url.URL) (*repo.Repo, error) {
	repo := &repo.Repo{
		Remote: remote,
		Server: remote.Host,
	}

	// trim .git suffix
	remote.Path = strings.TrimSuffix(remote.Path, ".git")

	// assume url path ends with <owner>/<name>
	// e.g. ssh://git@github.com/microhod/repo.git
	//                           ^^^^^^^^^^^^^^
	parts := strings.Split(remote.Path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("url path has an invalid format, expected '.../<owner>/<name>' but got '%s'", remote.Path)
	}
	repo.Owner = parts[len(parts)-2]
	repo.Name = parts[len(parts)-1]

	return repo, nil
}

func (git *Git) getRemoteURL(path string) (string, error) {
	remotes, err := git.exec(nil, "-C", path, "remote")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(remotes) == "" {
		return "", nil
	}

	names := strings.Split(remotes, "\n")
	remote := names[0]
	// default to 'origin' if it exists
	for _, name := range names {
		if name == "origin" {
			remote = "origin"
		}
	}

	url, err := git.exec(nil, "-C", path, "remote", "get-url", remote)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(url), nil
}

func (git *Git) Clone(repo *repo.Repo, path string, options *CloneOptions) error {
	if options == nil {
		options = &CloneOptions{}
	}

	_, err := git.exec(options.Progress, "clone", repo.Remote.String(), path)
	if err != nil {
		return err
	}

	repo.Local = path
	return nil
}

func (git *Git) FindRepos(base string) ([]*repo.Repo, error) {
	gitdirs, err := path.FindDir(base, ".git")
	if err != nil {
		return nil, fmt.Errorf("finding paths to git repos: %w", err)
	}

	// remove '.git' from the end of the paths
	paths := []string{}
	for _, dir := range gitdirs {
		paths = append(paths, filepath.Dir(dir))
	}

	repos := []*repo.Repo{}
	for _, path := range paths {
		remote, err := git.getRemoteURL(path)
		if err != nil {
			return nil, fmt.Errorf("getting remote URL for repo '%s': %w", path, err)
		}
		if remote == "" {
			continue
		}
		repo, err := git.ParseRepoFromRemote(remote)
		repo.Local = path
		if err != nil {
			return nil, fmt.Errorf("parsing repo for path '%s': %w", path, err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func joinURL(baseURL string, paths ...string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	paths = append([]string{baseURL}, paths...)
	for i := range paths {
		paths[i] = strings.TrimPrefix(paths[i], "/")
		paths[i] = strings.TrimSuffix(paths[i], "/")
	}
	return strings.Join(paths, "/")
}

func (git *Git) exec(progress io.Writer, args ...string) (string, error) {
	cmd := cmd("git", args...)

	cmd.Stdin = os.Stdin
	stdout := new(strings.Builder)
	cmd.Stdout = stdout
	stderr := new(strings.Builder)
	cmd.Stderr = stderr
	if progress != nil {
		cmd.Stderr = io.MultiWriter(stderr, progress)
		cmd.Stdout = io.MultiWriter(stdout, progress)
	}

	if err := cmd.Run(); err != nil {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}
