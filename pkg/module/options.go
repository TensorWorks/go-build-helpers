package module

// Represents the options used for building binaries
type BuildOptions struct {
	
	// The list of build tags to define when building the binaries
	BuildTags []string
	
	// The naming scheme for the built binaries
	Scheme NamingScheme
	
}
