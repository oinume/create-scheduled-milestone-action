package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	githubToken := os.Getenv("GITHUB_TOKEN")
	client := newGitHubClient(ctx, githubToken)

	a := &app{
		githubClient: client,
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	title := os.Getenv("INPUT_TITLE")
	description := os.Getenv("INPUT_DESCRIPTION")
	dueOn := os.Getenv("INPUT_DUE_ON")
	m, err := newMilestone(githubRepository, title, description, dueOn)
	if err != nil {
		os.Exit(1)
	}

	number, err := a.run(ctx, m)
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("::set-output name=myOutput::%d\n", number)
	os.Exit(0)
}

type app struct {
	githubClient *github.Client
	// TODO: outStream, errStream
}

type milestone struct {
	owner       string
	repo        string
	title       string
	description string
	dueOn       time.Time
}

func (m *milestone) toGitHub() *github.Milestone {
	ghm := &github.Milestone{
		Title:       &m.title,
		Description: &m.description,
	}
	if !m.dueOn.IsZero() {
		ghm.DueOn = &m.dueOn
	}
	return ghm
}

func newMilestone(githubRepository, title, description, dueOn string) (*milestone, error) {
	r := strings.Split(githubRepository, "/")
	if len(r) != 2 {
		return nil, errors.New("hoge")
	}
	if title == "" {
		return nil, errors.New("'title' is required")
	}

	dueOnTime, err := time.Parse(time.RFC3339, dueOn)
	if err != nil {
		return nil, fmt.Errorf("time.Parse failed: %v", err)
	}

	return &milestone{
		owner:       r[0],
		repo:        r[1],
		title:       title,
		description: description,
		dueOn:       dueOnTime,
	}, nil
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (c *app) run(ctx context.Context, m *milestone) (int, error) {
	milestone, _, err := c.githubClient.Issues.CreateMilestone(
		ctx,
		m.owner,
		m.repo,
		m.toGitHub(),
	)
	if err != nil {
		return 0, err
	}
	return milestone.GetNumber(), nil
}
