// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package commoncmd

import (
	"io/ioutil"
	"os"

	"github.com/palantir/go-license/golicense/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func LoadConfig(cfgFile string) (config.ProjectConfig, error) {
	cfgYML, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		return config.ProjectConfig{}, nil
	}
	if err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
	}

	upgradedBytes, err := config.UpgradeConfig(cfgYML)
	if err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
	}

	var cfg config.ProjectConfig
	if err := yaml.Unmarshal(upgradedBytes, &cfg); err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to unmarshal configuration as YAML")
	}
	return cfg, nil
}
