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
	"fmt"
	"log"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

// Chapter represents a chapter in a book
type Chapter struct {
	Number string
	Name   string
	File   string
	Link   string
	Links  []string
}

// Book represents the book getting created
type Book struct {
	Title     string
	Subtitle  string
	Published string
	Authors   []string
	Chapters  []*Chapter
}

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
		defer exec.Command("code", config.Paths.Base, getIndexPath(dir)).Run()

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
	const indexTmplFile = "/templates/book/index.md"
	const chapterTmplFile = "/templates/book/chapter.md"

	rootCmd.AddCommand(bookCmd)
	var err error
	if indexTmpl, err = prepareTemplate(indexTmplFile); err != nil {
		panic(err)
	}
	if chapterTmpl, err = prepareTemplate(chapterTmplFile); err != nil {
		panic(err)
	}
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
		chNum := fmt.Sprintf("%03d", i+1)
		chFile = chNum + "-" + chFile + ".md"
		chLink := "./" + chFile
		book.Chapters = append(book.Chapters, &Chapter{chNum, chName, chFile, chLink, nil})
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
	const imageDir = "/%s/images"
	return dir, createDirs(fmt.Sprintf(config.Paths.Book+imageDir, dir))
}

func createIndex(book *Book, dir string) (err error) {
	filename := getIndexPath(dir)
	_, err = writeTemplate(indexTmpl, filename, book)
	if err != nil {
		return err
	}
	return nil
}

func getIndexPath(dir string) string {
	const indexFile = "/%s/index.md"
	return fmt.Sprintf(config.Paths.Book+indexFile, dir)
}

func createChapters(chapters []*Chapter, dir string) error {
	const chapterFile = "/%s/%s"

	for _, ch := range chapters {
		filename := fmt.Sprintf(config.Paths.Book+chapterFile, dir, ch.File)
		if _, err := writeTemplate(chapterTmpl, filename, ch); err != nil {
			return err
		}
	}
	return nil
}
