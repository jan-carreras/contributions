package ghc

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

type Extract struct {
	emails map[string]bool
}

type Repo struct {
	Path    string
	Commits []Commit
}

type Commit struct {
	When time.Time
}

func NewExtract(emails []string) *Extract {
	s := &Extract{emails: make(map[string]bool)}
	for _, e := range emails {
		s.emails[e] = true
	}
	return s
}

func (s *Extract) Repository(repoPath string) (Repo, error) {
	repo := Repo{Path: repoPath}

	r, err := git.PlainOpen(repoPath)
	if err == git.ErrRepositoryNotExists {
		return repo, nil
	} else if err != nil {
		return repo, err
	}

	ref, err := r.Head()
	if err == plumbing.ErrReferenceNotFound {
		return repo, nil
	} else if err != nil {
		return repo, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return repo, err
	}

	commits := make([]Commit, 0)
	readCommit := func(c *object.Commit) error {
		if _, found := s.emails[c.Author.Email]; found {
			c := Commit{When: c.Author.When}
			commits = append(commits, c)
		}
		return nil
	}
	if err := cIter.ForEach(readCommit); err != nil {
		return repo, err
	}

	repo.Commits = commits

	return repo, nil
}
