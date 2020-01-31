package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"io/ioutil"

	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

)

var version = "undefined"

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func createManifest(ext string) (manifestInterface, error) {
	switch ext {
	case "hcl":
		return &manifestHcl{}, nil
	case "json":
		return &manifestJSON{}, nil
	case "yaml":
		return &manifestYaml{}, nil
	default:
		return nil, fmt.Errorf("Unhandled manifest type '%s'", ext)
	}
}

func openManifestFile(file string) (manifest manifestInterface, err error) {
	fmt.Printf("Opening manifest file %s\n", file)

	ext := strings.TrimLeft(path.Ext(file), ".")

	manifest, err = createManifest(ext)
	if err != nil {
		return
	}

	err = manifest.open(file)
	
	return
}

func usage(message string) {
	if message != "" {
		fmt.Println(message)
	}
	flag.Usage()
	os.Exit(1)
}

func main() {
	argGitURL := flag.String(
		"gitUrl",
		os.Getenv("GITOPS_GIT_URL"),
		"URL of git repository. Can also use GITOPS_GIT_URL environment variable")
	argGitUsername := flag.String(
		"gitUsername",
		os.Getenv("GITOPS_GIT_USERNAME"),
		"Username to authenticate with git. Can also use GITOPS_GIT_USERNAME environment variable ")
	argGitPassword := flag.String(
		"gitPassword",
		os.Getenv("GITOPS_GIT_PASSWORD"),
		"Password or token to authenticate with git. Can also use GITOPS_GIT_PASSWORD environment variable")
	argRepoFile := flag.String(
		"repoFile",
		os.Getenv("GITOPS_REPO_FILE"),
		"File in git repo to apply changes to. Can also use GITOPS_REPO_FILE environment variable")
	argValues := flag.String(
		"values",
		os.Getenv("GITOPS_VALUES"),
		"List of variables and coresponding values to update. Variables paths are a list of keys separated with periods. Each variable is separated with a colon. Example '-values=input.one=foo:input.two=bar'. Can also use GITOPS_VALUES environment variable")
	argDryRun := flag.Bool(
		"dry-run",
		false,
		"Turn on to disable making any changes to target repository")
	flagVersion := flag.Bool(
		"version",
		false,
		"Print version")

	flag.Parse()

	if *flagVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	repoPath, err := ioutil.TempDir("", "gitty-up")
	CheckIfError(err)

	filePath := repoPath + "/" + *argRepoFile

	if *argGitURL == "" {
		usage("ERROR: Git URL is required!")
	}

	if *argGitUsername == "" || argGitUsername == nil {
		usage("ERROR: Git username is required!")
	}

	if *argGitPassword == "" || argGitPassword == nil {
		usage("ERROR: Git password is required!")
	}

	if *argRepoFile == "" {
		usage("ERROR: File is required!")
	}

	if *argValues == "" {
		usage("ERROR: Values are required!")
	}
	valuePaths, err := parseValues(*argValues)
	if err != nil {
		usage("ERROR: Could not parse values")
	}

	gitAuth := &http.BasicAuth{
		Username: *argGitUsername,
		Password: *argGitPassword,
	}

	repo, err := gitClone(*argGitURL, gitAuth, repoPath)
	CheckIfError(err)

	branch, worktree, err := gitBranch(repo)
	CheckIfError(err)

	manifest, err := openManifestFile(filePath)
	CheckIfError(err)

	for _, value := range valuePaths {
		err = manifest.setValue(value.path, value.value)
		CheckIfError(err)
	}

	err = manifest.save()
	CheckIfError(err)

	err = gitCommit(worktree, *argRepoFile)
	CheckIfError(err)

	if *argDryRun == false {
		err = gitPush(repo, gitAuth)
		CheckIfError(err)
	}

	gitURLParts, err := url.Parse(*argGitURL)
	CheckIfError(err)

	if gitURLParts.Host == "github.com" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *argGitPassword})
		tokenClient := oauth2.NewClient(context.Background(), tokenSource)

		pathParts := strings.Split(strings.TrimSuffix(gitURLParts.Path, ".git"), "/")

		if *argDryRun == false {
			pullRequest, err := githubPullRequest(tokenClient, pathParts[1], pathParts[2], branch)
			CheckIfError(err)

			fmt.Printf("Pull Request created: %s\n", pullRequest.GetHTMLURL())
		}
	}
}
