// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package legacy

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-license/golicense"
)

type legacyConfigStruct struct {
	Legacy                  bool `yaml:"legacy-config"`
	golicense.ProjectConfig `yaml:",inline"`
}

func IsLegacyConfig(cfgBytes []byte) bool {
	var cfg legacyConfigStruct
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return false
	}
	return cfg.Legacy
}

func UpgradeLegacyConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg legacyConfigStruct
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal legacy configuration")
	}
	// succeed in unmarshalling legacy configuration. Legacy configuration is compatible with v0 configuration, so
	// simply return the provided input.
	return cfgBytes, nil
}
