// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	godelconfig "github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/pkg/matcher"
	"github.com/spf13/cobra"

	"github.com/palantir/go-license/commoncmd"
	"github.com/palantir/go-license/golicense"
)

var (
	runCmd = &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectCfg, err := commoncmd.LoadConfig(configFlagVal)
			if err != nil {
				return err
			}
			if godelConfigFileFlagVal != "" {
				cfgVal, err := godelconfig.ReadGodelConfigFromFile(godelConfigFileFlagVal)
				if err != nil {
					return err
				}
				projectCfg.Exclude.Add(cfgVal.Exclude)
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
			return golicense.RunLicense(goFiles, projectParam, verifyFlagVal, removeFlagVal, cmd.OutOrStdout())
		},
	}

	verifyFlagVal bool
	removeFlagVal bool
)

func init() {
	runCmd.Flags().BoolVar(&verifyFlagVal, "verify", false, "verify that files have proper license headers applied")
	runCmd.Flags().BoolVar(&removeFlagVal, "remove", false, "remove the license header from files (no-op if verify is true)")
	RootCmd.AddCommand(runCmd)
}
