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
	"log"
	"os"
	"text/template"

	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"
)

// Chapter represents a chapter in a book
type Chapter struct {
	Number int
	Name   string
	File   string
	Link   string
	Links  []string
}

// Book represents the book getting created
type Book struct {
	Title, Subtitle, Published string
	Authors                    []string
	Chapters                   []*Chapter
}

const indexTmplFile = "./book/index.md"
const chapterTmplFile = "./book/chapter.md"
const indexFile = "./%s/index.md"
const chapterFile = "./%s/%s"

var indexTmpl *template.Template
var chapterTmpl *template.Template

// bookCmd represents the book command
var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Creates a template for taking notes about a book",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		book, err := promptForBook()
		if err != nil {
			log.Fatalf("error prompting for book: %v, aborting", err)
		}

		dir, err := buildDirStruct(book.Title)
		if err != nil {
			log.Fatalf("error building directory structure: %v", err)
		}

		err = createIndex(book, dir)
		if err != nil {
			log.Fatalf("error creating index: %v", err)
		}

		err = createChapters(book.Chapters, dir)
		if err != nil {
			log.Fatalf("error creating chapters: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)
	box := packr.NewBox("../templates")
	indexTmpl = template.Must(template.New("index.md").Parse(box.String(indexTmplFile)))
	chapterTmpl = template.Must(template.New("chapter.md").Parse(box.String(chapterTmplFile)))
}

func promptForBook() (*Book, error) {
	book := &Book{
		Title:     prompt("title: "),
		Subtitle:  prompt("subtitle: "),
		Published: prompt("published: "),
		Authors:   repeatPrompt("author: "),
	}

	links := make([]string, 0)

	for i, chName := range repeatPrompt("chapter title: ") {
		chFile, err := formatDirName(chName)
		if err != nil {
			return nil, err
		}
		chFile = chFile + ".md"
		chLink := "./" + chFile
		book.Chapters = append(book.Chapters, &Chapter{i + 1, chName, chFile, chLink, nil})
		links = append(links, fmt.Sprintf("[%d](%s)", i+1, chLink))
	}

	for _, ch := range book.Chapters {
		ch.Links = links
	}

	return book, nil
}

func buildDirStruct(title string) (string, error) {
	dir, err := formatDirName(title)
	if err != nil {
		return "", err
	}
	err = createDirs([]string{fmt.Sprintf("./%s/images", dir)})
	return dir, err
}

func createIndex(book *Book, dir string) error {
	idxf, err := os.Create(fmt.Sprintf(indexFile, dir))
	if err != nil {
		return err
	}

	defer idxf.Close()
	w := bufio.NewWriter(idxf)
	defer w.Flush()
	return indexTmpl.Execute(w, book)
}

func createChapters(chapters []*Chapter, dir string) error {
	for _, ch := range chapters {
		chf, err := os.Create(fmt.Sprintf(chapterFile, dir, ch.File))
		if err != nil {
			return err
		}

		defer chf.Close()
		w := bufio.NewWriter(chf)
		defer w.Flush()

		err = chapterTmpl.Execute(w, ch)
		if err != nil {
			return err
		}
	}
	return nil
}
