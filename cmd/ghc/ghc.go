package main

import (
	"errors"
	"fgh/internal/ghc"
	"flag"
	"fmt"
	"os"
	"strings"
)

var errSrcRequired = errors.New("-src parameter is required")
var errDstRequired = errors.New("-dst parameter is required")
var errEmailsRequired = errors.New("-emails parameter is required")
var errGitHubRequired = errors.New("-github-email parameter is required")

var src string
var dst string
var emails []string
var gitHubEmail string

func main() {
	err := parseArgs()
	exitIfError(err)

	repositories, err := ghc.ScanDir(src)
	exitIfError(err)

	fmt.Println("Number of repositories found: ", len(repositories))

	extract := ghc.NewExtract(emails)
	for _, rp := range repositories {
		r, err := extract.Repository(rp)
		exitIfError(err)

		fmt.Printf("%v commits found in %v\n", len(r.Commits), r.Path)

		err = ghc.Merge(dst, r, gitHubEmail)
		exitIfError(err)
	}
}

func exitIfError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error! %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() error {
	flag.StringVar(&src, "src", "", "directory recursively scanned for GIT repositories")
	flag.StringVar(&dst, "dst", "", "the contributions repository")
	flag.StringVar(&gitHubEmail, "github-email", "", "the same email as your GitHub account")
	emailsLst := ""
	flag.StringVar(&emailsLst, "emails", "", "comma separated list of emails")

	flag.Parse()

	if src == "" {
		return errSrcRequired
	}
	if dst == "" {
		return errDstRequired
	}
	if gitHubEmail == "" {
		return errGitHubRequired
	}

	emails = strings.Split(emailsLst, ",")
	if len(emails) == 0 {
		return errEmailsRequired
	}

	return nil
}
