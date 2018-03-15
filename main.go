// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/palantir/pkg/cobracli"
	"github.com/spf13/cobra"

	"github.com/palantir/go-license/cmd"
)

var debugFlagVal bool

func main() {
	rootCmd, rootFlagVals := cmd.RootCmd()
	rootCmd.RunE = func(cobraCmd *cobra.Command, args []string) error {
		projectCfg, err := cmd.LoadConfig(rootFlagVals.CfgFlagVal)
		if err != nil {
			return err
		}
		projectParam, err := projectCfg.ToParam()
		if err != nil {
			return err
		}
		return cmd.RunLicense(*rootFlagVals, args, projectParam, cobraCmd.OutOrStdout())
	}
	os.Exit(cobracli.ExecuteWithDefaultParamsWithVersion(rootCmd, &debugFlagVal, ""))
}
