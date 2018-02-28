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

package golicense

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

func VerifyFiles(files []string, projectParam ProjectParam, stdout io.Writer) (bool, error) {
	// run verify
	modified, err := processFiles(files, projectParam, false, applyLicenseToFiles)
	if err != nil {
		return false, err
	}
	if len(modified) == 0 {
		return true, nil
	}

	var plural string
	if len(modified) == 1 {
		plural = "file does"
	} else {
		plural = "files do"
	}
	parts := append([]string{fmt.Sprintf("%d %s not have the correct license header:", len(modified), plural)}, modified...)
	fmt.Fprintln(stdout, strings.Join(parts, "\n\t"))
	return false, nil
}

func LicenseFiles(files []string, projectParam ProjectParam) ([]string, error) {
	return processFiles(files, projectParam, true, applyLicenseToFiles)
}

func UnlicenseFiles(files []string, projectParam ProjectParam) ([]string, error) {
	return processFiles(files, projectParam, true, removeLicenseFromFiles)
}

func processFiles(files []string, projectParam ProjectParam, modify bool, f func(files []string, header string, modify bool) ([]string, error)) ([]string, error) {
	// if header and matchers do not exist, return (nothing to check)
	if projectParam.Header == "" && len(projectParam.CustomHeaders) == 0 {
		return nil, nil
	}

	goFileMatcher := matcher.Name(`.*\.go`)
	var goFiles []string
	for _, f := range files {
		if goFileMatcher.Match(f) && (projectParam.Exclude == nil || !projectParam.Exclude.Match(f)) {
			goFiles = append(goFiles, f)
		}
	}

	// name of custom matcher -> files to process for the matcher
	m := make(map[string][]string)
	for _, f := range goFiles {
		var longestMatcher string
		longestMatchLen := 0
		for _, v := range projectParam.CustomHeaders {
			for _, p := range v.IncludePaths {
				if matcher.PathLiteral(p).Match(f) && len(p) >= longestMatchLen {
					longestMatcher = v.Name
					longestMatchLen = len(p)
				}
			}
		}
		// file may match multiple custom header params -- if that is the case, use the longest match. Allows
		// for hierarchical matching.
		if longestMatcher != "" {
			m[longestMatcher] = append(m[longestMatcher], f)
		}
	}

	// all files that were processed (considered by a matcher)
	processedFiles := make(map[string]struct{})
	// all files that were modified (or would have been modified)
	var modified []string

	// process custom matchers
	for _, v := range projectParam.CustomHeaders {
		currModified, err := f(m[v.Name], v.Header, modify)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to process headers for matcher %s", v.Name)
		}
		modified = append(modified, currModified...)
		for _, f := range m[v.Name] {
			processedFiles[f] = struct{}{}
		}
	}

	// process all "*.go" files not matched by custom matchers
	var unprocessedGoFiles []string
	for _, f := range goFiles {
		if _, ok := processedFiles[f]; !ok {
			unprocessedGoFiles = append(unprocessedGoFiles, f)
		}
	}
	currModified, err := f(unprocessedGoFiles, projectParam.Header, modify)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process headers for default *.go matcher")
	}
	modified = append(modified, currModified...)
	for _, f := range currModified {
		processedFiles[f] = struct{}{}
	}

	sort.Strings(modified)
	return modified, nil
}

func applyLicenseToFiles(files []string, header string, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if !strings.HasPrefix(content, header+"\n") {
			if modify {
				content = header + "\n" + content
				if err := ioutil.WriteFile(path, []byte(content), fi.Mode()); err != nil {
					return false, errors.Wrapf(err, "failed to write file %s with new license", path)
				}
			}
			return true, nil
		}
		return false, nil
	})
}

func removeLicenseFromFiles(files []string, header string, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if strings.HasPrefix(content, header+"\n") {
			if modify {
				content = strings.TrimPrefix(content, header+"\n")
				if err := ioutil.WriteFile(path, []byte(content), fi.Mode()); err != nil {
					return false, errors.Wrapf(err, "failed to write file %s with license removed", path)
				}
			}
			return true, nil
		}
		return false, nil
	})
}

func visitFiles(files []string, visitor func(path string, fi os.FileInfo, content string) (bool, error)) ([]string, error) {
	var modified []string

	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to stat %s", f)
		}
		bytes, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", f)
		}
		content := string(bytes)
		if changed, err := visitor(f, fi, content); err != nil {
			return nil, errors.WithStack(err)
		} else if changed {
			modified = append(modified, f)
		}
	}

	return modified, nil
}
