package network

import "fmt"

// Represents a URL that contains platform-specific substrings
type PlatformSpecificURL struct {
	Pattern string
	Substrings map[string][]string
}

// Resolves the URL for the specified platform
func (url *PlatformSpecificURL) ResolveForPlatform(platform string, architecture string) (string) {
	context := fmt.Sprint(platform, "/", architecture)
	
	substrings := make([]interface{}, len(url.Substrings[context]))
	for i := range url.Substrings[context] {
		substrings[i] = url.Substrings[context][i]
	}
	
	return fmt.Sprintf(url.Pattern, substrings...)
}
