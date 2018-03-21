// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "license",
		Short: "Apply, verify and remove license headers from project files",
	}

	projectDirFlagVal      string
	godelConfigFileFlagVal string
	configFlagVal          string
)

func init() {
	pluginapi.AddProjectDirPFlagPtr(RootCmd.PersistentFlags(), &projectDirFlagVal)
	pluginapi.AddGodelConfigPFlagPtr(RootCmd.PersistentFlags(), &godelConfigFileFlagVal)
	pluginapi.AddConfigPFlagPtr(RootCmd.PersistentFlags(), &configFlagVal)
	if err := RootCmd.MarkPersistentFlagRequired(pluginapi.ProjectDirFlagName); err != nil {
		panic(err)
	}
}
