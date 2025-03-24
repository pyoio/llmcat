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
	Use:   "cat [files or patterns...]",
	Short: "Concatenate files for LLM input",
	Long: `Concatenate multiple files into a single output, suitable for LLM input.
You can customize the output format with various flags for file names and content.
Supports glob patterns including double-star patterns (e.g., "**/*.md") to find files recursively.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Expand glob patterns into a list of files
		files, err := expandGlobPatterns(args)
		if err != nil {
			return err
		}

		// Sort files for consistent output
		sort.Strings(files)

		for _, file := range files {
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
				fmt.Println(strings.Repeat("-", 80))
			}

			// Handle file name display
			if showFileName {
				// Use relative path from baseDir
				relPath, err := filepath.Rel(baseDir, file)
				if err == nil {
					if fileNamePrefix != "" {
						fmt.Print(fileNamePrefix)
					}
					fmt.Print(relPath)
					if fileNameSuffix != "" {
						fmt.Print(fileNameSuffix)
					}
					fmt.Println()
				}
			}

			// Add content prefix if specified
			if contentPrefix != "" {
				fmt.Print(contentPrefix)
			}

			// Copy file contents to stdout
			if _, err := io.Copy(os.Stdout, f); err != nil {
				return fmt.Errorf("error reading file %s: %w", file, err)
			}

			// Add content suffix if specified
			if contentSuffix != "" {
				fmt.Print(contentSuffix)
			}

			// Add dashed line after file if requested
			if useDashes {
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
