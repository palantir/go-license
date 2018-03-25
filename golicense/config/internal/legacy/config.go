// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package legacy

import (
	"github.com/palantir/godel/pkg/versionedconfig"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GoLicense struct {
	versionedconfig.ConfigWithLegacy `yaml:",inline"`

	// Header is the expected license header. All applicable files are expected to start with this header followed
	// by a newline.
	Header string `yaml:"header"`

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []License `yaml:"custom-headers"`

	// Exclude matches the files and directories that should be excluded from consideration for verifying or
	// applying licenses.
	Exclude matcher.NamesPathsCfg `yaml:"exclude"`
}

type License struct {
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

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg GoLicense
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal license-plugin legacy configuration")
	}
	// legacy configuration is completely compatible with v0 configuration
	return cfgBytes, nil
}
