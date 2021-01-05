package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func githubPullRequest(httpClient *http.Client, org string, repo string, base string, branch *plumbing.Reference) (pullRequest *github.PullRequest, err error) {
	fmt.Printf("Create pull request for branch %s\n", branch.Name().Short())
	client := github.NewClient(httpClient)

	newPR := &github.NewPullRequest{
		Title:               github.String("GitOps"),
		Head:                github.String(branch.Name().Short()),
		Base:                github.String(base),
		Body:                github.String(createGitMessage("pull request")),
		MaintainerCanModify: github.Bool(true),
	}

	pullRequest, _, err = client.PullRequests.Create(context.Background(), org, repo, newPR)

	return
}
