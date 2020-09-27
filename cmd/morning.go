// Copyright © 2019 Sean K Smith <ssmith2347@gmail.com>
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
	"os/exec"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// Morning is a morning pages entry
type Morning struct {
	Date   time.Time
	Author string
}

const morningTmplFile = "/templates/morning/morning.md"

var mornTmpl *template.Template

var morningCmd = &cobra.Command{
	Use:   "morning",
	Short: "Creates the shell for morning pages and opens vscode",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		morning := &Morning{
			Date:   time.Now().UTC(),
			Author: config.FullName,
		}

		fileName := config.Paths.Morning + "/" + morning.Date.Format("2006-01-02-0304") + ".md"
		defer exec.Command("code", config.Paths.Base, fileName).Run()

		createDirs(config.Paths.Morning)
		if _, err := writeTemplate(mornTmpl, fileName, morning); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(morningCmd)

	var err error
	if mornTmpl, err = prepareTemplate(morningTmplFile); err != nil {
		panic(err)
	}
}
