// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package golicense_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/palantir/go-license/golicense"
	"github.com/palantir/go-license/golicense/config"
	"github.com/palantir/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLicenseFiles(t *testing.T) {
	for _, tc := range []struct {
		name         string
		projectParam golicense.ProjectParam
		files        map[string]string
		wantModified []string
		wantContent  map[string]string
	}{
		{
			name: "license applied to Go files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
		},
		{
			name: "license substitutes current year in placeholder",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright {{YEAR}} Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": fmt.Sprintf(`// Copyright %d Palantir Technologies, Inc.
package foo`, time.Now().Year()),
				"bar/bar.go": fmt.Sprintf(`// Copyright %d Palantir Technologies, Inc.
// Original comment
package bar`, time.Now().Year()),
			},
		},
		{
			name: "license not applied to non-Go files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.txt": `package foo`,
			},
			wantContent: map[string]string{
				"foo.txt": `package foo`,
			},
		},
		{
			name: "license not applied to excluded files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
				Exclude:  matcher.Name("foo.go"),
			},
			files: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
			},
			wantContent: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
		},
		{
			name: "license not re-applied to files that already have license",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
			wantModified: []string{
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
		},
		{
			name: "custom license applied to files that match custom matchers",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
				CustomHeaders: []golicense.CustomHeaderParam{
					{
						Name:         "Custom Co.",
						Licenser:     golicense.NewLicenser("// Copyright 2016 Custom Co."),
						IncludePaths: []string{"bar/bar.go"},
					},
					{
						Name:         "Baz",
						Licenser:     golicense.NewLicenser("// Copyright 2006 Legacy Inc."),
						IncludePaths: []string{"baz/baz.go"},
					},
				},
			},
			files: map[string]string{
				"foo.go":     `package foo`,
				"bar/bar.go": `package bar`,
				"baz/baz.go": `package baz`,
			},
			wantModified: []string{
				"bar/bar.go",
				"baz/baz.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Custom Co.
package bar`,
				"baz/baz.go": `// Copyright 2006 Legacy Inc.
package baz`,
			},
		},
		{
			name: "custom matchers match hierarchically",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
				CustomHeaders: []golicense.CustomHeaderParam{
					{
						Name:         "Custom Co.",
						Licenser:     golicense.NewLicenser("// Copyright 2016 Custom Co."),
						IncludePaths: []string{"bar"},
					},
					{
						Name:     "Baz",
						Licenser: golicense.NewLicenser("// Copyright 2006 Legacy Inc."),
						IncludePaths: []string{
							"bar/baz.go",
							"bar/subdir",
						},
					},
				},
			},
			files: map[string]string{
				"foo.go":             `package foo`,
				"bar/bar.go":         `package bar`,
				"bar/baz.go":         `package bar`,
				"bar/subdir/main.go": `package main`,
			},
			wantModified: []string{
				"bar/bar.go",
				"bar/baz.go",
				"bar/subdir/main.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Custom Co.
package bar`,
				"bar/baz.go": `// Copyright 2006 Legacy Inc.
package bar`,
				"bar/subdir/main.go": `// Copyright 2006 Legacy Inc.
package main`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd := chdir(t, tmpDir)
			defer oldWd()

			files := writeFiles(t, tmpDir, tc.files)
			modified, err := golicense.LicenseFiles(files, tc.projectParam)
			require.NoError(t, err)

			assert.Equal(t, tc.wantModified, modified)
			for k, v := range tc.wantContent {
				bytes, err := os.ReadFile(filepath.Join(tmpDir, k))
				require.NoError(t, err)
				require.Equal(t, v, string(bytes))
			}
		})
	}
}

func TestUnlicenseFiles(t *testing.T) {
	for _, tc := range []struct {
		name         string
		projectParam golicense.ProjectParam
		files        map[string]string
		wantModified []string
		wantContent  map[string]string
	}{
		{
			name: "unlicense applied to Go files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
		},
		{
			name: "unlicense applied to Go files with year placeholder",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright {{YEAR}} Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `// Copyright 2018 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
		},
		{
			name: "unlicense not applied to non-Go files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.txt": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
			},
			wantContent: map[string]string{
				"foo.txt": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
			},
		},
		{
			name: "unlicense not applied to excluded files",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
				Exclude:  matcher.Name("foo.go"),
			},
			files: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
			},
			wantContent: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
		},
		{
			name: "unlicense not re-applied to files that already do not have license",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
			},
			files: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Copyright 2016 Palantir Technologies, Inc.
