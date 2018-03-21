// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/palantir/godel/framework/verifyorder"
	"github.com/palantir/pkg/cobracli"
)

var pluginInfo = pluginapi.MustNewPluginInfo(
	"com.palantir.go-license",
	"license-plugin",
	cobracli.Version,
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
		pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions(
			pluginapi.VerifyOptionsOrdering(intVar(verifyorder.License)),
			pluginapi.VerifyOptionsApplyFalseArgs("--verify"),
		)),
	),
	pluginapi.PluginInfoUpgradeConfigTaskInfo(
		pluginapi.UpgradeConfigTaskInfoCommand("upgrade-config"),
	),
)

func intVar(val int) *int {
	return &val
}
