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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// Meeting is a meet pages entry
type Meeting struct {
	Date         time.Time
	Title        string
	Author       string
	Agenda       []string
	Participants []string
}

const meetTmplFile = "/templates/meet/meeting.md"

var meetTmpl *template.Template

var meetCmd = &cobra.Command{
	Use:   "meet",
	Short: "Creates the shell for meet notes and opens vscode",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		meet := &Meeting{
			Date:         time.Now().UTC(),
			Author:       config.FullName,
			Title:        prompt("title: "),
			Agenda:       repeatPrompt("agenda item: "),
			Participants: repeatPrompt("participant: "),
		}

		fmtTitle, err := formatDirName(meet.Title)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fileName := config.Paths.Meeting + "/" + fmtTitle + ".md"
		defer exec.Command("code", config.Paths.Base, fileName).Run()

		createDirs(config.Paths.Meeting)

		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		if err = meetTmpl.Execute(w, meet); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(meetCmd)

	var err error
	meetTmpl, err = prepareTemplate(meetTmplFile)
	if err != nil {
		panic(err)
	}
}
