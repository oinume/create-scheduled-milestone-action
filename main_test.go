package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

//func Test_app_run(t *testing.T) {
//	type fields struct {
//		githubClient *github.Client
//	}
//	type args struct {
//		ctx context.Context
//		m   *milestone
//	}
//
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    int
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &app{
//				githubClient: tt.fields.githubClient,
//			}
//			got, err := c.run(tt.args.ctx, tt.args.m)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("run() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

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
