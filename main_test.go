package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v31/github"
)

func Test_app_run_ok(t *testing.T) {
	type fields struct {
		outStream    io.Writer
		errStream    io.Writer
	}
	tests := map[string]struct {
		envs          map[string]string
		handler       http.Handler
		fields        fields
		wantStatus    int
		wantOutStream string
		wantErrStream string
	}{
		"ok": {
			envs: map[string]string{
				"GITHUB_REPOSITORY": "oinume/create-scheduled-milestone-action",
				"INPUT_TITLE": "v1.0.0",
				"INPUT_STATE": "open",
				"INPUT_DESCRIPTION": "v1.0.0 release",
				"INPUT_DUE_ON": "2021-05-10T21:43:54+09:00",
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				body := `{"number": 111}`
				_, _ = fmt.Fprintln(w, body)
			}),
			fields: fields{
				outStream: os.Stdout,
				errStream: os.Stderr,
			},
			wantStatus:    0,
			wantOutStream: "::set-output name=number::111\n",
			wantErrStream: "",
		},
	}

	for name, tt := range tests {
		for k, v := range tt.envs {
			_ = os.Setenv(k, v)
		}
		ts := httptest.NewServer(tt.handler)
		defer ts.Close()
		githubClient := newFakeGitHubClient(t, ts.URL + "/")

		t.Run(name, func(t *testing.T) {
			var outStream, errStream bytes.Buffer
			a := newApp(githubClient, &outStream, &errStream)
			ctx := context.Background()
			if got := a.run(ctx); got != tt.wantStatus {
				t.Fatalf("run() status: got = %v, want = %v", got, tt.wantStatus)
			}
			if got := outStream.String(); got != tt.wantOutStream {
				t.Errorf("run() outStream: got = %v, want = %v", got, tt.wantOutStream)
			}
			if got := errStream.String(); got != tt.wantErrStream {
				t.Errorf("run() errStream: got = %v, want = %v", got, tt.wantErrStream)
			}
		})
	}
}

func Test_app_createMilestone(t *testing.T) {
	type args struct {
		ctx context.Context
		m   *milestone
	}
	wantNumber := 111

	tests := map[string]struct {
		args    args
		handler http.Handler
		want    *github.Milestone
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
			want: &github.Milestone{
				Number: &wantNumber,
			},
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

			// TODO: Use outStream
			got, err := c.createMilestone(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Fatalf("createMilestone(): error = %v, wantErrStream = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMilestone(): got = %v, wantStatus = %v", got, tt.want)
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
	c := github.NewClient(nil)
	u, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("url.Parse failed: %v", err)
	}
	c.BaseURL = u
	return c
}
