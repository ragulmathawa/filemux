package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
)

const (
	defaultMaxSize = 500 * 1024 // 500KB max file size
)

func main() {
	// Define flags
	clipboardFlag := flag.Bool("clipboard", false, "Copy output to clipboard instead of printing")
	clipShortFlag := flag.Bool("c", false, "Same as -clipboard")
	forceFlag := flag.Bool("f", false, "Force processing of large files")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: filemux [flags] <file-patterns-or-dirs...>\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  filemux *.txt\n")
		fmt.Fprintf(os.Stderr, "  filemux -c ./docs\n")
		fmt.Fprintf(os.Stderr, "  filemux -f largefile.txt dir/\n")
	}
	flag.Parse()

	// Check if arguments are provided
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Collect all files from arguments
	var files []string
	for _, arg := range flag.Args() {
		// Check if the argument exists as a file or directory
		_, err := os.Stat(arg)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: File or directory '%s' not found\n", arg)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Error checking %s: %v\n", arg, err)
			os.Exit(1)
		}

		info, err := os.Stat(arg)
		if err == nil && info.IsDir() {
			// If argument is a directory, walk it
			err = filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", arg, err)
				continue
			}
		} else {
			// Treat as glob pattern
			matches, err := filepath.Glob(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing pattern %s: %v\n", arg, err)
				continue
			}
			if len(matches) == 0 {
				fmt.Fprintf(os.Stderr, "Error: No files matched pattern '%s'\n", arg)
				os.Exit(1)
			}
			files = append(files, matches...)
		}
	}

	if len(files) == 0 {
		fmt.Println("No files found matching the provided patterns or directories")
		os.Exit(1)
	}

	// Process each file and compile output
	var output strings.Builder
	output.WriteString("Here's the compiled content from multiple files:\n\n")

	for i, file := range files {
		// Check file size
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking %s: %v\n", file, err)
			continue
		}
		if fileInfo.Size() > defaultMaxSize && !*forceFlag {
			fmt.Fprintf(os.Stderr, "Error: File %s is too large (%d bytes, max 500KB). Use -f to force processing.\n", file, fileInfo.Size())
			os.Exit(1)
		}

		content, err := readFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			continue
		}

		// Check if file appears to be binary
		if isBinary(content) {
			fmt.Fprintf(os.Stderr, "Warning: File %s appears to be binary and may not display correctly\n", file)
		}

		// Add file header and content
		output.WriteString(fmt.Sprintf("### File: %s\n", file))
		output.WriteString("```\n")
		output.WriteString(content)
		output.WriteString("\n```\n")

		// Add separator between files (except for the last one)
		if i < len(files)-1 {
			output.WriteString("\n---\n\n")
		}
	}

	// Handle output based on clipboard flags
	if *clipboardFlag || *clipShortFlag {
		err := clipboard.WriteAll(output.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Content successfully copied to clipboard")
	} else {
		fmt.Println(output.String())
	}
}

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// isBinary checks if the content appears to be binary by looking for non-printable characters
func isBinary(content string) bool {
	// Check first 1024 bytes or entire content if smaller
	checkLen := len(content)
	if checkLen > 1024 {
		checkLen = 1024
	}

	for _, r := range content[:checkLen] {
		if r < 32 && !unicode.IsSpace(r) && r != '\r' && r != '\n' && r != '\t' {
			return true
		}
	}
	return false
}
