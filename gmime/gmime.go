package gmime

/*
#cgo pkg-config: gmime-3.0
#include "gmime.h"
*/
import "C"

// This function call automatically by runtime
func init() {
	C.g_mime_init()
	format := C.g_mime_format_options_get_default()
	C.g_mime_format_options_set_newline_format(format, C.GMIME_NEWLINE_FORMAT_DOS)
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
