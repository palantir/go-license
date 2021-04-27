// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/palantir/go-license/golicense"
	v0 "github.com/palantir/go-license/golicense/config/internal/v0"
	"github.com/pkg/errors"
)

type ProjectConfig v0.ProjectConfig

func (cfg *ProjectConfig) ToParam() (golicense.ProjectParam, error) {
	customHeaders := make([]golicense.CustomHeaderParam, len(cfg.CustomHeaders))
	for i, v := range cfg.CustomHeaders {
		v := CustomHeaderConfig(v)
		headerVal, err := v.ToParam()
		if err != nil {
			return golicense.ProjectParam{}, err
		}
		customHeaders[i] = headerVal
	}

	if err := validateCustomHeaderParams(customHeaders); err != nil {
		return golicense.ProjectParam{}, err
	}
	return golicense.ProjectParam{
		Licenser:      golicense.NewLicenser(cfg.Header),
		CustomHeaders: customHeaders,
		Exclude:       cfg.Exclude.Matcher(),
	}, nil
}

func validateCustomHeaderParams(headerParams []golicense.CustomHeaderParam) error {
	allNames := make(map[string]struct{})
	collisions := make(map[string]struct{})
	for _, param := range headerParams {
		if _, seen := allNames[param.Name]; seen {
			collisions[param.Name] = struct{}{}
		}
		allNames[param.Name] = struct{}{}
	}
	if len(collisions) > 0 {
		var sortedNames []string
		for k := range collisions {
			sortedNames = append(sortedNames, k)
		}
		sort.Strings(sortedNames)
		return errors.Errorf("custom header(s) defined multiple times: %v", sortedNames)
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

type CustomHeaderConfig v0.CustomHeaderConfig

func ToCustomHeaderConfigs(in []CustomHeaderConfig) []v0.CustomHeaderConfig {
	if in == nil {
		return nil
	}
	out := make([]v0.CustomHeaderConfig, len(in))
	for i, v := range in {
		out[i] = v0.CustomHeaderConfig(v)
	}
	return out
}

func (cfg *CustomHeaderConfig) ToParam() (golicense.CustomHeaderParam, error) {
	if cfg.Name == "" {
		return golicense.CustomHeaderParam{}, errors.Errorf("custom header name cannot be blank")
	}
	return golicense.CustomHeaderParam{
		Name:         cfg.Name,
		Licenser:     golicense.NewLicenser(cfg.Header),
		IncludePaths: cfg.Paths,
	}, nil
}
