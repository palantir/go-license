// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"
	"path"

	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/palantir/pkg/cobracli"
	"github.com/palantir/pkg/matcher"
	"github.com/spf13/cobra"

	"github.com/palantir/go-license/cmd"
)

var (
	debugFlagVal           bool
	projectDirFlagVal      string
	godelConfigFileFlagVal string
)

func main() {
	if ok := pluginapi.InfoCmd(os.Args, os.Stdout, pluginInfo); ok {
		return
	}

	rootCmd, rootFlagVals := cmd.RootCmd()

	// add flags for godel config and project directory
	pluginapi.AddGodelConfigPFlagPtr(rootCmd.PersistentFlags(), &godelConfigFileFlagVal)
	pluginapi.AddProjectDirPFlagPtr(rootCmd.PersistentFlags(), &projectDirFlagVal)
	if err := rootCmd.MarkPersistentFlagRequired(pluginapi.ProjectDirFlagName); err != nil {
		panic(err)
	}

	rootCmd.RunE = func(cobraCmd *cobra.Command, args []string) error {
		projectCfg, err := cmd.LoadConfig(rootFlagVals.CfgFlagVal)
		if err != nil {
			return err
		}
		// if godel config is specified, add in the "exclude" configuration it provides
		if godelConfigFileFlagVal != "" {
			cfg, err := godellauncher.ReadGodelConfig(path.Dir(godelConfigFileFlagVal))
			if err != nil {
				return err
			}
			projectCfg.Exclude.Add(cfg.Exclude)
		}
		projectParam, err := projectCfg.ToParam()
		if err != nil {
			return err
		}
		// plugin matches all Go files in project except for those excluded by configuration
		goFiles, err := godellauncher.ListProjectPaths(projectDirFlagVal, matcher.Name(`.*\.go`), projectParam.Exclude)
		if err != nil {
			return err
		}
		return cmd.RunLicense(*rootFlagVals, goFiles, projectParam, cobraCmd.OutOrStdout())
	}
	os.Exit(cobracli.ExecuteWithDefaultParamsWithVersion(rootCmd, &debugFlagVal, ""))
}
