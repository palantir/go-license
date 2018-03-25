// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package golicense

import (
	"github.com/palantir/pkg/matcher"
)

type ProjectParam struct {
	// The default Licenser.
	Licenser Licenser

	// CustomHeaders specifies the custom header parameters. Custom header parameters can be used to specify that
	// certain directories or files in the project should use a header that is different from "Header".
	CustomHeaders []CustomHeaderParam

	// Exclude matches the files and directories that should be excluded from consideration for verifying or applying
	// licenses.
	Exclude matcher.Matcher
}

type CustomHeaderParam struct {
	// Name is the identifier used to identify this custom license parameter. Must be unique.
	Name string

	// Licenser for this parameter.
	Licenser Licenser

	// IncludePaths specifies the paths for which this custom license is applicable. If multiple custom parameters
	// match a file or directory, the parameter with the longest path match is used. If multiple custom parameters
	// match a file or directory exactly (match length is equal), it is treated as an error.
	IncludePaths []string
}
