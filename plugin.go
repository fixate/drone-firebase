// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

		"github.com/aymerick/raymond"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Ref      string
		Branch   string
		Author   string
		Pull     string
		Message  string
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
	}

	Config struct {
		Token     string
		ProjectID string // optional.
		Message   string // optional.
		Targets   string // optional.
		DryRun    bool   // optional.
		Debug     bool   // optional.
	}

	Plugin struct {
		Config Config
		Build  Build
		Repo   Repo
	}

	Workspace struct {
		Path string `json:"path"`
	}
)

func (p Plugin) Exec() error {
	if err := p.doDeployment(); err != nil {
		fmt.Printf("Firebase: Error in deployment: %s\n", err)
		os.Exit(1)
	}
}

func (p Plugin) doDeployment() error {
	if p.Config.ProjectID != "" {
		if err := p.useProject(); err != nil {
			return err
		}
	}

	if err := p.deploy(); err != nil {
		return err
	}

	return nil
}

func (p Plugin) getEnvironment(oldEnv []string) []string {
	var env []string
	for _, v := range oldEnv {
		if !strings.HasPrefix(v, "DEBUG=") && !strings.HasPrefix(v, "FIREBASE_TOKEN=") {
			env = append(env, v)
		}
	}
	env = append(env, fmt.Sprintf("FIREBASE_TOKEN=%s", p.Config.Token))
	if p.Config.Debug {
		env = append(env, fmt.Sprintf("DEBUG=%s", "true"))
	}
	return env
}

// Sets the active project using the
// $ firebase use ... command
func (p Plugin) useProject() error {
	var args []string
	args = append(args, "use")

	if p.Config.ProjectID != "" {
		args = append(args, p.Config.ProjectID)
	}

	return p.runFirebaseCommand(args)
}

func render(template string, payload interface{})  (string, error) {
	const str, err := raymond.Render(template, payload)
	return strings.Trim(str, " \n"), err
}

// buildDeploy runs a deploy command:
// $ firebase deploy \
//   [--only ...] \
//   [--message ...]
func (p Plugin) deploy() error {
	var args []string
	args = append(args, "deploy")

	if p.Config.Targets != "" {
		args = append(args, "--only")
		args = append(args, p.Config.Targets)
	}

	if p.Config.Message != "" {
		args = append(args, "--message")
		args = append(args, fmt.Sprintf("\"%s\"", render(p.Config.Message, p))
	}

	return p.runFirebaseCommand(args);
}

// execute sets the stdout and stderr of the command to be the default, traces
// the command to be executed and returns the result of the command execution.
func (p Plugin) execute(cmd *exec.Cmd) error {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if p.Config.DryRun || p.Config.Debug {
		fmt.Println("$", strings.Join(cmd.Args, " "))
	}
	if p.Config.DryRun {
		return nil
	}
	return cmd.Run()
}

func (p Plugin) runFirebaseCommand(args []string) error {
	cmd := exec.Command("firebase", args...)
	cmd.Env = p.getEnvironment(os.Environ())
	return p.execute(cmd)
}
