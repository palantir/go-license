// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/palantir/godel/framework/pluginapi/v2/pluginapi"
	"github.com/palantir/godel/framework/verifyorder"
)

var (
	Version    = "unspecified"
	PluginInfo = pluginapi.MustNewPluginInfo(
		"com.palantir.go-license",
		"license-plugin",
		Version,
		pluginapi.PluginInfoUsesConfigFile(),
		pluginapi.PluginInfoGlobalFlagOptions(
			pluginapi.GlobalFlagOptionsParamDebugFlag("--"+pluginapi.DebugFlagName),
			pluginapi.GlobalFlagOptionsParamProjectDirFlag("--"+pluginapi.ProjectDirFlagName),
			pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--"+pluginapi.GodelConfigFlagName),
			pluginapi.GlobalFlagOptionsParamConfigFlag("--"+pluginapi.ConfigFlagName),
		),
		pluginapi.PluginInfoTaskInfo(
			"license",
			"Run license task",
			pluginapi.TaskInfoCommand("run"),
			pluginapi.TaskInfoVerifyOptions(
				pluginapi.VerifyOptionsOrdering(intVar(verifyorder.License)),
				pluginapi.VerifyOptionsApplyFalseArgs("--verify"),
			),
		),
		pluginapi.PluginInfoUpgradeConfigTaskInfo(
			pluginapi.UpgradeConfigTaskInfoCommand("upgrade-config"),
			pluginapi.LegacyConfigFile("license.yml"),
		),
	)
)

func intVar(val int) *int {
	return &val
}
