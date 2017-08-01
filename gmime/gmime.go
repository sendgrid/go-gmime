package gmime

/*
#cgo pkg-config: gmime-3.0
#include "gmime.h"
*/
import "C"

// This function call automatically by runtime
func init() {
	C.g_mime_init()
}

// Shutdown is really needed only for valgrind
func Shutdown() {
	C.g_mime_shutdown()
}

// convert from Go bool to C gboolean
func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// convert from C gboolean to Go bool
func gobool(b C.gboolean) bool {
	return b != C.gboolean(0)
}

// free up memory
func unref(referee C.gpointer) {
	C.g_object_unref(referee)
}
