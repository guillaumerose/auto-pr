package main

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type user struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type git struct {
	ctx         context.Context
	client      *github.Client
	user        user
	orgnization string
	repository  string
}

func newGit(ctx context.Context, user user, orgnization string, repository string) *git {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &git{ctx, github.NewClient(tc), user, orgnization, repository}
}

func (g git) getContents(path string) ([]byte, string, error) {
	content, _, _, err := g.client.Repositories.GetContents(g.ctx, g.orgnization, g.repository, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return []byte{}, "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return []byte{}, "", err
	}
	return decoded, *content.SHA, nil
}

func (g git) createFile(commit string, branch string, path string, sha string, content []byte) (string, error) {
	response, _, err := g.client.Repositories.CreateFile(g.ctx, g.user.Username, g.repository, path, &github.RepositoryContentFileOptions{
		Message: github.String(commit + "\n\nSigned-off-by: " + g.user.Name + " <" + g.user.Email + ">"),
		Content: content,
		SHA:     github.String(sha),
		Branch:  github.String(branch),
	})
	if err != nil {
		return "", err
	}
	return response.Content.GetHTMLURL(), nil
}

func (g git) createBranch(newBranch string) error {
	branch, _, err := g.client.Repositories.GetBranch(g.ctx, g.orgnization, g.repository, "master")
	if err != nil {
		return err
	}
	_, _, err = g.client.Git.CreateRef(g.ctx, g.user.Username, g.repository, &github.Reference{
		Ref: github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{
			SHA: branch.Commit.SHA,
		},
	})
	return err
}

func (g git) createPR(title string, branch string) (string, error) {
	pr, _, err := g.client.PullRequests.Create(g.ctx, g.orgnization, g.repository, &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(g.user.Username + ":" + branch),
		Base:  github.String("master"),
	})
	if err != nil {
		return "", err
	}
	return pr.GetHTMLURL(), nil
}
