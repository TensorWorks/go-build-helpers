package module

// Represents the options used for building binaries
type BuildOptions struct {
	
	// Additional flags to pass to `go build` when building the binaries
	AdditionalFlags[] string
	
	// The list of build tags to define when building the binaries
	BuildTags []string
	
	// The naming scheme for the built binaries
	Scheme NamingScheme
	
}
