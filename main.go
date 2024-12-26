package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <mode> [args...]")
		fmt.Println("Modes:")
		fmt.Println("  snapshot <root_path> [ignore_patterns...]")
		fmt.Println("  outdir <desired_output_dirname> // write the snapshot to the input.md file before running this mode.")
		os.Exit(1)
	}

	// Get mode from first argument
	mode := os.Args[1]

	switch mode {
	case "snapshot":
		if len(os.Args) < 3 {
			fmt.Println("Usage: program snapshot <root_path> [ignore_patterns...]")
			os.Exit(1)
		}
		rootPath := os.Args[2]
		ignoreList := os.Args[3:]
		err := generateMarkdownSnapshot(rootPath, ignoreList)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Snapshot created successfully - go to the output.md file to see the results.")
	case "outdir":
		fmt.Println("outdir mode selected... initiating outdir processes...")

		if len(os.Args) < 3 {
			fmt.Println("Usage: program outdir <desired_output_dirname>")
			fmt.Println("(╯°□°)╯︵ ┻━┻  Please provide an output directory name!")
			os.Exit(1)
		}
		dirname := os.Args[2]
		err := generateOutdir(dirname)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Output directory created successfully at directory name '%s'\n", dirname)
	default:
		fmt.Printf("Unknown mode '%s'\n", mode)
		os.Exit(1)
	}
}

// generateMarkdownSnapshot creates a Markdown file documenting the directory contents
func generateMarkdownSnapshot(rootPath string, ignoreList []string) error {
	outputFile, err := os.Create("output.md")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// Function to determine if the path should be ignored
	shouldIgnore := func(path string) bool {
		// Normalize the path to use forward slashes for consistent handling
		normalizedPath := filepath.ToSlash(path)

		for _, ignore := range ignoreList {
			// Normalize the ignore pattern to use forward slashes
			normalizedIgnore := filepath.ToSlash(ignore)

			// Match ignore patterns exactly from the root relative path
			trimmedPath := strings.TrimPrefix(normalizedPath, filepath.ToSlash(rootPath)+"/")
			if trimmedPath == normalizedIgnore || strings.HasPrefix(trimmedPath, normalizedIgnore+"/") {
				return true
			}
		}
		return false
	}

	// Walk the directory tree
	// func WalkDirAndWrite(dir string, info os.FileInfo, err error) error
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}
		if shouldIgnore(relativePath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Write the path and file content to the Markdown file
			fmt.Fprintf(writer, "### %s\n```\n%s\n```\n\n", relativePath, fileContent)
		}
		return nil
	})

	return err
}

// generateOutdir creates a directory with the `desired_output_dirname` name that accords with the snapshot in the input.md file
func generateOutdir(dirname string) error {
	// Create the root directory (with style! 💃)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w (ノಠ益ಠ)ノ彡┻━┻", err)
	}

	// Read our precious input.md file (🤞 hoping it exists!)
	content, err := os.ReadFile("input.md")
	if err != nil {
		return fmt.Errorf("couldn't read input.md - did you forget to create it? ¯\\_(ツ)_/¯ : %w", err)
	}

	// Split content into sections (each file is a section)
	sections := strings.Split(string(content), "### ")

	// Skip the first empty section (it's empty because Split creates an empty first element)
	for _, section := range sections[1:] {
		// Split into filename and content
		parts := strings.SplitN(section, "\n```\n", 2)
		if len(parts) != 2 {
			continue // Skip malformed sections (we're forgiving! 🤗)
		}

		filename := strings.TrimSpace(parts[0])
		// Extract content between ``` markers
		content := strings.Split(parts[1], "```")[0]

		// Create the full path for the file
		fullPath := filepath.Join(dirname, filename)

		// Create parent directories if they don't exist
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directories for %s: %w (ಥ﹏ಥ)", filename, err)
		}

		// Write the file (drumroll please! 🥁)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w (╯°□°）╯︵ ┻━┻", filename, err)
		}
	}

	return nil // Mission accomplished! 🎉
}
