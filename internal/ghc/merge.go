package ghc

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var ErrDestinationMustBeDirectory = errors.New("destination path must be a directory")

const contributionRepositoryFile = ".contribution_repository"
const FirstCommitMsg = "init repository"

func Merge(dir string, repoToMerge Repo, gitHubEmail string) error {
	if len(repoToMerge.Commits) == 0 || isContributionRepository(repoToMerge.Path) {
		return nil
	}

	r, err := openCreateRepository(dir, gitHubEmail)
	if err != nil {
		return err
	}

	commitTime, err := getCommitDates(r)
	if err != nil {
		return err
	}

	commitsImported := 0
	for _, c := range repoToMerge.Commits {
		if _, found := commitTime[c.When.Unix()]; found {
			continue // Ignoring. This commit has already been imported
		}
		if err := commit(r, c, gitHubEmail); err != nil {
			return err
		}
		commitTime[c.When.Unix()] = true

		commitsImported++
	}

	return nil
}

func isContributionRepository(dir string) bool {
	_, err := os.Stat(path.Join(dir, contributionRepositoryFile))
	return !os.IsNotExist(err)
}

func initRepo(dir string, gitHubEmail string) (*git.Repository, error) {
	r, err := git.PlainInit(dir, false)
	if err != nil {
		return nil, err
	}

	filename := path.Join(dir, contributionRepositoryFile)
	if err := ioutil.WriteFile(filename, []byte(""), 0664); err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	if _, err := w.Add(contributionRepositoryFile); err != nil {
		return nil, err
	}

	return r, commitRaw(r, "init repository", time.Now(), gitHubEmail)
}

func openCreateRepository(dir string, gitHubEmail string) (*git.Repository, error) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return initRepo(dir, gitHubEmail)
	}

	if !info.IsDir() {
		return nil, ErrDestinationMustBeDirectory
	}

	isEmpty, err := isDirEmpty(dir)
	if err != nil {
		return nil, err
	}
	if isEmpty {
		return initRepo(dir, gitHubEmail)
	}

	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func isDirEmpty(dir string) (bool, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return false, err
	}

	filenames, err := fd.Readdirnames(1)
	if err != nil && err != io.EOF {
		return false, err
	}
	return len(filenames) == 0, nil
}

func getCommitDates(r *git.Repository) (map[int64]bool, error) {
	commitTime := make(map[int64]bool, 0)

	ref, err := r.Head()
	if err == plumbing.ErrReferenceNotFound {
		return commitTime, nil
	}

	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	readCommit := func(c *object.Commit) error {
		commitTime[c.Author.When.Unix()] = true
		return nil
	}
	if err := cIter.ForEach(readCommit); err != nil {
		return nil, err
	}
	return commitTime, nil
}

func commit(r *git.Repository, c Commit, gitHubEmail string) error {
	return commitRaw(r, "", c.When, gitHubEmail)
}

func commitRaw(r *git.Repository, msg string, when time.Time, gitHubEmail string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitHubEmail,
			Email: gitHubEmail,
			When:  when,
		},
	})
	return err
}
