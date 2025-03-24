# llmcat

A CLI utility to combine multiple files into a single text output, suitable for presenting to an LLM chatbot.

## Installation

```bash
go install github.com/pyoio/llmcat@latest
```

## Usage

The basic usage is:

```bash
llmcat cat <base-dir> [glob-patterns...]
```

### Options

- `--show-filename`: Show file name before content
- `--show-dashes`: Add dashed lines before and after each file
- `--content-prefix`: Text to print before file contents
- `--content-suffix`: Text to print after file contents
- `--filename-prefix`: Text to print before file name
- `--filename-suffix`: Text to print after file name
- `--debug`: Enable verbose debug output showing pattern resolution

### Path Support

The tool supports special characters in paths:
- `~` expands to the user's home directory
- Environment variables are expanded (e.g., `$HOME`, `$USER`)

### Examples

```bash
# Basic usage with a directory and glob pattern
llmcat cat . "*.txt"

# Show file names
llmcat cat --show-filename /path/to/project "**/*.go"

# Add dashed lines around files
llmcat cat --show-dashes ~/Documents "**/*.md"

# Custom file name formatting
llmcat cat --show-filename --filename-prefix "File: " --filename-suffix ":" /path/to/project "**/*.txt"

# Custom content formatting
llmcat cat --content-prefix "Content starts here:\n" --content-suffix "\nContent ends here" . "**/*.md"

# Combine multiple options
llmcat cat --show-filename --show-dashes --filename-prefix "File: " --content-prefix "Content:\n" /path/to/project "**/*.go"

# Find all markdown files recursively in a directory
llmcat cat /path/to/project "**/*.md"

# Multiple glob patterns in a directory
llmcat cat /path/to/project "**/*.md" "**/*.txt"

# Using home directory
llmcat cat ~/Documents "**/*.md"

# Using environment variables
llmcat cat "$HOME/Documents" "**/*.md"

# Recursive with custom formatting
llmcat cat --show-filename --show-dashes --filename-prefix "File: " /path/to/project "**/*.md"

# Debug output to see pattern resolution
llmcat cat --debug /path/to/project "**/*.go" "*.md"
```

## License

MIT License 