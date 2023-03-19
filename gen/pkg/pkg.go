package pkg

import (
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"
	"kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/reporter"
)

func (dep *Dependency) Download(conf *conf.Config) (string, error) {
	if d, ok := dep.GetDependency().(*Dependency_Git); ok {
		return d.Git.Download(conf)
	}
	return "", nil
}

func (dep *GitDependency) Download(conf *conf.Config) (string, error) {
	repoURL := dep.GetGit()

	localPath := git.ParseLocalPathFromGitUrl(filepath.Dir(conf.KclModPath), dep.GetGit())
	repo, err := git.Clone(repoURL, localPath)

	if err != nil {
		reporter.Report("kpm: git clone error:", err)
		return localPath, err
	}
	// checkout branch

	err = git.CheckoutBranch(repo, dep.GetBranch())

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	// checkout commit

	err = git.CheckoutCommit(repo, dep.GetCommit())

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	// checkout tag

	err = git.CheckoutTag(repo, dep.GetTag())

	if err != nil {
		return localPath, report_checkout_err(repo, err)
	}

	return localPath, err
}

func report_checkout_err(repo *gogit.Repository, origin_err error) error {
	ref, err := repo.Head()

	if err != nil {
		reporter.Fatal("kpm: internal bug, please contact us to fix it.")
	}

	reporter.Report("kpm: checkout error:", origin_err, ".")
	reporter.Report("kpm:", ref.Hash().String(), " is used.")

	return origin_err
}
