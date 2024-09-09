package bashlike

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Error types
var (
	ErrInvalidRegex     = errors.New("invalid regex pattern")
	ErrCommandExecution = errors.New("error executing command")
	ErrIO               = errors.New("I/O error")
	ErrInvalidArgument  = errors.New("invalid argument")
)

// Echo prints arguments to stdout.
func Echo(args ...interface{}) error {
	_, err := fmt.Println(args...)
	return err
}

// Cat reads a file and returns its content as a string.
func Cat(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrIO, err)
	}
	return string(content), nil
}

// Grep searches for a pattern in a string and returns matching lines.
func Grep(pattern, text string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRegex, err)
	}
	var matches []string
	for _, line := range strings.Split(text, "\n") {
		if re.MatchString(line) {
			matches = append(matches, line)
		}
	}
	return matches, nil
}

// Ls lists files in a directory.
func Ls(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIO, err)
	}
	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names, nil
}

// Mkdir creates a directory.
func Mkdir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// Rm removes a file or directory.
func Rm(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// Pwd returns the current working directory.
func Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrIO, err)
	}
	return dir, nil
}

// Cd changes the current working directory.
func Cd(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// Exec executes a command and returns its output.
func Exec(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCommandExecution, err)
	}
	return string(output), nil
}

// ReadLine reads a line from stdin.
func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrIO, err)
	}
	return strings.TrimSpace(line), nil
}

// WriteFile writes content to a file.
func WriteFile(filename, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// AppendFile appends content to a file.
func AppendFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// Env returns the value of an environment variable.
func Env(key string) string {
	return os.Getenv(key)
}

// SetEnv sets an environment variable.
func SetEnv(key, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrIO, err)
	}
	return nil
}

// Cut extracts sections from each line of input.
func Cut(input, delimiter string, fields []int) []string {
	var result []string
	for _, line := range strings.Split(input, "\n") {
		parts := strings.Split(line, delimiter)
		var selected []string
		for _, field := range fields {
			if field > 0 && field <= len(parts) {
				selected = append(selected, parts[field-1])
			}
		}
		result = append(result, strings.Join(selected, delimiter))
	}
	return result
}

// Sed performs simple string substitutions.
func Sed(input, old, new string) string {
	return strings.ReplaceAll(input, old, new)
}

// Awk simulates basic awk functionality.
func Awk(input string, pattern string, action func([]string) string) []string {
	var result []string
	re := regexp.MustCompile(pattern)
	for _, line := range strings.Split(input, "\n") {
		if re.MatchString(line) {
			fields := strings.Fields(line)
			result = append(result, action(fields))
		}
	}
	return result
}

// Find simulates the find command.
func Find(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIO, err)
	}
	return matches, nil
}

// Xargs simulates the xargs command.
func Xargs(ctx context.Context, input []string, command string, args ...string) (string, error) {
	var output strings.Builder
	for _, item := range input {
		cmdArgs := append(args, item)
		cmd := exec.CommandContext(ctx, command, cmdArgs...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrCommandExecution, err)
		}
		output.Write(out)
	}
	return output.String(), nil
}

// Sort sorts lines of text.
func Sort(input []string) []string {
	sorted := make([]string, len(input))
	copy(sorted, input)
	sort.Strings(sorted)
	return sorted
}

// Uniq removes adjacent duplicate lines.
func Uniq(input []string) []string {
	var result []string
	for i, line := range input {
		if i == 0 || line != input[i-1] {
			result = append(result, line)
		}
	}
	return result
}

// Wc counts lines, words, and characters.
func Wc(input string) (lines, words, chars int) {
	lines = strings.Count(input, "\n")
	words = len(strings.Fields(input))
	chars = len(input)
	return
}

// Basename returns the base name of a file path.
func Basename(path string) string {
	return filepath.Base(path)
}

// Dirname returns the directory name of a file path.
func Dirname(path string) string {
	return filepath.Dir(path)
}

// Tr translates or deletes characters.
func Tr(input, from, to string) (string, error) {
	if len(from) != len(to) {
		return "", fmt.Errorf("%w: 'from' and 'to' must have the same length", ErrInvalidArgument)
	}
	tr := make(map[rune]rune)
	for i, r := range from {
		tr[r] = rune(to[i])
	}
	return strings.Map(func(r rune) rune {
		if v, ok := tr[r]; ok {
			return v
		}
		return r
	}, input), nil
}

// Head returns the first n lines of input.
func Head(input string, n int) string {
	lines := strings.SplitN(input, "\n", n+1)
	return strings.Join(lines[:min(n, len(lines))], "\n")
}

// Tail returns the last n lines of input.
func Tail(input string, n int) string {
	lines := strings.Split(input, "\n")
	start := max(0, len(lines)-n)
	return strings.Join(lines[start:], "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Test simulates the test command for file operations and string comparisons.
func Test(condition string, args ...string) (bool, error) {
	switch condition {
	case "-e":
		_, err := os.Stat(args[0])
		return err == nil, nil
	case "-f":
		info, err := os.Stat(args[0])
		return err == nil && !info.IsDir(), nil
	case "-d":
		info, err := os.Stat(args[0])
		return err == nil && info.IsDir(), nil
	case "-z":
		return len(args[0]) == 0, nil
	case "-n":
		return len(args[0]) > 0, nil
	case "=":
		return args[0] == args[1], nil
	case "!=":
		return args[0] != args[1], nil
	default:
		return false, fmt.Errorf("%w: unsupported test condition: %s", ErrInvalidArgument, condition)
	}
}

// Expr evaluates a simple arithmetic expression.
func Expr(ctx context.Context, expression string) (int, error) {
	output, err := Exec(ctx, "expr", strings.Fields(expression)...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrCommandExecution, err)
	}
	result, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	return result, nil
}

// Pipe represents a command that can be piped.
type Pipe struct {
	Cmd  func(context.Context, io.Reader) (io.Reader, error)
	Next *Pipe
}

// Execute runs the pipe chain.
func (p *Pipe) Execute(ctx context.Context, input io.Reader) (io.Reader, error) {
	var err error
	for p != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			input, err = p.Cmd(ctx, input)
			if err != nil {
				return nil, err
			}
			p = p.Next
		}
	}
	return input, nil
}

// ConcurrentMap is a thread-safe map implementation.
type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// NewConcurrentMap creates a new ConcurrentMap.
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		items: make(map[string]interface{}),
	}
}

// Set adds or updates an item in the map.
func (m *ConcurrentMap) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = value
}

// Get retrieves an item from the map.
func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.items[key]
	return value, ok
}

// Delete removes an item from the map.
func (m *ConcurrentMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}
