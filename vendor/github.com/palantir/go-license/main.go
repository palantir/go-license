// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/palantir/go-license/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
