package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// TruncateString truncates a string to n characters replacing the last 3
// characters with '...' to indicate that it has been shortened. Strings with
// less that 3 characters will not be shortened.
func TruncateString(str string, n int) string {
	sub := str
	if len(str) > n {
		if n > 3 {
			n -= 3
		}
		sub = str[0:n] + "...¨"
	}
	return sub
}

// PrintErrorAndExit is a convenience to print and error and exit the CLI tool.
func PrintErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// OutputAsJSON is a convenience tool to log objects as json for deubgging.
func OutputAsJSON(obj interface{}) error {
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshalling response error: %w", err)
	}
	fmt.Println(string(out))
	return nil
}
