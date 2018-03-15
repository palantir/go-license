// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package golicense_test

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/palantir/go-license/golicense"
)

func Example() {
	yml := `
header: |
  // Copyright 2016 Palantir Technologies, Inc.
  //
  // License content.

custom-headers:
  - name: subproject
    header: |
      // Copyright 2016 Palantir Technologies, Inc. All rights reserved.
      // Subproject license.

    paths:
      - subprojectDir
`
	var cfg golicense.ProjectConfig
	if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{Header:// Copyright 2016 Palantir Technologies, Inc.\n//\n// License content.\n CustomHeaders:[{Name:subproject Header:// Copyright 2016 Palantir Technologies, Inc. All rights reserved.\n// Subproject license.\n Paths:[subprojectDir]}] Exclude:{Names:[] Paths:[]}}"
}
