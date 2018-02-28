// Copyright 2016 Palantir Technologies, Inc.
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

package golicense

import (
	"fmt"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

type ProjectParam struct {
	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []CustomHeaderParam

	// Exclude matches the files and directories that should be excluded from consideration for verifying or
	// applying licenses.
	Exclude matcher.Matcher
}

type ProjectConfig struct {
	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string `yaml:"header"`

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []CustomHeaderConfig `yaml:"custom-headers"`

	// Exclude matches the files and directories that should be excluded from consideration for verifying or
	// applying licenses.
	Exclude matcher.NamesPathsCfg `yaml:"exclude"`
}

func (c *ProjectConfig) ToParam() (ProjectParam, error) {
	customHeaders := make([]CustomHeaderParam, len(c.CustomHeaders))
	for i, v := range c.CustomHeaders {
		customHeaders[i] = v.ToParam()
	}
	if err := validateCustomHeaderParams(customHeaders); err != nil {
		return ProjectParam{}, err
	}
	return ProjectParam{
		Header:        c.Header,
		CustomHeaders: customHeaders,
		Exclude:       c.Exclude.Matcher(),
	}, nil
}

func validateCustomHeaderParams(headerParams []CustomHeaderParam) error {
	var emptyNameParams []CustomHeaderParam
	nameToParams := make(map[string][]CustomHeaderParam)

	for _, v := range headerParams {
		if v.Name == "" {
			emptyNameParams = append(emptyNameParams, v)
		}
		nameToParams[v.Name] = append(nameToParams[v.Name], v)
	}

	if len(emptyNameParams) > 0 {
		return errors.Errorf("custom header entries have blank names: %+v", emptyNameParams)
	}

	var nameCollisionMsgs []string
	for k, v := range nameToParams {
		if len(v) > 1 {
			nameCollisionMsgs = append(nameCollisionMsgs, fmt.Sprintf("%s: %+v", k, v))
		}
	}
	if len(nameCollisionMsgs) > 0 {
		return errors.Errorf(strings.Join(append([]string{"multiple custom header entries have the same name:"}, nameCollisionMsgs...), "\n\t"))
	}

	// map from path to custom header entries that have the path
	pathsToCustomEntries := make(map[string][]string)
	for _, ch := range headerParams {
		for _, path := range ch.IncludePaths {
			pathsToCustomEntries[path] = append(pathsToCustomEntries[path], ch.Name)
		}
	}
	var customPathCollisionMsgs []string
	sortedKeys := make([]string, 0, len(pathsToCustomEntries))
	for k := range pathsToCustomEntries {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		v := pathsToCustomEntries[k]
		if len(v) > 1 {
			customPathCollisionMsgs = append(customPathCollisionMsgs, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
		}
	}
	if len(customPathCollisionMsgs) > 0 {
		return errors.Errorf(strings.Join(append([]string{"the same path is defined by multiple custom header entries:"}, customPathCollisionMsgs...), "\n\t"))
	}
	return nil
}

type CustomHeaderParam struct {
	// Name is the identifier used to identify this custom license parameter. Must be unique.
	Name string

	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string

	// IncludePaths specifies the paths for which this custom license is applicable. If multiple custom parameters
	// match a file or directory, the parameter with the longest path match is used. If multiple custom parameters
	// match a file or directory exactly (match length is equal), it is treated as an error.
	IncludePaths []string
}

type CustomHeaderConfig struct {
	// Name is the identifier used to identify this custom license parameter. Must be unique.
	Name string `yaml:"name"`

	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string `yaml:"header"`

	// Paths specifies the paths for which this custom license is applicable. If multiple custom parameters match a
	// file or directory, the parameter with the longest path match is used. If multiple custom parameters match a
	// file or directory exactly (match length is equal), it is treated as an error.
	Paths []string `yaml:"paths"`
}

func (c *CustomHeaderConfig) ToParam() CustomHeaderParam {
	return CustomHeaderParam{
		Name:         c.Name,
		Header:       c.Header,
		IncludePaths: c.Paths,
	}
}
