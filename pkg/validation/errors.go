package validation

import (
	"log"
)

// Calls log.Fatal() with the supplied error if it is not nil
func ExitIfError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
