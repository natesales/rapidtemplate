package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// Configuration constants
const (
	PagesDirectory  = "pages/"     // Directory where the markdown pages are located
	OutputDirectory = "out/"       // Output directory for the HTML files and static assets. This should be your webserver root directory.
	TemplateString  = "{{ post }}" // The string to replace in your template file
)

// CLI Usage information
const (
	Usage    = "Usage: rapidtemplate [run/generate/clean]"
	HelpPage = `RapidTemplate (https://github.com/natesales/rapidtemplate)
Usage: rapidtemplate [run/generate/clean]

Commands:
	run      - Build all files and listen for changes 
	generate - Build all files and exit
	clean    - Clean previously built HTML files`
)

// File change operation watcher
var watcher *fsnotify.Watcher

// Basic error handler. If this is invoked, something is really wrong.
func handle(e error) {
	if e != nil {
		panic(e)
	}
}

// Does the file name (with or without path) have a .md extension?
func isMarkdownFile(filename string) bool {
	_filepath := strings.Split(filename, ".")
	if len(_filepath) < 2 {
		return false
	} else {
		return strings.Split(filename, ".")[1] == "md"
	}
}

// Convert a possibly problematic filename into something URL safe. TODO: Make this actually URL safe (special chars)
// Example: "pages/My Cool File.md" -> "my-cool-file.html"
func normalize(filename string) string {
	return strings.Replace(strings.ToLower(strings.Split(strings.Split(filename, "/")[1], ".")[0]), " ", "-", -1) + ".html"
}

// Find and delete existing HTML files
func clean() {
	err := filepath.Walk(OutputDirectory, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			err = os.Remove(path)
			handle(err)
		}
		return nil
	})
	handle(err)
}

// Convert markdown to html
func markdownToHtml(filename string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	file, err := ioutil.ReadFile(filename)
	handle(err)
	return strings.Replace(string(markdown.ToHTML(file, parser, nil)), "</code>\n<code>", "</code>\n<br>\n<code>", -1)
}

// Insert content into template file and return the resulting HTML
func insertIntoTemplate(filename string, content string) string {
	file, err := ioutil.ReadFile(filename)
	handle(err)

	if !strings.Contains(string(file), TemplateString) {
		fmt.Println("ERROR: Can't find \"" + TemplateString + "\" template string.")
		os.Exit(1)
	}

	template := strings.SplitN(string(file), TemplateString, -1)

	output := template[0]
	output += content
	output += template[1]

	return output
}

// Run the conversion process on a file
func update(filename string) {
	fmt.Println("Updating " + filename)

	htmlContent := insertIntoTemplate("template.html", markdownToHtml(filename))

	err := ioutil.WriteFile(OutputDirectory+normalize(filename), []byte(htmlContent), 0644)
	handle(err)
}

// fsnotify walk function
func watchDir(path string, fi os.FileInfo, err error) error {
	if isMarkdownFile(path) {
		update(path)
	}

	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func main() {
	var command string

	if len(os.Args) != 2 {
		fmt.Println(Usage)
		os.Exit(0)
	} else {
		command = os.Args[1]
	}

	switch command {
	case "run":
		// Remove old HTML pages
		clean()

		// Create file change listener
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		// Enumerate (walk) the PagesDirectory to find files
		err := filepath.Walk(PagesDirectory, watchDir)
		handle(err)

		// Blocking channel to allow fsnotify to catch filesystem events
		done := make(chan bool)

		// gorountine for fsnotify
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					if isMarkdownFile(event.Name) {
						update(event.Name)
					}

					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		<-done // Endless blocking channel
	case "clean":
		clean()
		fmt.Println("Done")
	case "generate":
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		err := filepath.Walk(PagesDirectory, watchDir)
		handle(err)
	default:
		fmt.Println("Command \"" + command + "\" not found")
		fmt.Println(HelpPage)
	}
}
