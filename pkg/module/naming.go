package module

// Represents a naming scheme for built binaries that determines how their build context is represented
type NamingScheme string

const (
	
	// The build context is not represented at all (i.e. the binary for "mytool" is called "bin/mytool${GOEXE}")
	Undecorated NamingScheme = "Undecorated"
	
	// The build context is represented as a directory prefix (i.e. the binary for "mytool" is called "bin/${GOOS}/${GOARCH}/mytool${GOEXE}")
	PrefixedDirs NamingScheme = "PrefixedDirs"
	
	// The build context is represented as a filename suffix (i.e. the binary for "mytool" is called "bin/mytool-${GOOS}-${GOARCH}${GOEXE}")
	SuffixedFilenames NamingScheme = "SuffixedFilenames"
)
