# llmcat

A CLI utility to combine multiple files into a single text output, suitable for presenting to an LLM chatbot.

## Installation

```bash
go install github.com/pyoio/llmcat@latest
```

## Usage

The basic usage is:

```bash
llmcat cat file1.txt file2.txt file3.txt
```

### Options

- `-f, --show-filename`: Show file name before content
- `-d, --dashes`: Add dashed lines before and after each file
- `--content-prefix`: Text to print before file contents
- `--content-suffix`: Text to print after file contents
- `--filename-prefix`: Text to print before file name
- `--filename-suffix`: Text to print after file name
- `-b, --base-dir`: Base directory for file search (default: ".")

### Examples

```bash
# Basic usage
llmcat cat file1.txt file2.txt

# Show file names
llmcat cat -f file1.txt file2.txt

# Add dashed lines around files
llmcat cat -d file1.txt file2.txt

# Custom file name formatting
llmcat cat -f --filename-prefix "File: " --filename-suffix ":" file1.txt file2.txt

# Custom content formatting
llmcat cat --content-prefix "Content starts here:\n" --content-suffix "\nContent ends here" file1.txt file2.txt

# Combine multiple options
llmcat cat -f -d --filename-prefix "File: " --content-prefix "Content:\n" file1.txt file2.txt

# Find all markdown files recursively
llmcat cat "**/*.md"

# Find all Go files in a specific directory
llmcat cat -b /path/to/project "**/*.go"

# Multiple glob patterns
llmcat cat "**/*.md" "**/*.txt"

# Recursive with custom formatting
llmcat cat -f -d --filename-prefix "File: " "**/*.md"
```

## License

MIT License 