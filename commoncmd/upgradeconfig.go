// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package commoncmd

import (
	"github.com/palantir/godel/pkg/versionedconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-license/golicense"
	"github.com/palantir/go-license/internal/legacy"
)

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	if legacy.IsLegacyConfig(cfgBytes) {
		return legacy.UpgradeLegacyConfig(cfgBytes)
	}

	version, err := versionedconfig.ConfigVersion(cfgBytes)
	if err != nil {
		return nil, err
	}
	switch version {
	case "", "0":
		var cfg golicense.ProjectConfig
		if err := yaml.UnmarshalStrict(cfgBytes, &cfg); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal input as v0 YAML")
		}
		// input is valid current configuration: return exactly
		return cfgBytes, nil
	default:
		return nil, errors.Errorf("unsupported version: %s", version)
	}
}
