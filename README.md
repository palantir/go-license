<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/go-license"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

go-license
==========
`go-license` is a tool that ensures that a license header is applied to Go files.

Usage
-----
Run `./go-license --config=license.yml [files]` to apply the license specified by the configuration in `license.yml` to all of the specified files (only the files that end in `.go` and are not excluded by configuration are processed).

Run `./go-license --config=license.yml --remove [files]` to remove the license specified by the configuration in `license.yml` from all of the specified files (only the files that end in `.go` and are not excluded by configuration are processed).

Run `./go-license --config=license.yml --verify [files]` to verify that the license specified by the configuration is applied to all of the specified files `*.go` files (only the files that end in `.go` and are not excluded by configuration are processed). If the license is not applied properly to any of the files, the files that do not match are printed and the program exits with a non-0 exit code.

Configuration
-------------
The configuration file specifies the header that should be applied as a `header` key. It also supports an `exclude` parameter that specifies files or paths that should be excluded from configuration.

Here is an example configuration file:

```yml
header: |
  // Copyright {{YEAR}} Palantir Technologies, Inc.
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

custom-headers:
  - name: subproject
    header: |
      // Copyright 2016 Palantir Technologies, Inc. All rights reserved.
      // Subproject license.

    paths:
      - subprojectDir
exclude:
  names:
    - "vendor"
```

The string `{{YEAR}}` indicates that, when a license is added by the tool, the current year will be used. For operations that match licenses (for verification or removal), `{{YEAR}}` will match any 4-digit number.

The `custom-headers` configuration allows custom headers to be specified for matching names or paths.
