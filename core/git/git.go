package git

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitOption struct {
	Url    string
	Branch string
	Commit string
	Tag    string
}

func NewGitOption() *GitOption {
	return &GitOption{}
}

func (opt *GitOption) SetUrl(url string) *GitOption {
	opt.Url = url
	return opt
}

func (opt *GitOption) SetCommit(commit string) *GitOption {
	opt.Commit = commit
	return opt
}

func (opt *GitOption) SetBranch(branch string) *GitOption {
	opt.Branch = branch
	return opt
}

func (opt *GitOption) SetTag(tag string) *GitOption {
	opt.Tag = tag
	return opt
}

func Clone(repoURL string, localPath string) (*git.Repository, error) {
	repo, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	return repo, err
}

func CheckoutCommit(repo *git.Repository, commit string) error {
	return checkout(repo, plumbing.NewHash(commit))
}

func CheckoutBranch(repo *git.Repository, branch string) error {
	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := repo.Reference(refName, true)
	if err == nil {
		return checkout(repo, ref.Hash())
	}
	return err
}

func CheckoutTag(repo *git.Repository, tag string) error {
	tagRefName := plumbing.NewTagReferenceName(tag)
	ref, err := repo.Reference(tagRefName, true)

	if err == nil {
		return checkout(repo, ref.Hash())
	}
	return err
}

func checkout(repo *git.Repository, ref plumbing.Hash) error {
	worktree, err := repo.Worktree()

	if err == nil {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: ref,
		})
	}
	return err
}

func ParseLocalPathFromGitUrl(rootPath string, gitUrl string) string {
	parsedUrl, _ := url.Parse(gitUrl)
	pathWithoutScheme := strings.TrimPrefix(gitUrl, parsedUrl.Scheme+"://")

	fileExt := filepath.Ext(pathWithoutScheme)
	return filepath.Join(rootPath, filepath.Base(pathWithoutScheme[:len(pathWithoutScheme)-len(fileExt)]))
}

func ParseRepoNameFromGitUrl(gitUrl string) string {
	name := filepath.Base(gitUrl)
	return name[:len(name)-len(filepath.Ext(name))]
}
