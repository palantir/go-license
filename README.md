golicense
=========
`golicense` is a tool that ensures that a license header is applied to all Go files in a project.

Usage
-----
Run `./golicense --config=license.yml` to apply the license specified by the configuration in `license.yml` to all of
the `*.go` files rooted in the current working directory.

Run `./golicense --config=license.yml --verify` to verify that the license specified by the configuration is applied to
all of the `*.go` files rooted in the current working directory. If the license is not applied properly to any of the
files, the files that do not match are printed and the program exits with a non-0 exit code.

Alternatively, a list of files to format may be provided as arguments. The `exclude` filter specified in configuration
will still be applied to paths that are provided as arguments.

Configuration
-------------
The configuration file specifies the header that should be applied as a `header` key. It also supports an `exclude`
parameter that specifies files or paths that should be excluded from configuration.

Here is an example configuration file:

```yml
header: |
        /*
        Copyright 2016 Palantir Technologies, Inc.

        Licensed under the Apache License, Version 2.0 (the "License");
        you may not use this file except in compliance with the License.
        You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

        Unless required by applicable law or agreed to in writing, software
        distributed under the License is distributed on an "AS IS" BASIS,
        WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
        See the License for the specific language governing permissions and
        limitations under the License.
        */
exclude:
  exclude-names:
    - "vendor"
```
