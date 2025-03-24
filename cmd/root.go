package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

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
	debugMode      bool
	version        string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "llmcat",
	Short: "A CLI utility for combining files for LLM input",
	Long:  `A CLI utility that concatenates files with customizable formatting options, suitable for use as input to Large Language Models.`,
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Get version from module info
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
		if version == "" {
			version = "(devel)"
		}
	} else {
		version = "(unknown)"
	}

	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Define global flags
	rootCmd.PersistentFlags().BoolVar(&showFileName, "show-filename", false, "Show file name before content")
	rootCmd.PersistentFlags().BoolVar(&useDashes, "show-dashes", false, "Add dashed lines before and after each file")
	rootCmd.PersistentFlags().StringVar(&contentPrefix, "content-prefix", "", "Text to print before file contents")
	rootCmd.PersistentFlags().StringVar(&contentSuffix, "content-suffix", "", "Text to print after file contents")
	rootCmd.PersistentFlags().StringVar(&fileNamePrefix, "filename-prefix", "", "Text to print before file name")
	rootCmd.PersistentFlags().StringVar(&fileNameSuffix, "filename-suffix", "", "Text to print after file name")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable verbose debug output")
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

// expandPath expands special characters in a path
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	// Handle home directory expansion
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error expanding home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Handle environment variable expansion
	path = os.ExpandEnv(path)

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error converting to absolute path: %w", err)
	}

	return absPath, nil
}

// expandGlobPatterns expands glob patterns into a list of files
func expandGlobPatterns(baseDir string, patterns []string) ([]string, error) {
	var files []string

	// Expand the base directory path
	baseDirExpanded, err := expandPath(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error expanding base directory: %w", err)
	}

	if debugMode {
		fmt.Fprintf(os.Stderr, "Base directory: %s\n", baseDirExpanded)
	}

	// Create a temporary directory to change into
	oldDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}
	defer os.Chdir(oldDir)

	// Change to the base directory
	if err := os.Chdir(baseDirExpanded); err != nil {
		return nil, fmt.Errorf("error changing to base directory: %w", err)
	}

	for _, pattern := range patterns {
		if debugMode {
			fmt.Fprintf(os.Stderr, "\nResolving pattern: %s\n", pattern)
		}

		// Expand special characters in the pattern
		expandedPattern, err := expandPath(pattern)
		if err != nil {
			return nil, fmt.Errorf("error expanding pattern %s: %w", pattern, err)
		}

		if debugMode {
			fmt.Fprintf(os.Stderr, "Expanded pattern: %s\n", expandedPattern)
		}

		// Get the relative pattern
		relPattern, err := filepath.Rel(baseDirExpanded, expandedPattern)
		if err != nil {
			return nil, fmt.Errorf("error getting relative pattern: %w", err)
		}

		if debugMode {
			fmt.Fprintf(os.Stderr, "Relative pattern: %s\n", relPattern)
		}

		matches, err := doublestar.Glob(os.DirFS("."), relPattern)
		if err != nil {
			return nil, fmt.Errorf("error expanding glob pattern %s: %w", pattern, err)
		}

		if debugMode {
			if len(matches) == 0 {
				fmt.Fprintf(os.Stderr, "No matches found\n")
			} else {
				fmt.Fprintf(os.Stderr, "Matches found:\n")
				for _, match := range matches {
					fmt.Fprintf(os.Stderr, "  %s\n", match)
				}
			}
		}

		// Convert matches to absolute paths
		for _, match := range matches {
			absPath := filepath.Join(baseDirExpanded, match)
			files = append(files, absPath)
		}
	}
	return files, nil
}
