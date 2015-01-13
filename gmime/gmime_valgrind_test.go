// NOTE: next line must be followed by blank line
// +build valgrind

package gmime_test

import (
	"fmt"
	"github.com/sendgrid/go-gmime/gmime"
	"runtime/debug"
)

// Special function, used for shutdown gmime and free all memory to OS
func ExampleZZZShutdown() {
	fmt.Println("cleanup")
	debug.FreeOSMemory()
	gmime.Shutdown()
	// Output: cleanup
}
