// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config

import (
	"github.com/palantir/godel/v2/pkg/versionedconfig"
	"github.com/pkg/errors"

	"github.com/palantir/go-license/golicense/config/internal/legacy"
	v0 "github.com/palantir/go-license/golicense/config/internal/v0"
)

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	if versionedconfig.IsLegacyConfig(cfgBytes) {
		v0Bytes, err := legacy.UpgradeConfig(cfgBytes)
		if err != nil {
			return nil, err
		}
		cfgBytes = v0Bytes
	}
	version, err := versionedconfig.ConfigVersion(cfgBytes)
	if err != nil {
		return nil, err
	}
	switch version {
	case "", "0":
		return v0.UpgradeConfig(cfgBytes)
	default:
		return nil, errors.Errorf("unsupported version: %s", version)
	}
}
