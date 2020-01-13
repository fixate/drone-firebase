package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "firebase deploy plugin"
	app.Usage = "firebase deploy plugin"
	app.Action = run
	app.Version = fmt.Sprintf("0.1.0+%s", build)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:   "token",
			Usage:  "Firebase token",
      EnvVars: []string{"FIREBASE_TOKEN"},
		},
		&cli.StringFlag{
			Name:   "project",
			Usage:  "firebase project id",
      EnvVars: []string{"FIREBASE_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:   "message",
			Usage:  "release message",
      EnvVars: []string{"PLUGIN_MESSAGE"},
		},
		&cli.StringFlag{
			Name:   "targets",
			Usage:  "targets to deploy",
      EnvVars: []string{"PLUGIN_TARGETS,FIREBASE_TARGETS"},
		},
		&cli.StringFlag{
			Name:   "dryrun",
			Usage:  "dry run",
      EnvVars: []string{"PLUGIN_DRY_RUN"},
		},
		&cli.StringFlag{
			Name:   "debug",
			Usage:  "debug",
      EnvVars: []string{"PLUGIN_DEBUG"},
		},
		&cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
      EnvVars: []string{"DRONE_COMMIT_SHA"},
			Value:  "00000000",
		},
		&cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
      EnvVars: []string{"DRONE_COMMIT_REF"},
		},
		&cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
      EnvVars: []string{"DRONE_COMMIT_BRANCH"},
		},
		&cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
      EnvVars: []string{"DRONE_COMMIT_AUTHOR"},
		},
		&cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
      EnvVars: []string{"DRONE_PULL_REQUEST"},
		},
		&cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
      EnvVars: []string{"DRONE_COMMIT_MESSAGE"},
		},
		&cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
      EnvVars: []string{"DRONE_BUILD_EVENT"},
		},
		&cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
      EnvVars: []string{"DRONE_BUILD_NUMBER"},
		},
		&cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
      EnvVars: []string{"DRONE_BUILD_STATUS"},
		},
		&cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
      EnvVars: []string{"DRONE_BUILD_LINK"},
		},
		&cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
      EnvVars: []string{"DRONE_BUILD_STARTED"},
		},
		&cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
      EnvVars: []string{"DRONE_BUILD_CREATED"},
		},
		&cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
      EnvVars: []string{"DRONE_TAG"},
		},
		&cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
      EnvVars: []string{"DRONE_DEPLOY_TO"},
		},
		&cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
      EnvVars: []string{"DRONE_JOB_STARTED"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Commit:   c.String("commit.sha"),
			Ref:      c.String("commit.ref"),
			Branch:   c.String("commit.branch"),
			Author:   c.String("commit.author"),
			Pull:     c.String("commit.pull"),
			Message:  c.String("commit.message"),
			DeployTo: c.String("build.deployTo"),
			Link:     c.String("build.link"),
			Started:  c.Int64("build.started"),
			Created:  c.Int64("build.created"),
		},
		Config: Config{
			Token:     c.String("token"),
			ProjectID: c.String("project"),
			Message:   c.String("message"),
			Targets:   c.String("targets"),
			DryRun:    c.Bool("dryrun"),
			Debug:     c.Bool("debug"),
		},
	}

	return plugin.Exec()
}
