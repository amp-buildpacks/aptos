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

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	if _, ok, err := pr.Resolve(PlanEntryAptos); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve Aptos plan entry\n%w", err)
	} else if ok {
		cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
		}

		dc, err := libpak.NewDependencyCache(context)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
		}
		dc.Logger = b.Logger

		dr, err := libpak.NewDependencyResolver(context)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
		}

		// install aptos
		// update dependency.id from metadata.dependencies for buildpack.toml
		v, _ := cr.Resolve("BP_APTOS_VERSION")
		dependency, err := dr.Resolve(PlanEntryAptos, v)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		aptosLayer := NewAptos(dependency, dc)
		aptosLayer.Logger = b.Logger

		result.Processes, err = aptosLayer.BuildProcessTypes(cr, context.Application)
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to build list of process types\n%w", err)
		}
		result.Layers = append(result.Layers, aptosLayer)
	}

	return result, nil
}
