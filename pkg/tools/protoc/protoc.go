package protoc

import (
	"fmt"
	"runtime"

	"github.com/tensorworks/go-build-helpers/pkg/network"
)

// Returns a PlatformSpecificURL for the specified release of protoc
func URLForRelease(release string) (*network.PlatformSpecificURL) {
	return &network.PlatformSpecificURL{
		Pattern: fmt.Sprint("https://github.com/protocolbuffers/protobuf/releases/download/v", release, "/protoc-", release, "-%s.zip"),
		Substrings: map[string][]string{
			"linux/arm64": []string{"linux-aarch_64"},
			"linux/ppc64le": []string{"linux-ppcle_64"},
			"linux/s390x": []string{"linux-s390x"},
			"linux/386": []string{"linux-x86_32"},
			"linux/amd64": []string{"linux-x86_64"},
			"darwin/amd64": []string{"osx-x86_64"},
			"windows/386": []string{"win32"},
			"windows/amd64": []string{"win64"},
		},
	}
}

// Resolves the URL for the specified release of protoc for the specified platform
func ReleaseForPlatform(release string, platform string, architecture string) (string) {
	url := URLForRelease(release)
	return url.ResolveForPlatform(platform, architecture)
}

// Resolves the URL for the specified release of protoc for the running program's operating system target
func ReleaseForHost(release string) (string) {
	return ReleaseForPlatform(release, runtime.GOOS, runtime.GOARCH)
}
