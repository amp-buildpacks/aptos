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
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/buildpacks/libcnb"
)

const (
	PlanEntryAptos = "aptos"
)

type Detect struct {
}

type Package struct {
	Dependencies map[string]*Dependency
}

type Dependency struct {
	Git    string
	Rev    string
	Subdir string
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	found, err := d.aptosProject(context.Application.Path)
	if err != nil {
		return libcnb.DetectResult{Pass: false}, fmt.Errorf("unable to detect Aptos requirements\n%w", err)
	}

	if !found {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: PlanEntryAptos},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: PlanEntryAptos},
				},
			},
		},
	}, nil
}

func (d Detect) aptosProject(appDir string) (bool, error) {
	projectName := filepath.Join(appDir, "Move.toml")
	_, err := os.Stat(projectName)
	if os.IsNotExist(err) || err != nil {
		return false, fmt.Errorf("unable to determine if %s exists\n%w", projectName, err)
	}

	var packageConfig Package
	if _, err := toml.DecodeFile(projectName, &packageConfig); err != nil {
		return false, fmt.Errorf("unable to parsing %s\n%w", projectName, err)
	}

	if _, ok := packageConfig.Dependencies["AptosFramework"]; !ok {
		return false, fmt.Errorf("unable to determine if Aptos project")
	}

	buildDirectory := filepath.Join(appDir, ".")
	extension := ".move"
	if err := existsFilesWithExtension(buildDirectory, extension); err != nil {
		return false, fmt.Errorf("unable to determine if '%s' exists\n%w", extension, err)
	}
	return true, nil
}

func existsFilesWithExtension(directory, extension string) error {
	var found bool
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has the specified extension.
		if !info.IsDir() && filepath.Ext(path) == extension {
			found = true
			return nil
		}
		return nil
	})

	if !found {
		return fmt.Errorf("no files with extension '%s' found", extension)
	}
	return err
}
