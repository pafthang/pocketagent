package runtime

import (
	"fmt"
	"log"
	"os"
)

// RunMain executes fn and exits with a non-zero status on error.
func RunMain(service string, fn func() error) {
	if err := fn(); err != nil {
		if service != "" {
			fmt.Fprintf(os.Stderr, "%s: %v\n", service, err)
		} else {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}