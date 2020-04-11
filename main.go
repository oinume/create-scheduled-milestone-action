package main

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"os"
	"time"

	"github.com/google/go-github/v31/github"
)

func main() {
	//myInput := os.Getenv("INPUT_MYINPUT")
	//
	//output := fmt.Sprintf("Hello %s", myInput)
	//
	//fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
	ctx := context.Background()
	githubToken := os.Getenv("GITHUB_TOKEN")
	client := newGitHubClient(ctx, githubToken)

	a := &app{
		githubClient: client,
	}

	title := os.Getenv("INPUT_TITLE")
	description := os.Getenv("INPUT_DESCRIPTION")
	dueOn := os.Getenv("INPUT_DUE_ON")
	if err := a.run(title, description, dueOn); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

type app struct {
	githubClient *github.Client
	// TODO: outStream, errStream
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (c *app) run(title, description string, dueOn string) error {
	if title == "" {
		return errors.New("'title' is required")
	}

	_, err := c.validateDueOn(dueOn)
	if err != nil {
		return err
	}

	return nil
}

func (c *app) validateDueOn(dueOn string) (time.Time, error) {
	return time.Time{}, nil
}
