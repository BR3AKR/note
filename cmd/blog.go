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

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// BlogEntry represents a single post to the blog
type BlogEntry struct {
	ID       string
	Date     time.Time
	Author   string
	Title    string
	Subtitle string
	Tags     []string
}

const entryTmplFile = "/templates/blog/entry.md"

var tmpl *template.Template

// blogCmd represents the blog command
var blogCmd = &cobra.Command{
	Use:   "blog",
	Short: "Creates the shell for a blog entry and opens vscode",
	Long:  `This command simply creates the .md file for a blog`,
	Run: func(cmd *cobra.Command, args []string) {
		uid, _ := uuid.NewRandom()

		blog := &BlogEntry{
			Date:     time.Now().UTC(),
			Author:   "Sean K Smith",
			Title:    prompt("title: "),
			Subtitle: prompt("subtitle: "),
			Tags:     repeatPrompt("tag: "),
			ID:       uid.String(),
		}

		blogpath := os.Getenv("BLOGPATH")
		if blogpath == "" {
			fmt.Println("BLOGPATH environment variable not set, using home/blog")
			blogpath, _ = os.UserHomeDir()
		}
		fileName := blogpath + "/" + blog.Date.Format("2006-01-02-0304") + ".md"
		defer exec.Command("code", blogpath, fileName).Run()

		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		defer w.Flush()

		if err = tmpl.Execute(w, blog); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(blogCmd)

	var err error
	tmpl, err = prepareTemplate(entryTmplFile)
	if err != nil {
		panic(err)
	}
}
