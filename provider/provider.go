// Copyright 2016-2023, Pulumi Corporation.
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

package provider

import (
	"math/rand"
	"time"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "teamcity"

func Provider() p.Provider {
	prv := infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			DisplayName: "Teamcity",
			License:     "Apache-2.0",
			Repository:  "https://github.com/oss4u/pulumi-teamcity-native",
			Publisher:   "Oss4u",
			LanguageMap: map[string]any{
				"nodejs": map[string]any{
					"packageName": "@oss4u/teamcity",
				},
				"go": map[string]any{
					"generateResourceContainerTypes": true,
					"importBasePath":                 "github.com/oss4u/pulumi-teamcity-native/sdk/go/teamcity",
				},
				"csharp": map[string]any{
					"rootNamespace": "Oss4u",
				},
			},
			PluginDownloadURL: "github://api.github.com/oss4u/pulumi-teamcity-native",
		},
		// Resources: []infer.InferredResource{
		// 	infer.Resource[unbound.HostAliasOverride, unbound.HostAliasOverrideArgs, unbound.HostAliasOverrideState](),
		// 	infer.Resource[unbound.HostOverride, unbound.HostOverrideArgs, unbound.HostOverrideState](),
		// },
		Resources: []infer.InferredResource{
			infer.Resource[Random, RandomArgs, RandomState](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
		// Config: infer.Config[*config.Config](),
	})
	prv.DiffConfig = diff()
	return prv

	// We tell the provider what resources it needs to support.
	// In this case, a single custom resource.
	// return infer.Provider(infer.Options{
	// 	Resources: []infer.InferredResource{
	// 		infer.Resource[Random, RandomArgs, RandomState](),
	// 	},
	// 	ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
	// 		"provider": "index",
	// 	},
	// })
}

func diff() func(ctx p.Context, req p.DiffRequest) (p.DiffResponse, error) {
	return func(ctx p.Context, req p.DiffRequest) (p.DiffResponse, error) {
		return p.DiffResponse{
			DeleteBeforeReplace: false,
			HasChanges:          false,
			DetailedDiff:        nil,
		}, nil
	}
}

// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type Random struct{}

// Each resource has an input struct, defining what arguments it accepts.
type RandomArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but it's generally a
	// good idea.
	Length int `pulumi:"length"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type RandomState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	RandomArgs
	// Here we define a required output called result.
	Result string `pulumi:"result"`
}

// All resources must implement Create at a minimum.
func (Random) Create(ctx p.Context, name string, input RandomArgs, preview bool) (string, RandomState, error) {
	state := RandomState{RandomArgs: input}
	if preview {
		return name, state, nil
	}
	state.Result = makeRandom(input.Length)
	return name, state, nil
}

func makeRandom(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
