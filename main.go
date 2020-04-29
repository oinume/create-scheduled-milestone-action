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
	githubToken := os.Getenv("GITHUB_TOKEN")
	status := newApp(
		newGitHubClient(ctx, githubToken),
		os.Stdout,
		os.Stderr,
	).run(ctx)
	os.Exit(status)
}

type app struct {
	githubClient *github.Client
	outStream, errStream io.Writer
}

func newApp(githubClient *github.Client, outStream, errStream io.Writer) *app {
	return &app{
		githubClient: githubClient,
		outStream: outStream,
		errStream: errStream,
	}
}

type milestone struct {
	owner       string
	repo        string
	title       string
	state       string
	description string
	dueOn       time.Time
}

func newMilestone(repository, title, state, description, dueOn string) (*milestone, error) {
	r := strings.Split(repository, "/")
	if len(r) != 2 {
		return nil, errors.New("invalid repository format")
	}
	if title == "" {
		return nil, errors.New("'title' is required")
	}
	if !(state == "open" || state == "closed") {
		return nil, errors.New("'state' must be open or closed")
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

func (a *app) run(ctx context.Context) int {
	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	title := os.Getenv("INPUT_TITLE")
	state := os.Getenv("INPUT_STATE")
	description := os.Getenv("INPUT_DESCRIPTION")
	dueOn := os.Getenv("INPUT_DUE_ON")
	m, err := newMilestone(githubRepository, title, state, description, dueOn)
	if err != nil {
		fmt.Fprintf(a.errStream, "%v\n", err)
		return 1
	}

	created, err := a.createMilestone(ctx, m)
	if err != nil {
		fmt.Fprintf(a.errStream, "%v\n", err)
		return 1
	}
	fmt.Fprintf(a.outStream, "::set-output name=number::%d\n", created.GetNumber())

	return 0
}

// createMilestone creates a new milestone
func (a *app) createMilestone(ctx context.Context, m *milestone) (*github.Milestone, error) {
	created, _, err := a.githubClient.Issues.CreateMilestone(
		ctx,
		m.owner,
		m.repo,
		m.toGitHub(),
	)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
