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
