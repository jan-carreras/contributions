package ghc_test

import (
	"github.com/go-git/go-git/v5"
	"github.com/jan-carreras/contributions/internal/ghc"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

var tz = time.FixedZone("", 0)
var commit1 = ghc.Commit{When: time.Date(2020, 11, 13, 11, 35, 00, 0, tz)}
var gitHubEmail = "text@example.com"

func TestMerge_NoCommits(t *testing.T) {
	repo := ghc.Repo{}
	if err := ghc.Merge("", repo, gitHubEmail); err != nil {
		t.Fatal(err)
	}
}

func TestMerge_NewRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	repo := ghc.Repo{Commits: []ghc.Commit{commit1}}
	err = ghc.Merge(dir, repo, gitHubEmail)
	checkIfError(t, err)

	r, err := git.PlainOpen(dir)
	checkIfError(t, err)

	commits := countCommits(t, r)
	if commits != 1 {
		t.Fatalf("number of commits expected 1 found %v", commits)
	}
}

func TestMerge_ExistingRepository(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	r, err := git.PlainInit(dir, false)
	checkIfError(t, err)

	repo := ghc.Repo{Commits: []ghc.Commit{commit1}}
	err = ghc.Merge(dir, repo, gitHubEmail)
	checkIfError(t, err)

	commits := countCommits(t, r)
	if commits != len(repo.Commits) {
		t.Fatalf("number of commits expected %v found %v", len(repo.Commits), commits)
	}
}

func TestMerge_ImportMultipleCommits(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	repo := ghc.Repo{
		Commits: []ghc.Commit{
			{When: time.Date(2020, 11, 13, 11, 35, 00, 0, tz)},
			{When: time.Date(2020, 11, 13, 11, 40, 00, 0, tz)},
			{When: time.Date(2020, 11, 13, 11, 45, 00, 0, tz)},
		},
	}
	err = ghc.Merge(dir, repo, gitHubEmail)
	checkIfError(t, err)

	r, err := git.PlainOpen(dir)
	checkIfError(t, err)

	commits := countCommits(t, r)
	if commits != len(repo.Commits) {
		t.Fatalf("number of commits expected %v found %v", len(repo.Commits), commits)
	}
}

func TestMerge_MultipleImports(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	repo1 := ghc.Repo{
		Commits: []ghc.Commit{
			{When: time.Date(2020, 11, 13, 11, 35, 00, 0, tz)},
			{When: time.Date(2020, 11, 13, 11, 40, 00, 0, tz)},
			{When: time.Date(2020, 11, 13, 11, 45, 00, 0, tz)},
		},
	}
	repo2 := ghc.Repo{
		Commits: []ghc.Commit{
			{When: time.Date(2020, 11, 20, 11, 35, 00, 0, tz)},
			{When: time.Date(2020, 11, 20, 11, 40, 00, 0, tz)},
		},
	}
	err = ghc.Merge(dir, repo1, gitHubEmail)
	checkIfError(t, err)
	err = ghc.Merge(dir, repo2, gitHubEmail)
	checkIfError(t, err)

	r, err := git.PlainOpen(dir)
	checkIfError(t, err)

	commits := countCommits(t, r)
	totalCommits := len(repo1.Commits) + len(repo2.Commits)
	if commits != totalCommits {
		t.Fatalf("number of commits expected %v found %v", totalCommits, commits)
	}
}

func TestMerge_NotCountDuplicateDates(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	repo := ghc.Repo{Commits: []ghc.Commit{commit1, commit1, commit1}}

	err = ghc.Merge(dir, repo, gitHubEmail)
	err = ghc.Merge(dir, repo, gitHubEmail) // Merging commits twice. Should be idempotent
	checkIfError(t, err)

	r, err := git.PlainOpen(dir)
	checkIfError(t, err)

	commits := countCommits(t, r)
	if commits != 1 {
		t.Fatalf("number of commits expected %v found %v", 1, commits)
	}
}

func TestMerge_CommitWithProvidedEmail(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	repo := ghc.Repo{Commits: []ghc.Commit{commit1, commit1}}
	err = ghc.Merge(dir, repo, gitHubEmail)
	checkIfError(t, err)

	r, err := git.PlainOpen(dir)
	checkIfError(t, err)

	commits := readCommits(t, r)
	for _, c := range commits {
		if c.Author.Email != gitHubEmail {
			t.Fatal("Committed email does not match expected email")
		}
	}
}

func TestMerge_ErrorWhenDestinationIsFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "gitMerge")
	checkIfError(t, err)
	defer os.RemoveAll(dir)

	testFile := path.Join(dir, "testFile")
	err = ioutil.WriteFile(testFile, []byte(""), 0600)
	checkIfError(t, err)

	repo := ghc.Repo{Commits: []ghc.Commit{commit1, commit1}}
	err = ghc.Merge(testFile, repo, gitHubEmail)
	if !strings.HasSuffix(err.Error(), "not a directory") {
		t.Fatalf("error thrown is not expected. Having: %v", err)
	}
}
