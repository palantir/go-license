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

// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/palantir/godel/framework/verifyorder"
	"github.com/palantir/pkg/cobracli"
)

var pluginInfo = pluginapi.MustNewInfo(
	"com.palantir",
	"license-plugin",
	cobracli.Version,
	"license.yml",
	pluginapi.MustNewTaskInfo(
		"license",
		"Run license task",
		pluginapi.TaskInfoGlobalFlagOptions(pluginapi.NewGlobalFlagOptions(
			pluginapi.GlobalFlagOptionsParamDebugFlag("--"+pluginapi.DebugFlagName),
			pluginapi.GlobalFlagOptionsParamProjectDirFlag("--"+pluginapi.ProjectDirFlagName),
			pluginapi.GlobalFlagOptionsParamGodelConfigFlag("--"+pluginapi.GodelConfigFlagName),
			pluginapi.GlobalFlagOptionsParamConfigFlag("--"+pluginapi.ConfigFlagName),
		)),
		pluginapi.TaskInfoVerifyOptions(pluginapi.NewVerifyOptions(
			pluginapi.VerifyOptionsOrdering(intVar(verifyorder.License)),
			pluginapi.VerifyOptionsApplyFalseArgs("--verify"),
		)),
	),
)

func intVar(val int) *int {
	return &val
}
