package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	showFileName   bool
	useDashes      bool
	contentPrefix  string
	contentSuffix  string
	fileNamePrefix string
	fileNameSuffix string
	baseDir        string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "llmcat",
	Short: "A utility to combine multiple files for LLM input",
	Long: `llmcat is a CLI utility that combines multiple files into a single text output,
suitable for presenting to an LLM chatbot. It supports various formatting options for file names and content,
and can find files using glob patterns including double-star patterns (e.g., "**/*.md").`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Define global flags
	rootCmd.PersistentFlags().BoolVarP(&showFileName, "show-filename", "f", false, "Show file name before content")
	rootCmd.PersistentFlags().BoolVarP(&useDashes, "dashes", "d", false, "Add dashed lines before and after each file")
	rootCmd.PersistentFlags().StringVar(&contentPrefix, "content-prefix", "", "Text to print before file contents")
	rootCmd.PersistentFlags().StringVar(&contentSuffix, "content-suffix", "", "Text to print after file contents")
	rootCmd.PersistentFlags().StringVar(&fileNamePrefix, "filename-prefix", "", "Text to print before file name")
	rootCmd.PersistentFlags().StringVar(&fileNameSuffix, "filename-suffix", "", "Text to print after file name")
	rootCmd.PersistentFlags().StringVarP(&baseDir, "base-dir", "b", ".", "Base directory for file search")
}

// checkFileExists checks if a file exists and is not a directory
func checkFileExists(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("error accessing file %s: %w", filename, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", filename)
	}
	return nil
}

// expandGlobPatterns expands glob patterns into a list of files
func expandGlobPatterns(patterns []string) ([]string, error) {
	var files []string
	for _, pattern := range patterns {
		matches, err := doublestar.Glob(os.DirFS(baseDir), pattern)
		if err != nil {
			return nil, fmt.Errorf("error expanding glob pattern %s: %w", pattern, err)
		}
		// Convert matches to absolute paths
		for _, match := range matches {
			files = append(files, filepath.Join(baseDir, match))
		}
	}
	return files, nil
}
