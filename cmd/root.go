// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/palantir/go-license/commoncmd"
	"github.com/palantir/go-license/golicense"
)

var (
	RootCmd = &cobra.Command{
		Use:   "go-license [flags] [files]",
		Short: "Write or verify license headers for Go files",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectCfg, err := commoncmd.LoadConfig(cfgFlagVal)
			if err != nil {
				return err
			}
			projectParam, err := projectCfg.ToParam()
			if err != nil {
				return err
			}
			return golicense.RunLicense(args, projectParam, verifyFlagVal, removeFlagVal, cmd.OutOrStdout())
		},
	}

	cfgFlagVal    string
	verifyFlagVal bool
	removeFlagVal bool
)

func init() {
	RootCmd.Flags().StringVar(&cfgFlagVal, "config", "", "the YAML configuration file for the license check")
	RootCmd.Flags().BoolVar(&verifyFlagVal, "verify", false, "verify that files have proper license headers applied")
	RootCmd.Flags().BoolVar(&removeFlagVal, "remove", false, "remove the license header from files (no-op if verify is true)")
}
