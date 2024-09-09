# BashLike Go Package

This package provides Go functions that mimic common bash commands, allowing you to write bash-like scripts in Go.

## But why?

Here are some reasons why you might want to use this package:

1. **Type safety**: Go provides strong typing, catching errors at compile-time.
2. **Better tooling**: Go offers excellent IDE support, testing frameworks, and debugging tools.
3. **Cross-platform compatibility**: Go programs can be easily compiled for different operating systems.
4. **Performance**: Go typically outperforms Bash scripts, especially for complex operations.
5. **Maintainability**: Go's structured approach and package system make it easier to manage large codebases.
6. **Standard library**: Go's rich standard library reduces dependency on external tools.
7. **Concurrency**: Go's goroutines and channels make parallel processing simpler.
8. **Error handling**: Go's explicit error handling leads to more robust code.

This project aims to combine the simplicity of Bash-like commands with the power and safety of Go, creating a more reliable and efficient scripting alternative.

## Usage

Import the package in your Go script:

```go
import "github.com/zerocorebeta/bashlike"
```

## Installation

To install the BashLike package, use `go get`:

```bash
go get github.com/zerocorebeta/bashlike
```

## Features

- Bash-like command functions (e.g., `Cat`, `Grep`, `Ls`, `Mkdir`, `Rm`)
- I/O operations (`ReadLine`, `WriteFile`, `AppendFile`)
- Text processing (`Cut`, `Sed`, `Awk`, `Sort`, `Uniq`, `Wc`)
- File path operations (`Basename`, `Dirname`)
- Command execution with context support
- Concurrent map implementation
- Pipe-like command chaining

## Example

Here's an example that demonstrates how to use this package to analyze log files, `log_analyzer.go`:

```go

#!/usr/bin/env go run

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/zerocorebeta/bashlike"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./log_analyzer <log_directory>")
		os.Exit(1)
	}

	logDir := os.Args[1]
	files, err := bashlike.Ls(logDir)
	if err != nil {
		log.Fatal(err)
	}

	statusCodes := make(map[string]int)

	for _, file := range files {
		if filepath.Ext(file) == ".log" {
			content, err := bashlike.Cat(filepath.Join(logDir, file))
			if err != nil {
				log.Printf("Error reading file %s: %v", file, err)
				continue
			}

			lines, err := bashlike.Grep("HTTP/1", content)
			if err != nil {
				log.Printf("Error grepping file %s: %v", file, err)
				continue
			}

			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) > 8 {
					statusCode := fields[8]
					statusCodes[statusCode]++
				}
			}
		}
	}

	var results []string
	for code, count := range statusCodes {
		results = append(results, fmt.Sprintf("%s: %d", code, count))
	}

	sort.Strings(results)

	report := strings.Join(results, "\n")
	err = bashlike.WriteFile("status_report.txt", report)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Report generated: status_report.txt")
}
```

We can run it like this:
```bash
chmod +x log_analyzer.go
./log_analyzer.go
```

This script analyzes log files in a specified directory, extracts information about HTTP status codes, and generates a summary report.

## Error Handling

The package uses custom error types for better error handling:

- `ErrInvalidRegex`: Invalid regular expression
- `ErrCommandExecution`: Error executing a command
- `ErrIO`: I/O-related errors
- `ErrInvalidArgument`: Invalid function arguments

Always check returned errors and handle them appropriately in your code.

## Documentation

For detailed documentation of all available functions, please refer to the [GoDoc](https://pkg.go.dev/github.com/zerocorebeta/bashlike) page.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
