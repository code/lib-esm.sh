package cli

import (
	"encoding/base64"
	"flag"
	"os"
	"strings"

	"golang.org/x/term"
)

var (
	moduleExts = []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx", ".svelte", ".vue"}
)

// termRaw implements the github.com/ije/gox/term.Raw interface.
type termRaw struct{}

func (t *termRaw) Next() byte {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		panic(err)
	}

	// The third byte is the key specific value we are looking for.
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if n == 3 {
		return buf[2]
	}

	return buf[0]
}

// isHttpSepcifier returns true if the specifier is a remote URL.
func isHttpSepcifier(specifier string) bool {
	return strings.HasPrefix(specifier, "https://") || strings.HasPrefix(specifier, "http://")
}

// isRelPathSpecifier returns true if the specifier is a local path.
func isRelPathSpecifier(specifier string) bool {
	return strings.HasPrefix(specifier, "./") || strings.HasPrefix(specifier, "../")
}

// isAbsPathSpecifier returns true if the specifier is an absolute path.
func isAbsPathSpecifier(specifier string) bool {
	return strings.HasPrefix(specifier, "/") || strings.HasPrefix(specifier, "file://")
}

// endsWith returns true if the given string ends with any of the suffixes.
func endsWith(s string, suffixs ...string) bool {
	for _, suffix := range suffixs {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// btoaUrl converts a string to a base64 string.
func btoaUrl(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// atobUrl converts a base64 string to a string.
func atobUrl(s string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// parseCommandFlag parses the command flag.
func parseCommandFlag() (string, []string) {
	flag.CommandLine.Parse(os.Args[2:])

	args := make([]string, 0, len(os.Args)-2)
	nextVaule := false
	for _, arg := range os.Args[2:] {
		if !strings.HasPrefix(arg, "-") {
			if !nextVaule {
				args = append(args, arg)
			} else {
				nextVaule = false
			}
		} else if !strings.Contains(arg, "=") {
			nextVaule = true
		}
	}
	if len(args) == 0 {
		return "", nil
	}
	return args[0], args[1:]
}

// stringInSlice returns true if the given value is included in the array.
func stringInSlice(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
