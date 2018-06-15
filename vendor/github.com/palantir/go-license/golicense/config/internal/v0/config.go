// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package v0

import (
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ProjectConfig struct {
	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline. Any occurrences of the string {{YEAR}} is treated specially: when generating a license, the current
	// year will be substituted for it, and when verifying a license, any 4-digit string will be considered a match.
	Header string `yaml:"header,omitempty"`

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []CustomHeaderConfig `yaml:"custom-headers,omitempty"`

	// Exclude matches the files and directories that should be excluded from consideration for verifying or applying
	// licenses.
	Exclude matcher.NamesPathsCfg `yaml:"exclude,omitempty"`
}

type CustomHeaderConfig struct {
	// Name is the identifier used to identify this custom license parameter. Must be unique.
	Name string `yaml:"name,omitempty"`

	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline. Any occurrences of the string {{YEAR}} is treated specially: when generating a license, the current
	// year will be substituted for it, and when verifying a license, any 4-digit string will be considered a match.
	Header string `yaml:"header,omitempty"`

	// Paths specifies the paths for which this custom license is applicable. If multiple custom parameters match a
	// file or directory, the parameter with the longest path match is used. If multiple custom parameters match a
	// file or directory exactly (match length is equal), it is treated as an error.
	Paths []string `yaml:"paths,omitempty"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var cfg ProjectConfig
	if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal license-plugin v0 configuration")
	}
	return cfgBytes, nil
}
