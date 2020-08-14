package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func gitClone(url string, auth transport.AuthMethod, repoPath string) (*git.Repository, error) {
	fmt.Printf("Cloning git repo %s with %s\n", url, auth.String())
	return git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: auth,
		URL:  url,
	})
}

func gitBranch(repo *git.Repository) (branch *plumbing.Reference, worktree *git.Worktree, err error) {
	name := fmt.Sprintf("gitops-%x", time.Now().Unix())
	fmt.Printf("Create branch %s\n", name)

	head, err := repo.Head()
	if err != nil {
		return
	}

	branch = plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), head.Hash())

	err = repo.Storer.SetReference(branch)

	worktree, err = repo.Worktree()

	worktree.Checkout(&git.CheckoutOptions{
		Branch: branch.Name(),
	})

	return
}

func gitCommit(worktree *git.Worktree, file string) (err error) {
	fmt.Printf("Add changes to repo %s\n", file)
	_, err = worktree.Add(file)
	if err != nil {
		return
	}

	val := createGitMessage("commit")

	fmt.Println("Commiting changes")
	worktree.Commit(fmt.Sprintf("GitOps: Update %s\n%s", file, val), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "GitOps Automation",
			Email: "gitops@liatr.io",
			When:  time.Now(),
		},
	})

	return
}

func gitPush(repo *git.Repository, auth transport.AuthMethod) (err error) {
	fmt.Println("Pushing changes")
	err = repo.Push(&git.PushOptions{
		Auth: auth,
	})
	return
}

func createGitMessage(gitAction string) (pullRequestMessage string) {
	val, exists := os.LookupEnv("GIT_URL")
	if exists {
		source := strings.Split(strings.SplitN(val, "/", 4)[3], ".")
		message := "This " + gitAction + " was automatically generated from the pipeline of the repo " + source[0] + "\n\n"

		gitAuth := &http.BasicAuth{
			Username: os.Getenv("GITOPS_GIT_USERNAME"),
			Password: os.Getenv("GITOPS_GIT_PASSWORD"),
		}

		sourceRepo, err := gitClone(val, gitAuth, "/home/jenkins/tempSourceRepo/")
		CheckIfError(err)

		ref, err := sourceRepo.Head()
		CheckIfError(err)

		cIter, err := sourceRepo.Log(&git.LogOptions{From: ref.Hash()})
		CheckIfError(err)

		os.RemoveAll("/home/jenkins/tempSourceRepo/")

		c, err := cIter.Next()
		return message + c.String()
	}
	return "This " + gitAction + " was automatically generated https://github.com/liatrio/builder-images/tree/master/builder-image-gitops"
}
