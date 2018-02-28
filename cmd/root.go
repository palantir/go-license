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

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/go-license/golicense"
)

type RootFlags struct {
	CfgFlagVal    string
	VerifyFlagVal bool
	RemoveFlagVal bool
}

// RootCmd returns a cobra.Command that acts as the root command for the "license" program. Sets the usage strings and
// hooks up common flags and returns a struct that will contain the value of the flags on execution. Note that the
// command does not have any Run actions defined, so the caller should define the run action.
func RootCmd() (*cobra.Command, *RootFlags) {
	rootFlags := &RootFlags{}
	cmd := &cobra.Command{
		Use:   "license [flags] [files]",
		Short: "Write or verify license headers for Go files",
	}
	cmd.Flags().StringVar(&rootFlags.CfgFlagVal, "config", "", "the YAML configuration file for the license check")
	cmd.Flags().BoolVar(&rootFlags.VerifyFlagVal, "verify", false, "verify that files have proper license headers applied")
	cmd.Flags().BoolVar(&rootFlags.RemoveFlagVal, "remove", false, "remove the license header from files (no-op if verify is true)")
	return cmd, rootFlags
}

// RunLicense runs the license CLI operation using the provided arguments. Meant to be used as the core of any "Run"
// actions defined for RootCmd().
func RunLicense(flags RootFlags, args []string, projectParam golicense.ProjectParam, stdout io.Writer) error {
	switch {
	case flags.VerifyFlagVal:
		if ok, err := golicense.VerifyFiles(args, projectParam, stdout); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("")
		}
		return nil
	case flags.RemoveFlagVal:
		_, err := golicense.UnlicenseFiles(args, projectParam)
		return err
	default:
		_, err := golicense.LicenseFiles(args, projectParam)
		return err
	}
}

func LoadConfig(cfgFile string) (golicense.ProjectConfig, error) {
	cfgYML, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return golicense.ProjectConfig{}, errors.Wrapf(err, "failed to read file %s", cfgFile)
	}
	var cfg golicense.ProjectConfig
	if err := yaml.Unmarshal(cfgYML, &cfg); err != nil {
		return golicense.ProjectConfig{}, errors.Wrapf(err, "failed to unmarshal configuration as YAML")
	}
	return cfg, nil
}
