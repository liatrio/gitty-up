package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

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
	case "yaml", "yml":
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
	gitURL := flag.String(
		"gitUrl",
		os.Getenv("GITOPS_GIT_URL"),
		"URL of git repository. Can also use GITOPS_GIT_URL environment variable")
	gitUsername := flag.String(
		"gitUsername",
		os.Getenv("GITOPS_GIT_USERNAME"),
		"Username to authenticate with git. Can also useGITOPS_GIT_USERNAME environment variable ")
	gitPassword := flag.String(
		"gitPassword",
		os.Getenv("GITOPS_GIT_PASSWORD"),
		"Password or token to authenticate with git. Can also use GITOPS_GIT_PASSWORD environment variable")
	repoFile := flag.String(
		"repoFile",
		os.Getenv("GITOPS_REPO_FILE"),
		"File in git repo to apply changes to. Can also use GITOPS_REPO_FILE environment variable")
	values := flag.String(
		"values",
		os.Getenv("GITOPS_VALUES"),
		"List of variables and coresponding values to update. Variables paths are a list of keys separated with periods. Each variable is separated with a colon. Example '-values=input.one=foo:input.two=bar'. Can also use GITOPS_VALUES environment variable")
	dryrun := flag.Bool(
		"dry-run",
		false,
		"Turn on to disable making any changes to target repository")

	flag.Parse()

	repoPath, err := ioutil.TempDir("", "gitty-up")
	CheckIfError(err)

	filePath := repoPath + "/" + *repoFile

	fmt.Println("Start GitOps")

	if *gitURL == "" {
		usage("ERROR: Git URL is required!")
	}

	if *gitUsername == "" || gitUsername == nil {
		usage("ERROR: Git username is required!")
	}

	if *gitPassword == "" || gitPassword == nil {
		usage("ERROR: Git password is required!")
	}

	if *repoFile == "" {
		usage("ERROR: File is required!")
	}

	if *values == "" {
		usage("ERROR: Values are required!")
	}
	valuePaths, err := parseValues(*values)
	if err != nil {
		usage("ERROR: Could not parse values")
	}

	gitAuth := &http.BasicAuth{
		Username: *gitUsername,
		Password: *gitPassword,
	}

	repo, err := gitClone(*gitURL, gitAuth, repoPath)
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

	err = gitCommit(worktree, *repoFile)
	CheckIfError(err)

	if *dryrun == false {
		err = gitPush(repo, gitAuth)
		CheckIfError(err)
	}

	gitURLParts, err := url.Parse(*gitURL)
	CheckIfError(err)

	if gitURLParts.Host == "github.com" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *gitPassword})
		tokenClient := oauth2.NewClient(context.Background(), tokenSource)

		pathParts := strings.Split(strings.TrimSuffix(gitURLParts.Path, ".git"), "/")

		if *dryrun == false {
			pullRequest, err := githubPullRequest(tokenClient, pathParts[1], pathParts[2], branch)
			CheckIfError(err)

			fmt.Printf("Pull Request created: %s\n", pullRequest.GetHTMLURL())
		}
	}
}
