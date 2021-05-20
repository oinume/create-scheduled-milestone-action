package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v35/github"
)

func Test_app_run(t *testing.T) {
	type wants struct {
		status int
		out    string
		err    string
	}

	tests := map[string]struct {
		envs    map[string]string
		handler http.Handler
		wants   wants
	}{
		"ok": {
			envs: map[string]string{
				"GITHUB_REPOSITORY": "oinume/create-scheduled-milestone-action",
				"INPUT_TITLE":       "v1.0.0",
				"INPUT_STATE":       "open",
				"INPUT_DESCRIPTION": "v1.0.0 release",
				"INPUT_DUE_ON":      "2021-05-10T21:43:54+09:00",
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				body := `{"number": 111}`
				_, _ = fmt.Fprintln(w, body)
			}),
			wants: wants{
				status: 0,
				out:    "::set-output name=number::111\n",
				err:    "",
			},
		},
		"error_invalid_github_repository": {
			envs: map[string]string{
				"GITHUB_REPOSITORY": "invalid",
				"INPUT_TITLE":       "v1.0.0",
				"INPUT_STATE":       "open",
				"INPUT_DESCRIPTION": "v1.0.0 release",
				"INPUT_DUE_ON":      "2021-05-10T21:43:54+09:00",
			},
			wants: wants{
				status: 1,
				out:    "",
				err:    "invalid repository format\n",
			},
		},
		"error_empty_title": {
			envs: map[string]string{
				"GITHUB_REPOSITORY": "oinume/create-scheduled-milestone-action",
				"INPUT_TITLE":       "",
				"INPUT_STATE":       "open",
				"INPUT_DESCRIPTION": "v1.0.0 release",
				"INPUT_DUE_ON":      "2021-05-10T21:43:54+09:00",
			},
			wants: wants{
				status: 1,
				out:    "",
				err:    "'title' is required\n",
			},
		},
	}

	for name, tt := range tests {
		for k, v := range tt.envs {
			_ = os.Setenv(k, v)
		}

		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()
			githubClient := newFakeGitHubClient(t, ts.URL+"/")

			var outStream, errStream bytes.Buffer
			a := newApp(githubClient, &outStream, &errStream)
			ctx := context.Background()
			if got := a.run(ctx); got != tt.wants.status {
				t.Fatalf("run() status: got = %v, want = %v", got, tt.wants.status)
			}
			if got := outStream.String(); got != tt.wants.out {
				t.Errorf("run() out: got = %v, want = %v", got, tt.wants.out)
			}
			if got := errStream.String(); got != tt.wants.err {
				t.Errorf("run() err: got = %v, want = %v", got, tt.wants.err)
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
				t.Fatalf("newMilestone() error = %v, wantErrStream %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMilestone() got = %v, wantStatus = %v", got, tt.want)
			}
		})
	}
}

func newFakeGitHubClient(t *testing.T, baseURL string) *github.Client {
	t.Helper()
	c := newGitHubClient(context.Background(), "")
	u, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}
	c.BaseURL = u
	return c
}
