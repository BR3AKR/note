// Copyright © 2020 Sean K Smith <ssmith2347@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/markbates/pkger"
)

var s *bufio.Scanner

func prepareTemplate(filename string) (*template.Template, error) {
	f, err := pkger.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return template.Must(template.New(f.Info().Name).Parse(string(d))), nil
}

func prompt(prompt string) string {
	fmt.Print(prompt)
	if s == nil {
		s = bufio.NewScanner(os.Stdin)
	}
	s.Scan()
	return strings.TrimSpace(s.Text())
}

func repeatPrompt(promptStr string) []string {
	values := make([]string, 0)

	for c := prompt(promptStr); c != ""; c = prompt(promptStr) {
		values = append(values, c)
	}

	return values
}

func createDirs(dirs ...string) error {
	for i, dir := range dirs {
		err := os.MkdirAll(dir, filePermissions)
		if err != nil {
			if delErr := deleteDirs(dirs[:i]...); delErr != nil {
				err = fmt.Errorf(`error creating directory "%s": %s, unable to clean up: %s`, dirs[i], err.Error(), delErr.Error())
			}
			return err
		}
	}
	return nil
}

func deleteDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func formatDirName(dir string) (string, error) {
	if len(dir) == 0 {
		return "", fmt.Errorf("Invalid directory %s", dir)
	}
	dir = strings.Trim(dir, " ")

	re := regexp.MustCompile(`[^a-zA-Z0-9- ]`)
	dir = re.ReplaceAllString(dir, "")

	re = regexp.MustCompile(`[ ]+`)
	dir = re.ReplaceAllString(dir, "-")

	if len(dir) > 128 {
		dir = dir[:128]
	}

	dir = strings.ToLower(dir)
	return dir, nil
}
