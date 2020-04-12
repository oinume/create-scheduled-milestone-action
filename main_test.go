package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v31/github"
)

func Test_app_run(t *testing.T) {
	type args struct {
		ctx context.Context
		m   *milestone
	}

	tests := map[string]struct {
		args    args
		handler http.Handler
		want    int
		wantErr bool
	}{
		"status created": {
			args: args{
				ctx: context.Background(),
				m: &milestone{
					owner:       "oinume",
					repo:        "create-milestone-action",
					title:       "v1.0.0",
					state:       "open",
					description: "v1.0.0 release",
					dueOn:       time.Date(2012, 10, 9, 23, 39, 1, 0, time.UTC),
				},
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				body := `{"number": 111}`
				_, _ = fmt.Fprintln(w, body)
			}),
			want: 111,
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()
			githubClient := newFakeGitHubClient(t, ts.URL + "/")
			c := &app{
				githubClient: githubClient,
			}

			got, err := c.run(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Fatalf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("run() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_newMilestone(t *testing.T) {
	type args struct {
		githubRepository string
		title            string
		state            string
		description      string
		dueOn            string
	}
	tests := map[string]struct {
		args    args
		want    *milestone
		wantErr bool
		err     error
	}{
		"ok": {
			args: args{
				githubRepository: "oinume/create-milestone-action",
				title:            "v1.0.0",
				state:            "open",
				description:      "v1.0.0 release",
				dueOn:            "2012-10-09T23:39:01Z",
			},
			want: &milestone{
				owner:       "oinume",
				repo:        "create-milestone-action",
				title:       "v1.0.0",
				state:       "open",
				description: "v1.0.0 release",
				dueOn:       time.Date(2012, 10, 9, 23, 39, 1, 0, time.UTC),
			},
			wantErr: false,
		},
		"invalid repository format": {
			args: args{
				githubRepository: "oinume$create-milestone-action",
				title:            "",
				state:            "",
				description:      "",
				dueOn:            "",
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("invalid repository format"),
		},
		"empty title error": {
			args: args{
				githubRepository: "oinume/create-milestone-action",
				title:            "",
				state:            "",
				description:      "",
				dueOn:            "",
			},
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("'title' is required"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := newMilestone(tt.args.githubRepository, tt.args.title, tt.args.state, tt.args.description, tt.args.dueOn)
			if (err != nil) != tt.wantErr {
				t.Fatalf("newMilestone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMilestone() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func newFakeGitHubClient(t *testing.T, baseURL string) *github.Client {
	t.Helper()
	c := github.NewClient(nil)
	u, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}
	c.BaseURL = u
	return c
}