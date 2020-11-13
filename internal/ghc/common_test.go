package ghc_test

import (
	"fgh/internal/ghc"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"testing"
)

func checkIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func readCommits(t *testing.T, r *git.Repository) []*object.Commit {
	ref, err := r.Head()
	checkIfError(t, err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	checkIfError(t, err)

	commits := make([]*object.Commit, 0)
	readCommit := func(c *object.Commit) error {
		if c.Message == ghc.FirstCommitMsg {
			return nil
		}
		commits = append(commits, c)
		return nil
	}
	err = cIter.ForEach(readCommit)
	checkIfError(t, err)

	return commits
}

func countCommits(t *testing.T, r *git.Repository) int {
	return len(readCommits(t, r))
}