// Original comment
package bar`,
			},
			wantModified: []string{
				"bar/bar.go",
			},
			wantContent: map[string]string{
				"foo.go": `package foo`,
				"bar/bar.go": `// Original comment
package bar`,
			},
		},
		{
			name: "custom license removed from files that match custom matchers",
			projectParam: golicense.ProjectParam{
				Licenser: golicense.NewLicenser(`// Copyright 2016 Palantir Technologies, Inc.`),
				CustomHeaders: []golicense.CustomHeaderParam{
					{
						Name:         "Custom Co.",
						Licenser:     golicense.NewLicenser("// Copyright 2016 Custom Co."),
						IncludePaths: []string{"bar/bar.go"},
					},
					{
						Name:         "Baz",
						Licenser:     golicense.NewLicenser("// Copyright 2006 Legacy Inc."),
						IncludePaths: []string{"baz/baz.go"},
					},
				},
			},
			files: map[string]string{
				"foo.go": `// Copyright 2016 Palantir Technologies, Inc.
package foo`,
				"bar/bar.go": `// Copyright 2016 Custom Co.
package bar`,
				"baz/baz.go": `// Copyright 2006 Legacy Inc.
package baz`,
			},
			wantModified: []string{
				"bar/bar.go",
				"baz/baz.go",
				"foo.go",
			},
			wantContent: map[string]string{
				"foo.go":     `package foo`,
				"bar/bar.go": `package bar`,
				"baz/baz.go": `package baz`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd := chdir(t, tmpDir)
			defer oldWd()

			files := writeFiles(t, tmpDir, tc.files)
			modified, err := golicense.UnlicenseFiles(files, tc.projectParam)
			require.NoError(t, err)

			assert.Equal(t, tc.wantModified, modified)
			for k, v := range tc.wantContent {
				bytes, err := os.ReadFile(path.Join(tmpDir, k))
				require.NoError(t, err)
				assert.Equal(t, v, string(bytes))
			}
		})
	}
}

func TestValidateCustomLicenseParams(t *testing.T) {
	for _, tc := range []struct {
		name          string
		projectConfig config.ProjectConfig
		wantErr       string
	}{
		{
			name:          "empty configuration valid",
			projectConfig: config.ProjectConfig{},
		},
		{
			name: "empty custom configuration name invalid",
			projectConfig: config.ProjectConfig{
				CustomHeaders: config.ToCustomHeaderConfigs([]config.CustomHeaderConfig{
					{
						Header: "// Header",
						Paths:  []string{""},
					},
				}),
			},
			wantErr: "custom header name cannot be blank",
		},
		{
			name: "non-unique custom configuration names invalid",
			projectConfig: config.ProjectConfig{
				CustomHeaders: config.ToCustomHeaderConfigs([]config.CustomHeaderConfig{
					{
						Name:   "foo",
						Header: "// Header",
						Paths:  []string{""},
					},
					{
						Name:   "foo",
						Header: "// Header",
						Paths:  []string{""},
					},
				}),
			},
			wantErr: "custom header(s) defined multiple times: [foo]",
		},
		{
			name: "custom configurations with same paths invalid",
			projectConfig: config.ProjectConfig{
				CustomHeaders: config.ToCustomHeaderConfigs([]config.CustomHeaderConfig{
					{
						Name:   "foo",
						Header: "// Header",
						Paths: []string{
							"foo",
							"bar",
						},
					},
					{
						Name:   "bar",
						Header: "// Header",
						Paths: []string{
							"bar",
							"baz",
						},
					},
					{
						Name:   "ok",
						Header: "// Header",
						Paths: []string{
							"ok",
						},
					},
					{
						Name:   "collides",
						Header: "// Header",
						Paths: []string{
							"bar",
						},
					},
				}),
			},
			wantErr: "the same path is defined by multiple custom header entries:\n\tbar: foo, bar, collides",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.projectConfig.ToParam()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

func chdir(t *testing.T, dest string) func() {
	orig, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(dest)
	require.NoError(t, err)
	return func() {
		if err := os.Chdir(orig); err != nil {
			panic(err)
		}
	}
}

func writeFiles(t *testing.T, root string, files map[string]string) []string {
	dir, err := filepath.Abs(root)
	require.NoError(t, err)

	var writtenFiles []string
	for relPath, content := range files {
		filePath := filepath.Join(dir, relPath)
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
		writtenFiles = append(writtenFiles, relPath)
	}
	return writtenFiles
}
