package cmd

import (
	"bufio"
	"os"
	"testing"
)

func TestFormatDirName(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"A simple chapter", "a-simple-chapter"},
		{"A chapter with \"quotes\"", "a-chapter-with-quotes"},
		{"A chapt-er with-hyphens", "a-chapt-er-with-hyphens"},
		{"A chapter   with too     many spaces", "a-chapter-with-too-many-spaces"},
		{"This extraordinary chapter has so many characters that it exceeds the 128 character limit. Full file lengths are typically limited to 255, but that includes files I believe.", "this-extraordinary-chapter-has-so-many-characters-that-it-exceeds-the-128-character-limit-full-file-lengths-are-typically-limite"},
		{"This chapter: has some `strange characters` -- which \"should be removed\"", "this-chapter-has-some-strange-characters----which-should-be-removed"},
		{`  Some  Directory -Name

		dir-name`, "some-directory--namedir-name"},
	}

	for _, test := range tests {
		got, err := formatDirName(test.input)
		if err != nil {
			t.Errorf("formatDirName(%s) returned err %v", test.input, err)
			continue
		}
		if got != test.want {
			t.Errorf("formatDirName(%s) = %s wanted %s", test.input, got, test.want)
		}
	}
}

func TestFormatDirNameInvalid(t *testing.T) {
	_, err := formatDirName("")
	if err == nil {
		t.Errorf("expected error not returned")
	}
}

func TestCreateDirectories(t *testing.T) {
	const testDir = "test-dir"

	if err := createDirs(testDir); err != nil {
		t.Fatalf("unexpected error creating directories: %v", err)
	}

	if _, err := os.Stat("./" + testDir); os.IsNotExist(err) {
		t.Fatalf("directory not created: %v", err)
	}

	const filepath = "./" + testDir + "/test.txt"
	f, err := os.Create(filepath)
	defer os.RemoveAll("./" + testDir)
	defer f.Close()

	if err != nil {
		t.Fatalf("unable to create file in directory: %v", err)
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString("some data")

	if err != nil {
		t.Errorf("unable to write to file: %v", err)
	}

	w.Flush()
}
