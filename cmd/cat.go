package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:   "cat <base-dir> [glob-patterns...]",
	Short: "Concatenate files for LLM input",
	Long: `Concatenate multiple files into a single output, suitable for LLM input.
You can customize the output format with various flags for file names and content.
The first argument must be a base directory, followed by glob patterns to match files.
Supports glob patterns including double-star patterns (e.g., "**/*.md") to find files recursively.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// First argument is the base directory
		baseDir := args[0]
		patterns := args[1:]

		// Verify base directory exists and is accessible
		baseDirInfo, err := os.Stat(baseDir)
		if err != nil {
			return fmt.Errorf("error accessing base directory %s: %w", baseDir, err)
		}
		if !baseDirInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", baseDir)
		}

		// Expand glob patterns into a list of files
		files, err := expandGlobPatterns(baseDir, patterns)
		if err != nil {
			return err
		}

		// Sort files for consistent output
		sort.Strings(files)

		// Track unique files to avoid duplicates
		processedFiles := make(map[string]bool)

		for _, file := range files {
			// Skip if we've already processed this file
			if processedFiles[file] {
				continue
			}
			processedFiles[file] = true

			// Check if file exists and is not a directory
			if err := checkFileExists(file); err != nil {
				return err
			}

			// Open the file
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", file, err)
			}
			defer f.Close()

			// Add dashed line before file if requested
			if useDashes {
				// Ensure we're on a new line before printing dashes
				fmt.Println()
				fmt.Println(strings.Repeat("-", 80))
			}

			// Handle file name display
			if showFileName {
				// Get current working directory
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("error getting current directory: %w", err)
				}

				// Use relative path from current directory
				relPath, err := filepath.Rel(cwd, file)
				if err != nil {
					relPath = file // Fallback to absolute path if relative path fails
				}

				if fileNamePrefix != "" {
					fmt.Print(fileNamePrefix)
				}
				fmt.Print(relPath)
				if fileNameSuffix != "" {
					fmt.Print(fileNameSuffix)
				}
				fmt.Println()

				// Add blank line after filename
				fmt.Println()
			}

			// Add content prefix if specified
			if contentPrefix != "" {
				fmt.Printf("%s", strings.ReplaceAll(contentPrefix, "\\n", "\n"))
			}

			// Copy file contents to stdout and track if it ends with newline
			content, err := io.ReadAll(f)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", file, err)
			}
			if _, err := os.Stdout.Write(content); err != nil {
				return fmt.Errorf("error writing file contents: %w", err)
			}

			// Add content suffix if specified
			if contentSuffix != "" {
				fmt.Printf("%s", strings.ReplaceAll(contentSuffix, "\\n", "\n"))
			}

			// Add dashed line after file if requested
			if useDashes {
				// Ensure we're on a new line before printing dashes
				if len(content) > 0 && !strings.HasSuffix(string(content), "\n") &&
					(contentSuffix == "" || !strings.HasSuffix(contentSuffix, "\n")) {
					fmt.Println()
				}
				fmt.Println(strings.Repeat("-", 80))
			}
		}

		// Add final newline if not present
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catCmd)
}
