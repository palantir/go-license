// Copyright (c) 2016 Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package golicense

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

type Licenser interface {
	// Add adds the license to the provided content.
	Add(content string) string
	// Remove removes the license to the provided content.
	Remove(content string) string
	// Matches returns true if the provided content starts with the license in this Licenser. This is not necessarily
	// a literal prefix match of LicenseHeader (because any year may match).
	Matches(content string) bool
	// Empty returns true if no license header exists.
	Empty() bool
}

// RunLicense runs the license operation using the provided arguments.
func RunLicense(files []string, projectParam ProjectParam, verify, remove bool, stdout io.Writer) error {
	switch {
	case verify:
		if ok, err := VerifyFiles(files, projectParam, stdout); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("")
		}
		return nil
	case remove:
		_, err := UnlicenseFiles(files, projectParam)
		return err
	default:
		_, err := LicenseFiles(files, projectParam)
		return err
	}
}

type licenserImpl struct {
	// literal license to add for new files
	newLicenseHeader string
	// regular expression that matches the license (if nil, the literal content of newLicenseHeader is used)
	matchRegexp *regexp.Regexp
}

func (l *licenserImpl) Add(content string) string {
	return l.newLicenseHeader + "\n" + content
}

func (l *licenserImpl) Remove(content string) string {
	if l.matchRegexp == nil {
		return strings.TrimPrefix(content, l.newLicenseHeader+"\n")
	}
	matchLoc := l.matchRegexp.FindStringIndex(content)
	return content[matchLoc[1]:]
}

func (l *licenserImpl) Matches(content string) bool {
	if l.matchRegexp == nil {
		return strings.HasPrefix(content, l.newLicenseHeader+"\n")
	}
	matchLoc := l.matchRegexp.FindStringIndex(content)
	return len(matchLoc) > 0 && matchLoc[0] == 0
}

func (l *licenserImpl) Empty() bool {
	return l.newLicenseHeader == "" && l.matchRegexp == nil
}

func NewLicenser(license string) Licenser {
	// if special "{{YEAR}}" replacement string is not present, use literal only
	if !strings.Contains(license, "{{YEAR}}") {
		return &licenserImpl{
			newLicenseHeader: license,
		}
	}

	// create a regexp that matches the provided literal header and `\d\d\d\d` for `{{YEAR}}` with a final newline
	parts := strings.Split(license, "{{YEAR}}")
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}

	return &licenserImpl{
		newLicenseHeader: strings.Replace(license, "{{YEAR}}", strconv.Itoa(time.Now().Year()), -1),
		matchRegexp:      regexp.MustCompile(`^` + strings.Join(parts, `\d\d\d\d`) + "\n"),
	}
}

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

func processFiles(files []string, projectParam ProjectParam, modify bool, f func(files []string, licenser Licenser, modify bool) ([]string, error)) ([]string, error) {
	// if header and matchers do not exist, return (nothing to check)
	if projectParam.Licenser.Empty() && len(projectParam.CustomHeaders) == 0 {
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
		currModified, err := f(m[v.Name], v.Licenser, modify)
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
	currModified, err := f(unprocessedGoFiles, projectParam.Licenser, modify)
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

func applyLicenseToFiles(files []string, licenser Licenser, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if !licenser.Matches(content) {
			if modify {
				content = licenser.Add(content)
				if err := ioutil.WriteFile(path, []byte(content), fi.Mode()); err != nil {
					return false, errors.Wrapf(err, "failed to write file %s with new license", path)
				}
			}
			return true, nil
		}
		return false, nil
	})
}

func removeLicenseFromFiles(files []string, licenser Licenser, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if licenser.Matches(content) {
			if modify {
				content = licenser.Remove(content)
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
