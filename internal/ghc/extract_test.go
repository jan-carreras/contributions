package ghc_test

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jan-carreras/contributions/internal/ghc"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testingFile = "testing-git-file"

func TestSearch_Repository(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitSearch")
	checkIfError(t, err)

	defer os.RemoveAll(dir)

	prepareTestRepository(t, dir)

	s := ghc.NewExtract([]string{"a@example.com", "b@example.com"})
	r, err := s.Repository(dir)
	checkIfError(t, err)

	if r.Path != dir {
		t.Errorf("Repository path missmatch. Expect %v, have %v", dir, r.Path)
	}

	if len(r.Commits) != 4 {
		t.Errorf("expecting %v commits but only %v found", 4, len(r.Commits))
	}
}
func prepareTestRepository(t *testing.T, dir string) {
	r, err := git.PlainInit(dir, false)
	checkIfError(t, err)

	w, err := r.Worktree()
	checkIfError(t, err)

	createFile(t, dir, "test 1")
	_, err = w.Add(testingFile)
	checkIfError(t, err)
	commit(t, w, "commit #1", "a@example.com")

	createFile(t, dir, "test 2")
	_, err = w.Add(testingFile)
	checkIfError(t, err)
	commit(t, w, "commit #2", "b@example.com")

	createFile(t, dir, "test 3")
	_, err = w.Add(testingFile)
	checkIfError(t, err)
	commit(t, w, "commit #3", "c@example.com")

	createFile(t, dir, "test 4")
	_, err = w.Add(testingFile)
	checkIfError(t, err)
	commit(t, w, "commit #4", "a@example.com")

	createFile(t, dir, "test 5")
	_, err = w.Add(testingFile)
	checkIfError(t, err)
	commit(t, w, "commit #5", "b@example.com")
}

func commit(t *testing.T, w *git.Worktree, msg, email string) {
	_, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  email,
			Email: email,
			When:  time.Now(),
		},
	})
	checkIfError(t, err)
}

func createFile(t *testing.T, dir, content string) {
	filename := filepath.Join(dir, testingFile)
	err := ioutil.WriteFile(filename, []byte(content), 0600)
	checkIfError(t, err)
}
