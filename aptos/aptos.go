// Copyright (c) The Amphitheatre Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aptos

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/crush"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Aptos struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
	Executor         effect.Executor
}

func NewAptos(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) Aptos {
	contributor := libpak.NewDependencyLayerContributor(dependency, cache, libcnb.LayerTypes{
		Build:  true,
		Cache:  true,
		Launch: true,
	})
	return Aptos{
		LayerContributor: contributor,
		Executor:         effect.NewExecutor(),
	}
}

func (r Aptos) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	r.LayerContributor.Logger = r.Logger
	return r.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		bin := filepath.Join(layer.Path, "bin")

		r.Logger.Bodyf("Expanding %s to %s", artifact.Name(), bin)
		if err := crush.Extract(artifact, bin, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand %s\n%w", artifact.Name(), err)
		}

		// Must be set to executable
		file := filepath.Join(bin, PlanEntryAptos)
		r.Logger.Bodyf("Setting %s as executable", file)
		if err := os.Chmod(file, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to chmod %s\n%w", file, err)
		}

		// Must be set to PATH
		r.Logger.Bodyf("Setting %s in PATH", bin)
		if err := os.Setenv("PATH", sherpa.AppendToEnvVar("PATH", ":", bin)); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to set $PATH\n%w", err)
		}

		// get version
		buf, err := r.Execute(PlanEntryAptos, []string{"--version"})
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to get %s version\n%w", PlanEntryAptos, err)
		}
		version := strings.Split(strings.TrimSpace(buf.String()), " ")[1]
		r.Logger.Bodyf("Checking %s version: %s", PlanEntryAptos, version)

		// set MOVE_HOME
		moveHome := "/workspace"
		r.Logger.Bodyf("Setting MOVE_HOME = %s", moveHome)
		if err := os.Setenv("MOVE_HOME", moveHome); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to set MOVE_HOME\n%w", err)
		}

		// compile contract
		args := []string{"move", "compile"}
		r.Logger.Bodyf("Compiling contracts by '%s %s'", PlanEntryAptos, args)
		if _, err := r.Execute(PlanEntryAptos, args); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to compile contract\n%w", err)
		}

		layer.LaunchEnvironment.Append("PATH", ":", bin)
		layer.LaunchEnvironment.Default("MOVE_HOME", moveHome)
		return layer, nil
	})
}

func (r Aptos) Execute(command string, args []string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if err := r.Executor.Execute(effect.Execution{
		Command: command,
		Args:    args,
		Stdout:  buf,
		Stderr:  buf,
	}); err != nil {
		return buf, fmt.Errorf("%s: %w", buf.String(), err)
	}
	return buf, nil
}

func (r Aptos) BuildProcessTypes(cr libpak.ConfigurationResolver, app libcnb.Application) ([]libcnb.Process, error) {
	processes := []libcnb.Process{}

	enableDeploy := cr.ResolveBool("BP_ENABLE_APTOS_PROCESS")
	if enableDeploy {
		processes = append(processes, libcnb.Process{
			Type:      PlanEntryAptos,
			Command:   PlanEntryAptos,
			Arguments: []string{"move", "publish", "--skip-fetch-latest-git-deps", "--assume-yes"},
			Default:   true,
		})
	}
	return processes, nil
}

func (r Aptos) Name() string {
	return r.LayerContributor.LayerName()
}
