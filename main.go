package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	title := os.Getenv("INPUT_TITLE")
	state := os.Getenv("INPUT_STATE")
	description := os.Getenv("INPUT_DESCRIPTION")
	dueOn := os.Getenv("INPUT_DUE_ON")
	m, err := newMilestone(githubRepository, title, state, description, dueOn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	client := newGitHubClient(ctx, githubToken)
	a := &app{
		githubClient: client,
	}

	if err := a.run(ctx, m); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

type app struct {
	githubClient *github.Client
	outStream, errStream io.Writer
}

type milestone struct {
	owner       string
	repo        string
	title       string
	state       string
	description string
	dueOn       time.Time
}

func (m *milestone) toGitHub() *github.Milestone {
	ghm := &github.Milestone{
		Title:       &m.title,
	}
	if m.state != "" {
		ghm.State = &m.state
	}
	if m.description != "" {
		ghm.Description = &m.description
	}
	if !m.dueOn.IsZero() {
		ghm.DueOn = &m.dueOn
	}
	return ghm
}

func newMilestone(repository, title, state, description, dueOn string) (*milestone, error) {
	r := strings.Split(repository, "/")
	if len(r) != 2 {
		return nil, errors.New("invalid repository format")
	}
	if title == "" {
		return nil, errors.New("'title' is required")
	}

	var dueOnTime time.Time
	if dueOn == "" {
		dueOnTime = time.Time{}
	} else {
		t, err := time.Parse(time.RFC3339, dueOn)
		if err != nil {
			return nil, fmt.Errorf("time.Parse failed: %v", err)
		}
		dueOnTime = t
	}

	return &milestone{
		owner:       r[0],
		repo:        r[1],
		title:       title,
		state:       state,
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

// run creates a new milestone
func (c *app) run(ctx context.Context, m *milestone) error {
	milestone, _, err := c.githubClient.Issues.CreateMilestone(
		ctx,
		m.owner,
		m.repo,
		m.toGitHub(),
	)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "::set-output name=number::%d\n", milestone.GetNumber())
	return nil
}
